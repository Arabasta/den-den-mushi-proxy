package server

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/internal/control/proxy_lb"
	"den-den-mushi-Go/internal/control/pty_sessions"
	"den-den-mushi-Go/internal/control/pty_token"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func initDependencies(ddmDb *gorm.DB, cfg *config.Config, log *zap.Logger) *Deps {
	issuer := jwt.New(cfg, log)

	// repos and services ========================================================================================
	proxyLbRepo := proxy_lb.NewGormRepository(ddmDb, log)
	proxyLbService := proxy_lb.NewService(proxyLbRepo, log)

	//proxyHostRepo := proxy_host.NewGormRepository(ddmDb, log)
	//proxyHostService := proxy_host.NewService(proxyHostRepo, log)

	changeRepo := change_request.NewGormRepository(ddmDb, log)
	changeService := change_request.NewService(changeRepo, log)

	impGroupsRepo := implementor_groups.NewGormRepository(ddmDb, log)
	impGroupsService := implementor_groups.NewService(impGroupsRepo, log)

	ptySessionRepo := pty_sessions.NewGormRepository(ddmDb, log)
	ptySessionService := pty_sessions.NewService(ptySessionRepo, log)

	hostRepo := host.NewGormRepository(ddmDb, log)
	hostService := host.NewService(hostRepo, log)

	// policies ==================================================================================================

	changePolicy := policy.NewChangePolicy[request.Ctx](impGroupsService, log)
	healthcheckPolicy := policy.NewHealthcheckPolicy[request.Ctx](hostService, impGroupsService, log)
	ouPolicy := policy.NewOUPolicy[request.Ctx](hostService, log)
	ptySessionPolicy := policy.NewPtySessionPolicy[request.Ctx](ptySessionService, log)

	// policy chains ============================================================================================

	changeRequestPolicyChain := policy.NewBuilder[request.Ctx]().
		Add(changePolicy).
		Add(ouPolicy).
		Add(ptySessionPolicy).
		Build()

	healthcheckPolicyChain := policy.NewBuilder[request.Ctx]().
		Add(healthcheckPolicy).
		Add(ouPolicy).
		Add(ptySessionPolicy).
		Build()

	if cfg.App.Environment == "dev" && cfg.Development.SkipPolicyChecks {
		log.Info("Using noop policy in development mode")
		changeRequestPolicyChain = policy.NewBuilder[request.Ctx]().
			Add(policy.NewNoopPolicy[request.Ctx](log)).Build()
		healthcheckPolicyChain = policy.NewBuilder[request.Ctx]().
			Add(policy.NewNoopPolicy[request.Ctx](log)).Build()
	}

	// pass policy chains to service
	ptyTokenService := pty_token.NewService(ptySessionService, proxyLbService, hostService, issuer, changeService,
		changeRequestPolicyChain, healthcheckPolicyChain,
		log, cfg)

	return &Deps{
		Issuer:                   issuer,
		ProxyService:             proxyLbService,
		ChangeService:            changeService,
		ImplementorGroupsService: impGroupsService,
		PtySessionService:        ptySessionService,
		HostService:              hostService,
		PtyTokenService:          ptyTokenService,
	}
}
