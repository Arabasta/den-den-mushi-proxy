package pty_token

import (
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

func (s *Service) mintStartTokenHealth(r wrapper.WithAuth[request.StartRequest]) (string, string, error) {
	adapter := &request.StartAdapter{Req: r}

	s.log.Debug("Starting healthcheck policy check")
	if err := s.healthCheckPolicyChain.Check(adapter); err != nil {
		s.log.Warn("Health check policy check failed", zap.Error(err))
		return "", "", err
	}

	filter := types.Blacklist // todo: get filter type by host type
	//filter, err = s.hostSvc.FindFilterTypeByHostType(hostType)
	//if err != nil {
	//	s.log.Error("Failed to find filter type by host type", zap.String("hostType", string(hostType)), zap.Error(err))
	//	return "", "", err
	//}

	hostConnMethod, hostType := types.SshTestKey, types.OS
	// todo: enable all these
	// hostConnMethod, err := s.hostS.FindHostConnectionMethodByIp(r.Server.IP) todo: grab server conn method how?

	//hostType, err := s.hostSvc.FindTypeByIp(r.Body.Server.IP)
	//if err != nil {
	//	s.log.Error("Failed to find host type by IP", zap.String("ip", r.Body.Server.IP), zap.Error(err))
	//	return "", "", err
	//}

	s.log.Debug("Building connection for start")
	conn := jwt.BuildConnForStart(hostConnMethod, r, nil, filter, s.cfg.Development.TargetSshPort)

	return s.mintTokenAndGetProxyUrl(r.AuthCtx, conn, hostType)
}

func (s *Service) mintStartTokenCR(r wrapper.WithAuth[request.StartRequest]) (string, string, error) {
	adapter := &request.StartAdapter{Req: r}

	cr, err := s.checkCRAndAdapt(adapter, r.Body.ChangeID)
	if err != nil {
		return "", "", err
	}

	hostConnMethod, hostType := types.SshTestKey, types.OS

	// todo: enable all these
	// hostConnMethod, err := s.hostS.FindHostConnectionMethodByIp(r.Server.IP) todo: grab server conn method how?
	//hostType, err := s.hostSvc.FindTypeByIp(r.Body.Server.IP)
	//if err != nil {
	//	s.log.Error("Failed to find host type by IP", zap.String("ip", r.Body.Server.IP), zap.Error(err))
	//	return "", "", err
	//}

	conn := jwt.BuildConnForStart(hostConnMethod, r, cr, "", s.cfg.Development.TargetSshPort)
	return s.mintTokenAndGetProxyUrl(r.AuthCtx, conn, hostType)
}

func (s *Service) mintTokenAndGetProxyUrl(authCtx middleware.AuthContext, conn *dto.Connection, proxyType types.Proxy) (string, string, error) {
	token, err := s.issuer.Mint(authCtx, conn, proxyType)
	if err != nil {
		s.log.Error("Failed to mint token", zap.Error(err))
		return "", "", err
	}

	url, err := s.plbSvc.GetLBEndpointByProxyType(proxyType)
	if err != nil {
		s.log.Error("Failed to get proxy load balancer URL", zap.String("hostType", string(proxyType)), zap.Error(err))
		return "", "", err
	}

	return token, url, nil
}

func (s *Service) checkCRAndAdapt(adapter *request.StartAdapter, changeID string) (*changerequestpkg.Record, error) {
	cr, err := s.crSvc.FindByTicketNumber(changeID)
	if err != nil || cr == nil {
		s.log.Error("Failed to find change request", zap.String("changeID", changeID), zap.Error(err))
		return nil, err
	}
	adapter.CR = cr

	s.log.Debug("Starting CR policy check")
	if err := s.changeRequestPolicyChain.Check(adapter); err != nil {
		s.log.Warn("CR policy check failed", zap.Error(err))
		return nil, err
	}
	return cr, nil
}
