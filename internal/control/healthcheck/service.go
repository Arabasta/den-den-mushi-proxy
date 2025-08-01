package healthcheck

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/filters"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/pty_sessions"
	oapi "den-den-mushi-Go/openapi/control"
	ptysessionspkg "den-den-mushi-Go/pkg/dto/pty_sessions"
	"go.uber.org/zap"
)

type Service struct {
	ptySessionsSvc *pty_sessions.Service
	hostSvc        *host.Service

	log *zap.Logger
	cfg *config.Config
}

func NewService(ptySessionsSvc *pty_sessions.Service,
	hostSvc *host.Service, log *zap.Logger, cfg *config.Config) *Service {
	log.Info("Initializing Make Change Service")
	return &Service{
		ptySessionsSvc: ptySessionsSvc,
		hostSvc:        hostSvc,
		log:            log,
		cfg:            cfg,
	}
}

func (s *Service) getHostsAndAssociatedPtySessions(f filters.HealthcheckPtySession) (*[]oapi.HostSessionDetails, error) {
	s.log.Debug("Fetching hosts and PTY sessions", zap.Any("filter", f))

	// todo: FindAllByFilter will eventually require ougroup
	hosts, err := s.hostSvc.FindAllByFilter(f)
	s.log.Debug("Got hosts", zap.Any("filter", f), zap.Any("hosts", len(hosts)))
	if err != nil {
		return nil, err
	}

	hostips := make([]string, len(hosts))
	for i, h := range hosts {
		hostips[i] = h.IpAddress
	}

	sessions, err := s.ptySessionsSvc.FindAllByStartConnServerIpsAndState(hostips, f.PtySessionState)
	s.log.Debug("Got sessions", zap.Any("filter", f), zap.Any("sessions", len(sessions)))
	if err != nil {
		s.log.Error("Failed to fetch PTY sessions", zap.Error(err))
		return nil, err
	}

	var result []oapi.HostSessionDetails
	for _, h := range hosts {
		hostSessions := filterSessionsForHost(sessions, h.IpAddress)

		details := oapi.HostSessionDetails{
			Host: &oapi.Host{
				Name:        h.HostName,
				IpAddress:   h.IpAddress,
				AppCode:     h.Appcode,
				Environment: h.Environment,
				Country:     &h.Country,
			},
			PtySessions: convertToPtySessionSummaries(hostSessions),
			OsUsers:     &s.cfg.Development.HealthcheckOsUsers,
		}

		result = append(result, details)
	}

	s.log.Debug("Returning host session details", zap.Int("count", len(result)))
	return &result, nil
}

func filterSessionsForHost(sessions []*ptysessionspkg.Record, ip string) []*ptysessionspkg.Record {
	var filtered []*ptysessionspkg.Record
	for _, s := range sessions {
		if s.StartConnServerIP == ip {
			filtered = append(filtered, s)
		}
	}
	return filtered
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
