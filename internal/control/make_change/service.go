package make_change

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/filters"
	"den-den-mushi-Go/internal/control/pty_sessions"
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/middleware"
	"go.uber.org/zap"
)

type Service struct {
	crSvc          *change_request.Service
	ptySessionsSvc *pty_sessions.Service

	log *zap.Logger
}

func NewService(crSvc *change_request.Service, ptySessionsSvc *pty_sessions.Service, log *zap.Logger) *Service {
	log.Info("Initializing Make Change Service")
	return &Service{
		crSvc:          crSvc,
		ptySessionsSvc: ptySessionsSvc,
		log:            log,
	}
}

func (s *Service) ListChangeRequestsWithSessions(filter filters.ListCR, authCtx *middleware.AuthContext) ([]oapi.ChangeRequestSessionsResponse, error) {
	//var r []oapi.ChangeRequestSessionsResponse
	//
	//// todo: verify user permissions ? is this needed?
	//
	//// fetch CRs using filter
	//crs, err := s.crSvc.FindChangeRequests(filter)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil

	//
	//for _, cr := range crs {
	//	// 2. extract hosts ips from CR cyberark objects
	//
	//
	//	// 3. get details for each host
	//
	//	// 4. For each host, get the PTY sessions using start_conn_cr_id and start_conn_server_ip
	//
	//	sessions, err := s.ptySessionsSvc.FindSessionsByChangeRequestID(cr.CRNumber)
	//	if err != nil {
	//		continue
	//	}
	//
	//	// for each session, get its connections
	//	r = append(r, oapi.ChangeRequestSessionsResponse{
	//		ChangeId:          &cr.CRNumber,
	//		ChangeStartTime:   parseTimePtr(cr.ChangeStartTime),
	//		ChangeEndTime:     parseTimePtr(cr.ChangeEndTime),
	//		ChangeRequestStatus: &cr.State,
	//		Summary:           &cr.Summary,
	//		Country:           &cr.Country,
	//		ImplementorGroups: splitAndPtr(cr.ImplementorGroups),
	//		Lob:               &cr.Lob,
	//		HostSessionDetails: &[]HostSessionDetails{...},
	//	})
	//
	//
	//	type hostAggregate struct {
	//		sessions []*pty_sessions.Model
	//		users    map[string]struct{}
	//	}
	//
	//	hostMap := map[string]*hostAggregate{}
	//	hostDetails := []HostSessionDetails{}
	//	for ip, agg := range hostMap {
	//		host := fetchProxyHostByIP(ip)
	//
	//		osUsers := make([]string, 0, len(agg.users))
	//		for user := range agg.users {
	//			osUsers = append(osUsers, user)
	//		}
	//
	//		hostDetails = append(hostDetails, HostSessionDetails{
	//			Host: &struct {
	//				AppCode     string `json:"app_code"`
	//				Environment string `json:"environment"`
	//				IpAddress   string `json:"ip_address"`
	//				Name        string `json:"name"`
	//			}{
	//				AppCode:     host.AppCode,
	//				Environment: host.Environment,
	//				IpAddress:   host.IpAddress,
	//				Name:        host.HostName,
	//			},
	//			OsUsers:     &osUsers,
	//			PtySessions: convertToPtySessionSummaries(agg.sessions),
	//		})
	//	}
	//}
	//
	//return r, nil
}
