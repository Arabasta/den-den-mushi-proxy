package server

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/internal/control/proxy_lb"
	"den-den-mushi-Go/internal/control/pty_sessions"
	"den-den-mushi-Go/internal/control/pty_token"
	"go.uber.org/zap"
)

type Deps struct {
	Issuer                   *jwt.Issuer
	ProxyService             *proxy_lb.Service
	ChangeService            *change_request.Service
	ImplementorGroupsService *implementor_groups.Service
	PtySessionService        *pty_sessions.Service
	HostService              *host.Service
	PtyTokenService          *pty_token.Service
}

func initDependencies(cfg *config.Config, log *zap.Logger) *Deps {
	issuer := jwt.New(cfg, log)

	// init repos and services ========================================================================================

	proxyRepo := proxy_lb.NewInMemRepository()
	proxyService := proxy_lb.NewService(proxyRepo, log)

	changeRepo := change_request.NewInMemRepository()
	changeService := change_request.NewService(changeRepo, log)

	impGroupsRepo := implementor_groups.NewInMemRepository()
	impGroupsService := implementor_groups.NewService(impGroupsRepo, log)

	ptySessionRepo := pty_sessions.NewInMemRepository()
	ptySessionService := pty_sessions.NewService(ptySessionRepo, log)

	hostRepo := host.NewInMemRepository()
	hostService := host.NewService(hostRepo, log)

	// init policies ==================================================================================================

	changePolicy := policy.NewChangePolicy[dto.RequestCtx](changeService, impGroupsService, hostService, log)
	healthcheckPolicy := policy.NewHealthcheckPolicy[dto.RequestCtx](hostService, impGroupsService, log)
	ouPolicy := policy.NewOUPolicy[dto.RequestCtx](hostService, log)
	ptySessionPolicy := policy.NewPtySessionPolicy[dto.RequestCtx](ptySessionService, log)

	// build policy chains ============================================================================================

	changeRequestPolicyChain := policy.NewPolicyBuilder[dto.RequestCtx]().
		Add(changePolicy).
		Add(ouPolicy).
		Add(ptySessionPolicy).
		Build()

	healthcheckPolicyChain := policy.NewPolicyBuilder[dto.RequestCtx]().
		Add(healthcheckPolicy).
		Add(ouPolicy).
		Add(ptySessionPolicy).
		Build()

	// done
	ptyTokenService := pty_token.NewService(ptySessionService, proxyService, hostService, issuer, changeService,
		changeRequestPolicyChain, healthcheckPolicyChain,
		log, cfg)

	return &Deps{
		Issuer:                   issuer,
		ProxyService:             proxyService,
		ChangeService:            changeService,
		ImplementorGroupsService: impGroupsService,
		PtySessionService:        ptySessionService,
		HostService:              hostService,
		PtyTokenService:          ptyTokenService,
	}
}
