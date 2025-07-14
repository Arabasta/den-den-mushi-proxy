package pty_token

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/internal/control/proxy_lb"
	"den-den-mushi-Go/internal/control/pty_sessions"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type Service struct {
	issuer *jwt.Issuer

	psS   *pty_sessions.Service
	plbS  *proxy_lb.Service
	hostS *host.Service
	crS   *change_request.Service

	changeRequestPolicyChain policy.Policy[dto.RequestCtx]
	healthCheckPolicyChain   policy.Policy[dto.RequestCtx]

	log *zap.Logger
	cfg *config.Config
}

func NewService(psS *pty_sessions.Service, plbS *proxy_lb.Service, hostS *host.Service,
	issuer *jwt.Issuer, crS *change_request.Service,
	changeRequestPolicyChain policy.Policy[dto.RequestCtx],
	healthCheckPolicyChain policy.Policy[dto.RequestCtx],
	log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		psS:                      psS,
		plbS:                     plbS,
		hostS:                    hostS,
		issuer:                   issuer,
		changeRequestPolicyChain: changeRequestPolicyChain,
		healthCheckPolicyChain:   healthCheckPolicyChain,
		crS:                      crS,
		log:                      log,
		cfg:                      cfg,
	}
}

func (s *Service) mintJoinToken(r *dto.JoinRequest) (string, string, error) {
	ps, err := s.psS.FindById(r.PtySessionId)
	if err != nil {
		s.log.Error("Failed to find pty session", zap.String("ptySessionId", r.PtySessionId), zap.Error(err))
		return "", "", err
	}
	if ps == nil {
		s.log.Error("Pty session not found", zap.String("ptySessionId", r.PtySessionId))
		return "", "", errors.New("pty session not found")
	}

	adapter := &dto.JoinAdapter{
		Req:      r,
		Purpose:  ps.StartConnectionDetails.Purpose,
		ChangeID: ps.StartConnectionDetails.ChangeRequest.Id,
		Server:   ps.StartConnectionDetails.Server,
	}

	hostType, err := s.hostS.FindTypeByIp(adapter.Server.IP)
	if err != nil {
		s.log.Error("Failed to find host type by IP", zap.String("ip", adapter.Server.IP), zap.Error(err))
		return "", "", err
	}

	proxyLbUrl, err := s.plbS.GetLBEndpointByProxyType(hostType)
	if err != nil {
		s.log.Error("Failed to get proxy load balancer URL", zap.String("hostType", string(hostType)), zap.Error(err))
		return "", "", err
	}

	if adapter.Purpose == types.Change {
		if err := s.changeRequestPolicyChain.Check(adapter); err != nil {
			s.log.Error("Change request policy check failed", zap.Error(err))
			return "", "", err
		}
	} else if adapter.Purpose == types.Healthcheck {
		if err := s.healthCheckPolicyChain.Check(adapter); err != nil {
			s.log.Error("Health check policy check failed", zap.Error(err))
			return "", "", err
		}
	}

	conn := jwt.BuildConnForJoin(ps, r)

	tok, err := s.issuer.Mint(r.UserId, conn, hostType)

	return tok, proxyLbUrl, nil
}

func (s *Service) mintStartToken(r *dto.StartRequest) (string, string, error) {
	// hostConnMethod, err := s..s.hostS.FindHostConnectionMethodByIp(r.Server.IP) todo
	hostConnMethod := types.SshTestKey

	hostType, err := s.hostS.FindTypeByIp(r.Server.IP)
	if err != nil {
		s.log.Error("Failed to find host type by IP", zap.String("ip", r.Server.IP), zap.Error(err))
		return "", "", err
	}

	proxyLbUrl, err := s.plbS.GetLBEndpointByProxyType(hostType)
	if err != nil {
		s.log.Error("Failed to get proxy load balancer URL", zap.String("hostType", string(hostType)), zap.Error(err))
		return "", "", err
	}

	var f types.Filter
	var cr *change_request.Entity

	if r.Purpose == types.Change {
		if err := s.changeRequestPolicyChain.Check(r); err != nil {
			s.log.Error("Change request policy check failed", zap.Error(err))
			return "", "", err
		}

		cr, err = s.crS.FindById(r.ChangeID)
		if err != nil {
			s.log.Error("Failed to find change request by ID", zap.String("changeID", r.ChangeID), zap.Error(err))
			return "", "", err
		}
	} else if r.Purpose == types.Healthcheck {
		if err := s.healthCheckPolicyChain.Check(r); err != nil {
			s.log.Error("Health check policy check failed", zap.Error(err))
			return "", "", err
		}

		f, err = s.hostS.FindFilterTypeByHostType(hostType)
		if err != nil {
			s.log.Error("Failed to find filter type by host type", zap.String("hostType", string(hostType)), zap.Error(err))
			return "", "", err
		}
	}

	conn := jwt.BuildConnForStart(hostConnMethod, r, cr, f)

	tok, err := s.issuer.Mint(r.UserId, conn, hostType)
	if err != nil {
		s.log.Error("Failed to mint token", zap.Error(err))
		return "", "", err
	}

	return tok, proxyLbUrl, nil
}
