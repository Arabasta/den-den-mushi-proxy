package pty_token

import (
	request2 "den-den-mushi-Go/internal/control/app/pty_token/request"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/core/certname"
	"den-den-mushi-Go/internal/control/core/change_request"
	"den-den-mushi-Go/internal/control/core/host"
	"den-den-mushi-Go/internal/control/core/iexpress"
	"den-den-mushi-Go/internal/control/core/os_adm_users"
	"den-den-mushi-Go/internal/control/core/proxy_lb"
	"den-den-mushi-Go/internal/control/core/pty_sessions"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	iexpress2 "den-den-mushi-Go/pkg/dto/iexpress"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"den-den-mushi-Go/pkg/util/cyberark"
	"errors"
	"go.uber.org/zap"
)

type Service struct {
	issuer *jwt.Issuer

	psSvc                    *pty_sessions.Service
	plbSvc                   *proxy_lb.Service
	hostSvc                  *host.Service
	crSvc                    *change_request.Service
	certNameSvc              *certname.Service
	changeRequestPolicyChain policy.Policy[request2.Ctx]
	healthCheckPolicyChain   policy.Policy[request2.Ctx]
	iexpressPolicyChain      policy.Policy[request2.Ctx]
	osAdmUsersSvc            *os_adm_users.Service
	iexpressSvc              *iexpress.Service

	log *zap.Logger
	cfg *config.Config
}

func NewService(psS *pty_sessions.Service, plbS *proxy_lb.Service, hostS *host.Service, certNameSvc *certname.Service,
	issuer *jwt.Issuer, crS *change_request.Service, osAdmUsersSvc *os_adm_users.Service, iexpressSvc *iexpress.Service,
	changeRequestPolicyChain policy.Policy[request2.Ctx],
	healthCheckPolicyChain policy.Policy[request2.Ctx],
	iexpressPolicyChain policy.Policy[request2.Ctx],
	log *zap.Logger, cfg *config.Config) *Service {
	log.Info("Initializing PTY Token Service")
	return &Service{
		psSvc:                    psS,
		plbSvc:                   plbS,
		hostSvc:                  hostS,
		certNameSvc:              certNameSvc,
		issuer:                   issuer,
		changeRequestPolicyChain: changeRequestPolicyChain,
		healthCheckPolicyChain:   healthCheckPolicyChain,
		iexpressPolicyChain:      iexpressPolicyChain,
		crSvc:                    crS,
		osAdmUsersSvc:            osAdmUsersSvc,
		iexpressSvc:              iexpressSvc,
		log:                      log,
		cfg:                      cfg,
	}
}

// todo: split into healthcheck and cr endpoints
func (s *Service) mintStartToken(r wrapper.WithAuth[request2.StartRequest]) (string, string, error) {
	hostConnMethod, hostType := types.LocalSshKey, types.OS

	if !s.cfg.Development.IsLocalSshKeyIfNotIsPuppetKey {
		hostConnMethod = types.SshOrchestratorKey
	}

	var err error

	// find hostname by ip and then certname in puppettrusted what the heck
	puppetTrusted, err := s.certNameSvc.FindCertnameByIp(r.Body.Server.IP)
	if err != nil || puppetTrusted == nil || puppetTrusted.Certname == "" {
		s.log.Error("Failed to find host certname by IP", zap.String("ip", r.Body.Server.IP), zap.Error(err))
		return "", "", errors.New("failed to find host certname")
	}

	// todo: enable all these
	// hostConnMethod, err := s.hostS.FindHostConnectionMethodByIp(r.Server.IP) todo: grab server conn method how?

	//hostType, err := s.hostSvc.FindTypeByIp(r.Body.Server.IP)
	//if err != nil {
	//	s.log.Error("Failed to find host type by IP", zap.String("ip", r.Body.Server.IP), zap.Error(err))
	//	return "", "", err
	//}

	var filter types.Filter
	var cr *changerequestpkg.Record
	var exp *iexpress2.Record

	adapter := &request2.StartAdapter{
		Req: r,
		AdapterFields: request2.AdapterFields{
			Purpose:  r.Body.Purpose,
			ChangeID: r.Body.ChangeID,
			Server: dto.ServerInfo{
				IP:     r.Body.Server.IP,
				OSUser: r.Body.Server.OSUser,
			},
		},
	}

	var allowedSuOsUsers []string

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

		// todo update this need to su? or just root only?
		allowedSuOsUsers = cyberark.ExtractAllOsUsers(cr.CyberArkObjects)
		s.log.Debug("CR allowed OS users extracted", zap.Strings("allowedSuOsUsers", allowedSuOsUsers))
	} else if r.Body.Purpose == types.Healthcheck {
		s.log.Debug("Starting healthcheck policy check")
		if err = s.healthCheckPolicyChain.Check(adapter); err != nil {
			s.log.Warn("Health check policy check failed", zap.Error(err))
			return "", "", err
		}

		if s.cfg.Development.IsBlacklistFilter {
			filter = types.Blacklist // todo: get filter type by host type or OU group?
		} else {
			filter = types.Whitelist
		}

		// todo: this guy should be from db, each perso n has their own os user
		// todo: maybe dont need seems liek su for root only
		//	allowedSuOsUsers = s.osAdmUsersSvc.GetNonCrOsUsers(r.AuthCtx.UserID)
		s.log.Debug("Healthcheck allowed OS users", zap.Strings("allowedSuOsUsers", allowedSuOsUsers))

		// todo: this should be by OU group what they want idk
		//filter, err = s.hostSvc.FindFilterTypeByHostType(hostType)
		//if err != nil {
		//	s.log.Error("Failed to find filter type by host type", zap.String("hostType", string(hostType)), zap.Error(err))
		//	return "", "", err
		//}
	} else if r.Body.Purpose == types.IExpress {
		s.log.Debug("Starting IExpress policy check")
		exp, err = s.iexpressSvc.FindByTicketNumber(r.Body.ChangeID)
		if err != nil || exp == nil {
			s.log.Error("Failed to find iexpress request by ID", zap.String("iexpress", r.Body.ChangeID), zap.Error(err))
			return "", "", err
		}

		adapter.Iexpress = exp

		s.log.Debug("Starting IExpress policy check")
		if err = s.iexpressPolicyChain.Check(adapter); err != nil {
			s.log.Warn("IExpress policy check failed", zap.Error(err))
			return "", "", err
		}

		// todo update this need to su? or just root only?
		allowedSuOsUsers = cyberark.ExtractAllOsUsers(exp.CyberArkObjects)
		s.log.Debug("IExpress allowed OS users extracted", zap.Strings("allowedSuOsUsers", allowedSuOsUsers))

		if s.cfg.Development.IsBlacklistFilter {
			filter = types.Blacklist // todo: get filter type by host type or OU group?
		} else {
			filter = types.Whitelist
		}
	} else {
		s.log.Error("Invalid connection purpose", zap.String("purpose", string(r.Body.Purpose)))
		return "", "", errors.New("invalid connection purpose")
	}

	s.log.Debug("Building connection for start")
	conn := jwt.BuildConnForStart(hostConnMethod, r, cr, exp, filter, s.cfg.Development.TargetSshPort, allowedSuOsUsers, puppetTrusted.Certname)

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

func (s *Service) mintJoinToken(r wrapper.WithAuth[request2.JoinRequest]) (string, string, error) {
	ps, err := s.psSvc.FindById(r.Body.PtySessionId)
	if err != nil || ps == nil {
		s.log.Error("Failed to find pty session", zap.String("ptySessionId", r.Body.PtySessionId), zap.Error(err))
		return "", "", errors.New("failed to find pty session")
	}

	ticketId, err := s.getTicketIDOrError(ps.StartConnPurpose, ps.StartConnChangeRequestID)
	if err != nil {
		s.log.Error("Invalid connection details", zap.String("ptySessionId", r.Body.PtySessionId), zap.Error(err))
		return "", "", err
	}

	adapter := &request2.JoinAdapter{
		Req: r,
		AdapterFields: request2.AdapterFields{
			Purpose:  ps.StartConnPurpose,
			ChangeID: ticketId,
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
	} else if adapter.Purpose == types.IExpress {
		exp, err := s.iexpressSvc.FindByTicketNumber(adapter.ChangeID)
		if err != nil || exp == nil {
			s.log.Error("Failed to find IExpress request by ID", zap.String("IExpress", adapter.ChangeID), zap.Error(err))
			return "", "", err
		}

		adapter.Iexpress = exp

		s.log.Debug("Starting IExpress policy check")
		if err := s.iexpressPolicyChain.Check(adapter); err != nil {
			s.log.Warn("IExpress requesst policy check failed", zap.Error(err))
			return "", "", err
		}
	} else {
		s.log.Error("Invalid connection purpose", zap.String("purpose", string(adapter.Purpose)))
		return "", "", errors.New("invalid connection purpose")
	}

	s.log.Debug("Building connection for join", zap.String("ptySessionId", r.Body.PtySessionId))
	conn := jwt.BuildConnForJoin(ps, r)

	// todo: get this from start conn
	proxyType := types.OS

	tok, err := s.issuer.Mint(r.AuthCtx, conn, proxyType)
	if err != nil {
		s.log.Error("Failed to mint token", zap.Error(err))
		return "", "", err
	}

	// return X-Proxy-Host ps.ProxyHostName, to be passed to load balancer for routing?
	// or maybe just straight up manual routing
	return tok, ps.ProxyHostName, nil
}

func (s *Service) getTicketIDOrError(p types.ConnectionPurpose, id string) (string, error) {
	switch p {
	case types.Change:
		if id == "" {
			return "", errors.New("missing change request ID for change purpose")
		}
		return id, nil
	case types.Healthcheck:
		return "", nil
	case types.IExpress:
		if id == "" {
			return "", errors.New("missing IExpress request ID for IExpress purpose")
		}
		return id, nil
	default:
		return "", errors.New("invalid connection purpose: " + string(p))
	}
}
