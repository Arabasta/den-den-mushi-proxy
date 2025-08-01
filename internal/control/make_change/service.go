package make_change

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/filters"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/pty_sessions"
	oapi "den-den-mushi-Go/openapi/control"
	hostpkg "den-den-mushi-Go/pkg/dto/host"
	ptysessionspkg "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/util/cyberark"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service struct {
	crSvc          *change_request.Service
	ptySessionsSvc *pty_sessions.Service
	hostSvc        *host.Service
	impGrpSvc      *implementor_groups.Service

	log *zap.Logger
}

func NewService(crSvc *change_request.Service, ptySessionsSvc *pty_sessions.Service,
	hostSvc *host.Service, impGrpSvc *implementor_groups.Service, log *zap.Logger) *Service {
	log.Info("Initializing Make Change Service")
	return &Service{
		crSvc:          crSvc,
		ptySessionsSvc: ptySessionsSvc,
		hostSvc:        hostSvc,
		impGrpSvc:      impGrpSvc,
		log:            log,
	}
}

// todo refactor garbage, need to make it 1 query ... will do it when there are no more changes or maybe not
// todo return cyberark objects
func (s *Service) ListChangeRequestsWithSessions(filter filters.ListCR, c *gin.Context) ([]oapi.ChangeRequestSessionsResponse, error) {
	var r []oapi.ChangeRequestSessionsResponse

	// only show crs for user's implementor groups
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		s.log.Error("Auth context missing in request")
		return nil, errors.New("auth context missing in request")
	}
	userImplGroups, err := s.impGrpSvc.FindAllByUserId(authCtx.UserID)
	if err != nil {
		s.log.Error("Failed to fetch user implementor groups", zap.Error(err))
		return nil, err
	}
	s.log.Debug("implementor groups for user", zap.Any("groups", userImplGroups))
	impGroups := make([]string, 0)
	for _, group := range userImplGroups {
		impGroups = append(impGroups, group.GroupName)
	}
	filter.ImplementorGroups = &impGroups

	// fetch CRs using filter
	crs, err := s.crSvc.FindChangeRequestsByFilter(filter)
	//s.log.Debug("CRs fetched", zap.Int("count", len(crs)))
	if err != nil {
		s.log.Error("FindChangeRequestsByFilter", zap.Error(err))
		return nil, err
	}

	for _, cr := range crs {
		//	s.log.Debug("Mapping CR", zap.Any("cr", cr))
		// extract ips from cr
		ipToUsers := cyberark.MapIPToOSUsers(cr.CyberArkObjects)
		s.log.Debug("Mapped CyberArk object", zap.Any("object", cr.CyberArkObjects))

		ips := make([]string, 0, len(ipToUsers))
		for ip := range ipToUsers {
			ips = append(ips, ip)
		}
		s.log.Debug("Mapped IPs to OS Users", zap.Strings("ips", ips))

		// get host details
		hosts, err := s.hostSvc.FindAllLinuxOsByIps(ips)
		if err != nil {
			s.log.Error("Failed to fetch hosts by IPs", zap.Error(err))
			continue
		}
		//	s.log.Debug("Fetched hosts by IPs", zap.Any("hosts", hosts))

		hostMap := make(map[string]*hostpkg.Record)
		for _, h := range hosts {
			hostMap[h.IpAddress] = h
		}
		s.log.Debug("Mapped hosts by IPs", zap.Any("hostMap", hostMap))

		sessions, err := s.ptySessionsSvc.FindAllByChangeRequestIDAndServerIPs(cr.ChangeRequestId, ips)
		if err != nil {
			s.log.Error("Failed to fetch PTY sessions", zap.Error(err))
			continue
		}
		//s.log.Debug("Fetched PTY sessions", zap.Any("sessions", sessions))

		// group sessions by StartConnServerIP
		ipSessionsMap := map[string]*hostAggregate{}
		for _, s := range sessions {
			if _, ok := ipSessionsMap[s.StartConnServerIP]; !ok {
				ipSessionsMap[s.StartConnServerIP] = &hostAggregate{
					sessions: []*ptysessionspkg.Record{},
				}
			}
			ipSessionsMap[s.StartConnServerIP].sessions = append(ipSessionsMap[s.StartConnServerIP].sessions, s)
		}
		s.log.Debug("Grouped PTY sessions by IP", zap.Any("ipSessionsMap", ipSessionsMap))

		var hostDetails []oapi.HostSessionDetails
		for ip, hostRec := range hostMap {
			hostInfo := &oapi.Host{
				AppCode:     hostRec.Appcode,
				Environment: hostRec.Environment,
				IpAddress:   hostRec.IpAddress,
				Name:        hostRec.HostName,
			}

			osUsers, ok := ipToUsers[ip]
			if !ok {
				osUsers = []string{}
			}

			var sessions []*ptysessionspkg.Record
			if agg, ok := ipSessionsMap[ip]; ok {
				sessions = agg.sessions
			}

			hostDetails = append(hostDetails, oapi.HostSessionDetails{
				Host:        hostInfo,
				OsUsers:     &osUsers,
				PtySessions: convertToPtySessionSummaries(sessions),
			})
		}

		r = append(r, oapi.ChangeRequestSessionsResponse{
			ChangeId:            &cr.ChangeRequestId,
			ChangeStartTime:     cr.ChangeStartTime,
			ChangeEndTime:       cr.ChangeEndTime,
			ChangeRequestStatus: &cr.State,
			Summary:             &cr.Summary,
			Country:             &cr.Country,
			ImplementorGroups:   &cr.ImplementorGroups,
			Lob:                 &cr.Lob,
			HostSessionDetails:  &hostDetails,
		})
	}

	//s.log.Debug("Returning change request sessions response", zap.Any("response", r))
	return r, nil
}

type hostAggregate struct {
	sessions []*ptysessionspkg.Record
}

func convertToPtySessionSummaries(sessions []*ptysessionspkg.Record) *[]oapi.PtySessionSummary {
	out := make([]oapi.PtySessionSummary, 0, len(sessions))
	for _, s := range sessions {
		var conns []oapi.Connection
		for _, c := range s.Connections {
			conns = append(conns, oapi.Connection{
				Id:           &c.ID,
				JoinTime:     c.JoinTime,
				LeaveTime:    c.LeaveTime,
				PtySessionId: &c.PtySessionID,
				StartRole:    (*oapi.StartRole)(&c.StartRole),
				Status:       (*oapi.ConnectionStatus)(&c.Status),
				UserId:       &c.UserID,
			})
		}
		out = append(out, oapi.PtySessionSummary{
			ChangeId:     &s.StartConnChangeRequestID,
			Connections:  conns,
			CreatedBy:    s.CreatedBy,
			EndTime:      s.EndTime,
			Id:           s.ID,
			LastActivity: s.LastActivity,
			Purpose:      oapi.ConnectionPurpose(s.StartConnPurpose),
			StartTime:    s.StartTime,
			State:        oapi.PtySessionState(s.State),
		})
	}
	return &out
}
