package pty_token

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/internal/control/proxy_lb"
	"den-den-mushi-Go/internal/control/pty_sessions"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type Service struct {
	issuer *jwt.Issuer

	psSvc   *pty_sessions.Service
	plbSvc  *proxy_lb.Service
	hostSvc *host.Service
	crSvc   *change_request.Service

	changeRequestPolicyChain policy.Policy[request.Ctx]
	healthCheckPolicyChain   policy.Policy[request.Ctx]

	log *zap.Logger
	cfg *config.Config
}

func NewService(psS *pty_sessions.Service, plbS *proxy_lb.Service, hostS *host.Service,
	issuer *jwt.Issuer, crS *change_request.Service,
	changeRequestPolicyChain policy.Policy[request.Ctx],
	healthCheckPolicyChain policy.Policy[request.Ctx],
	log *zap.Logger, cfg *config.Config) *Service {
	log.Info("Initializing PTY Token Service")
	return &Service{
		psSvc:                    psS,
		plbSvc:                   plbS,
		hostSvc:                  hostS,
		issuer:                   issuer,
		changeRequestPolicyChain: changeRequestPolicyChain,
		healthCheckPolicyChain:   healthCheckPolicyChain,
		crSvc:                    crS,
		log:                      log,
		cfg:                      cfg,
	}
}

// todo: split into healthcheck and cr endpoints
func (s *Service) mintStartToken(r wrapper.WithAuth[request.StartRequest]) (string, string, error) {
	hostConnMethod, hostType := types.LocalSshKey, types.OS

	if !s.cfg.Development.IsLocalSshKeyIfNotIsPuppetKey {
		hostConnMethod = types.SshOrchestratorKey
	}

	var err error

	// todo: enable all these
	// hostConnMethod, err := s.hostS.FindHostConnectionMethodByIp(r.Server.IP) todo: grab server conn method how?

	//hostType, err := s.hostSvc.FindTypeByIp(r.Body.Server.IP)
	//if err != nil {
	//	s.log.Error("Failed to find host type by IP", zap.String("ip", r.Body.Server.IP), zap.Error(err))
	//	return "", "", err
	//}

	var filter types.Filter
	var cr *changerequestpkg.Record
	adapter := &request.StartAdapter{Req: r}

	if r.Body.Purpose == types.Change {
		cr, err = s.crSvc.FindByTicketNumber(r.Body.ChangeID)
		if err != nil || cr == nil {
			s.log.Error("Failed to find change request by ID", zap.String("changeID", r.Body.ChangeID), zap.Error(err))
			return "", "", err
		}

		adapter.CR = cr

		s.log.Debug("Starting CR policy check")
		if err = s.changeRequestPolicyChain.Check(adapter); err != nil {
			s.log.Warn("Change request policy check failed", zap.Error(err))
			return "", "", err
		}
	} else if r.Body.Purpose == types.Healthcheck {
		s.log.Debug("Starting healthcheck policy check")
		if err = s.healthCheckPolicyChain.Check(adapter); err != nil {
			s.log.Warn("Health check policy check failed", zap.Error(err))
			return "", "", err
		}

		if s.cfg.Development.IsBlacklistFilter {
			filter = types.Blacklist // todo: get filter type by host type
		} else {
			filter = types.Whitelist
		}

		//filter, err = s.hostSvc.FindFilterTypeByHostType(hostType)
		//if err != nil {
		//	s.log.Error("Failed to find filter type by host type", zap.String("hostType", string(hostType)), zap.Error(err))
		//	return "", "", err
		//}
	}

	s.log.Debug("Building connection for start")
	conn := jwt.BuildConnForStart(hostConnMethod, r, cr, filter, s.cfg.Development.TargetSshPort)

	tok, err := s.issuer.Mint(r.AuthCtx, conn, hostType)
	if err != nil {
		s.log.Error("Failed to mint token", zap.Error(err))
		return "", "", err
	}

	proxyLbUrl, err := s.plbSvc.GetLBEndpointByProxyType(hostType)
	if err != nil {
		s.log.Error("Failed to get proxy load balancer URL", zap.String("hostType", string(hostType)), zap.Error(err))
		return "", "", err
	}

	return tok, proxyLbUrl, nil
}

func (s *Service) mintJoinToken(r wrapper.WithAuth[request.JoinRequest]) (string, string, error) {
	ps, err := s.psSvc.FindById(r.Body.PtySessionId)
	if err != nil || ps == nil {
		s.log.Error("Failed to find pty session", zap.String("ptySessionId", r.Body.PtySessionId), zap.Error(err))
		return "", "", errors.New("failed to find pty session")
	}

	crId, err := s.getChangeRequestIDOrError(ps.StartConnPurpose, ps.StartConnChangeRequestID)
	if err != nil {
		s.log.Error("Invalid connection details", zap.String("ptySessionId", r.Body.PtySessionId), zap.Error(err))
		return "", "", err
	}

	adapter := &request.JoinAdapter{
		Req: r,
		AdapterFields: request.AdapterFields{
			Purpose:  ps.StartConnPurpose,
			ChangeID: crId,
			Server: dto.ServerInfo{
				IP:     ps.StartConnServerIP,
				OSUser: ps.StartConnServerOSUser,
			},
		},
	}

	if adapter.Purpose == types.Change {
		cr, err := s.crSvc.FindByTicketNumber(adapter.ChangeID)
		if err != nil || cr == nil {
			s.log.Error("Failed to find change request by ID", zap.String("changeID", adapter.ChangeID), zap.Error(err))
			return "", "", err
		}

		adapter.CR = cr

		s.log.Debug("Starting CR policy check")
		if err := s.changeRequestPolicyChain.Check(adapter); err != nil {
			s.log.Warn("Change request policy check failed", zap.Error(err))
			return "", "", err
		}
	} else if adapter.Purpose == types.Healthcheck {
		s.log.Debug("Starting healthcheck policy check")
		if err := s.healthCheckPolicyChain.Check(adapter); err != nil {
			s.log.Warn("Health check policy check failed", zap.Error(err))
			return "", "", err
		}
	}

	s.log.Debug("Building connection for join", zap.String("ptySessionId", r.Body.PtySessionId))
	conn := jwt.BuildConnForJoin(ps, r)

	tok, err := s.issuer.Mint(r.AuthCtx, conn, ps.ProxyDetails.ProxyType)
	if err != nil {
		s.log.Error("Failed to mint token", zap.Error(err))
		return "", "", err
	}

	// todo: return X-Proxy-Host ps.ProxyHostName, to be passed to load balancer for routing
	return tok, ps.ProxyDetails.LoadBalancerEndpoint, nil
}

func (s *Service) getChangeRequestIDOrError(p types.ConnectionPurpose, id string) (string, error) {
	switch p {
	case types.Change:
		if id == "" {
			return "", errors.New("missing change request ID for change purpose")
		}
		return id, nil
	case types.Healthcheck:
		return "", nil
	default:
		return "", errors.New("invalid connection purpose: " + string(p))
	}
}
