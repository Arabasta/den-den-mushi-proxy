package iexpress

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/core/host"
	"den-den-mushi-Go/internal/control/core/iexpress"
	"den-den-mushi-Go/internal/control/core/implementor_groups"
	"den-den-mushi-Go/internal/control/core/os_adm_users"
	"den-den-mushi-Go/internal/control/core/pty_sessions"
	"den-den-mushi-Go/internal/control/filters"
	oapi "den-den-mushi-Go/openapi/control"
	hostpkg "den-den-mushi-Go/pkg/dto/host"
	ptysessionspkg "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/util/cyberark"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"strings"
)

type Service struct {
	iexpressSvc    *iexpress.Service
	ptySessionsSvc *pty_sessions.Service
	hostSvc        *host.Service
	impGrpSvc      *implementor_groups.Service
	osAdmUsersSvc  *os_adm_users.Service

	log *zap.Logger
	cfg *config.Config
}

func NewService(iexpSvc *iexpress.Service, ptySessionsSvc *pty_sessions.Service,
	hostSvc *host.Service, impGrpSvc *implementor_groups.Service, osAdmUsersSvc *os_adm_users.Service, log *zap.Logger, cfg *config.Config) *Service {
	log.Info("Initializing IExpress Service")
	return &Service{
		iexpressSvc:    iexpSvc,
		ptySessionsSvc: ptySessionsSvc,
		hostSvc:        hostSvc,
		impGrpSvc:      impGrpSvc,
		osAdmUsersSvc:  osAdmUsersSvc,
		log:            log,
		cfg:            cfg,
	}
}

func (s Service) ListExpressRequests(f filters.ListIexpress, c *gin.Context) (*oapi.GetIExpressResponse, error) {
	// only show tickets for user's implementor groups
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
	//s.log.Debug("implementor groups for user", zap.Any("groups", userImplGroups))
	impGroups := make([]string, 0)
	for _, group := range userImplGroups {
		impGroups = append(impGroups, group.GroupName)
	}
	f.ImplementorGroups = EnrichImpGroupsWithGOV_(impGroups)

	// get total count
	var totalCount int
	if f.IsGetTotalCount {
		totalCount, err = s.iexpressSvc.CountApprovedByFilter(f)
		if err != nil {
			s.log.Error("CountApprovedIExpressByFilter", zap.Error(err))
			return nil, err
		}
		//	s.log.Debug("Total count of iexpress requests", zap.Int("count", totalCount))
	}

	// fetch IExpress requests using filter
	exps, err := s.iexpressSvc.FindApprovedByFilter(f)
	if err != nil {
		s.log.Error("FindApprovedByFilter", zap.Error(err))
		return nil, err
	}
	items := make([]oapi.IExpress, 0, len(exps))

	for _, exp := range exps {
		items = append(items, oapi.IExpress{
			RequestId:       exp.RequestId,
			Requestor:       exp.Requestor,
			Lob:             exp.Lob,
			OriginCountry:   exp.OriginCountry,
			AppImpacted:     exp.AppImpacted,
			StartTime:       *exp.StartTime,
			EndTime:         *exp.EndTime,
			State:           exp.State,
			ApproverGroup1:  &exp.ApproverGroup1,
			ApproverGroup2:  &exp.ApproverGroup2,
			MdApproverGroup: &exp.MDApproverGroup,
			RelatedTicket:   &exp.RelatedTicket,
			Action:          &exp.Action,
		})
	}

	var totalCountPtr *int
	if f.IsGetTotalCount {
		totalCountPtr = &totalCount
	}

	return &oapi.GetIExpressResponse{
		TotalCount: totalCountPtr,
		Page:       &f.Page,
		PageSize:   &f.PageSize,
		Items:      &items,
	}, nil
}

func (s Service) ListExpressRequestDetails(id string, c *gin.Context) (*oapi.GetIExpressHostsAndSessionDetailsResponse, error) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		s.log.Error("Auth context missing in request")
		return nil, errors.New("auth context missing in request")
	}
	exp, err := s.iexpressSvc.FindByTicketNumber(id)
	if err != nil {
		s.log.Error("Failed to find IExpress by ticket number", zap.String("ticket", id), zap.Error(err))
		return nil, err
	}

	// extract ips from exp
	ipToUsers := cyberark.MapIPToOSUsers(exp.CyberArkObjects)
	//s.log.Debug("Mapped CyberArk object", zap.Any("object", exp.CyberArkObjects))

	ips := make([]string, 0, len(ipToUsers))
	for ip := range ipToUsers {
		ips = append(ips, ip)
	}
	//s.log.Debug("Mapped IPs to OS Users", zap.Strings("ips", ips))

	// get host details
	hosts, err := s.hostSvc.FindAllLinuxOsByIps(ips)
	if err != nil {
		s.log.Error("Failed to fetch hosts by IPs", zap.Error(err))
		return nil, err
	}

	hostMap := make(map[string]*hostpkg.Record)
	for _, h := range hosts {
		hostMap[h.IpAddress] = h
	}
	//s.log.Debug("Mapped hosts by IPs", zap.Any("hostMap", hostMap))

	sessions, err := s.ptySessionsSvc.FindAllByChangeRequestIDAndServerIPs(exp.RequestId, ips)
	if err != nil {
		s.log.Error("Failed to fetch PTY sessions", zap.Error(err))
		return nil, err
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
	//s.log.Debug("Grouped PTY sessions by IP", zap.Any("ipSessionsMap", ipSessionsMap))

	var hostDetails []oapi.HostSessionDetailsV2
	for ip, hostRec := range hostMap {
		hostInfo := &oapi.Host{
			AppCode:     hostRec.Appcode,
			Environment: hostRec.Environment,
			IpAddress:   hostRec.IpAddress,
			Name:        hostRec.HostName,
		}

		// for prod environment, filter out non-prod hosts
		// for non-prod, yirong say don't care
		if s.cfg.App.Environment == "prod" {
			if hostInfo.Environment != "PROD" {
				continue
			}
		}

		// todo refactor this, it is garbage
		osUsers, ok := ipToUsers[ip]
		if !ok {
			osUsers = []string{}
		}
		extraUsers := s.osAdmUsersSvc.GetNonCrOsUsers(authCtx.UserID)
		userSet := make(map[string]struct{}, len(osUsers)+len(extraUsers))
		var merged []string
		for _, u := range append(osUsers, extraUsers...) {
			if _, exists := userSet[u]; !exists && u != "" {
				userSet[u] = struct{}{}
				merged = append(merged, u)
			}
		}
		osUsers = merged
		// remove root users so they dont click and err
		filtered := make([]string, 0, len(merged))
		for _, u := range merged {
			if u != "root" {
				filtered = append(filtered, u)
			}
		}

		osUsers = filtered

		var sessions []*ptysessionspkg.Record
		if agg, ok := ipSessionsMap[ip]; ok {
			sessions = agg.sessions
		}

		hostDetails = append(hostDetails, oapi.HostSessionDetailsV2{
			Host:            hostInfo,
			OsUsers:         &osUsers,
			CyberarkObjects: cyberark.ExtractObjectsForIp(hostInfo.IpAddress, exp.CyberArkObjects),
			PtySessions:     convertToPtySessionSummaries(sessions),
		})
	}

	return &oapi.GetIExpressHostsAndSessionDetailsResponse{
		HostSessionDetails: &hostDetails,
	}, nil
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

// EnrichImpGroupsWithGOV_ makes the implementor a GOV guy, if he is an INF guy
func EnrichImpGroupsWithGOV_(groups []string) *[]string {
	var r []string

	for _, g := range groups {
		r = append(r, g)
		infGroup := g
		if strings.Contains(g, "INF_") {
			infGroup = strings.Replace(g, "INF_", "GOV_", 1)
		}

		if len(infGroup) > 0 && infGroup[0] != 'P' {
			infGroup = "P" + infGroup
		}

		if infGroup != g {
			r = append(r, infGroup)
		}
	}

	log.Debug("EnrichImpGroupsWithGOV_", zap.Any("groups", r))
	return &r
}
