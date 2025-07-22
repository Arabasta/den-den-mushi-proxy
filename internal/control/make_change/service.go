package make_change

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/filters"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/pty_sessions"
	oapi "den-den-mushi-Go/openapi/control"
	hostpkg "den-den-mushi-Go/pkg/dto/host"
	ptysessionspkg "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/util/cyberark"
	"go.uber.org/zap"
)

type Service struct {
	crSvc          *change_request.Service
	ptySessionsSvc *pty_sessions.Service
	hostSvc        *host.Service

	log *zap.Logger
}

func NewService(crSvc *change_request.Service, ptySessionsSvc *pty_sessions.Service,
	hostSvc *host.Service, log *zap.Logger) *Service {
	log.Info("Initializing Make Change Service")
	return &Service{
		crSvc:          crSvc,
		ptySessionsSvc: ptySessionsSvc,
		hostSvc:        hostSvc,
		log:            log,
	}
}

// todo refactor garbage
func (s *Service) ListChangeRequestsWithSessions(filter filters.ListCR, authCtx *middleware.AuthContext) ([]oapi.ChangeRequestSessionsResponse, error) {
	var r []oapi.ChangeRequestSessionsResponse
	// todo: verify user permissions ? is this needed?

	// fetch CRs using filter
	crs, err := s.crSvc.FindChangeRequests(filter)
	if err != nil {
		return nil, err
	}

	for _, cr := range crs {
		// extract ips from cr
		ipToUsers := cyberark.MapIPToOSUsers(cr.CyberArkObjects)
		if len(ipToUsers) == 0 {
			continue
		}

		ips := make([]string, 0, len(ipToUsers))
		for ip := range ipToUsers {
			ips = append(ips, ip)
		}

		// get host details
		hosts, err := s.hostSvc.FindAllByIps(ips)
		if err != nil {
			s.log.Error("Failed to fetch hosts by IPs", zap.Error(err))
			continue
		}

		hostMap := make(map[string]*hostpkg.Record)
		for _, h := range hosts {
			hostMap[h.IpAddress] = h
		}

		sessions, err := s.ptySessionsSvc.FindAllByChangeRequestIDAndServerIPs(cr.ChangeRequestId, ips)
		if err != nil {
			s.log.Error("Failed to fetch PTY sessions", zap.Error(err))
			continue
		}

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

		var hostDetails []oapi.HostSessionDetails
		for k, v := range ipSessionsMap {
			h := hostMap[k]
			if h == nil {
				continue
			}

			osUsers, ok := ipToUsers[k]
			if !ok {
				osUsers = []string{}
			}

			hostDetails = append(hostDetails, oapi.HostSessionDetails{
				Host: &oapi.Host{
					AppCode:     h.Appcode,
					Environment: h.Environment,
					IpAddress:   h.IpAddress,
					Name:        h.HostName,
				},
				OsUsers:     &osUsers,
				PtySessions: convertToPtySessionSummaries(v.sessions),
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

	return r, nil
}

type hostAggregate struct {
	sessions []*ptysessionspkg.Record
}

func convertToPtySessionSummaries(sessions []*ptysessionspkg.Record) *[]oapi.PtySessionSummary {
	if len(sessions) == 0 {
		return nil
	}
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
