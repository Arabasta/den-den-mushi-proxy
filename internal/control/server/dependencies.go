package server

import (
	"den-den-mushi-Go/internal/control/app/healthcheck"
	iexpress2 "den-den-mushi-Go/internal/control/app/iexpress"
	"den-den-mushi-Go/internal/control/app/make_change"
	"den-den-mushi-Go/internal/control/app/pty_token"
	"den-den-mushi-Go/internal/control/app/pty_token/request"
	"den-den-mushi-Go/internal/control/app/whiteblacklist"
	"den-den-mushi-Go/internal/control/config"
	certname2 "den-den-mushi-Go/internal/control/core/certname"
	change_request2 "den-den-mushi-Go/internal/control/core/change_request"
	connection2 "den-den-mushi-Go/internal/control/core/connection"
	host2 "den-den-mushi-Go/internal/control/core/host"
	iexpress3 "den-den-mushi-Go/internal/control/core/iexpress"
	implementor_groups2 "den-den-mushi-Go/internal/control/core/implementor_groups"
	os_adm_users2 "den-den-mushi-Go/internal/control/core/os_adm_users"
	proxy_lb2 "den-den-mushi-Go/internal/control/core/proxy_lb"
	pty_sessions2 "den-den-mushi-Go/internal/control/core/pty_sessions"
	regex_filters2 "den-den-mushi-Go/internal/control/core/regex_filters"
	"den-den-mushi-Go/internal/control/jwt"
	"den-den-mushi-Go/internal/control/policy"
	"den-den-mushi-Go/internal/control/policy/validators"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Deps struct {
	Issuer                   *jwt.Issuer
	ProxyService             *proxy_lb2.Service
	ChangeService            *change_request2.Service
	ImplementorGroupsService *implementor_groups2.Service
	PtySessionService        *pty_sessions2.Service
	HostService              *host2.Service
	PtyTokenService          *pty_token.Service
	MakeChangeService        *make_change.Service
	RegexService             *regex_filters2.Service
	WhiteBlacklistService    *whiteblacklist.Service
	HealthcheckService       *healthcheck.Service
	IexpressService          *iexpress2.Service
}

func initDependencies(ddmDb *gorm.DB, cfg *config.Config, log *zap.Logger) *Deps {
	issuer := jwt.New(cfg, log)

	// repos and services ========================================================================================
	proxyLbRepo := proxy_lb2.NewGormRepository(ddmDb, log)
	proxyLbService := proxy_lb2.NewService(proxyLbRepo, log)

	//proxyHostRepo := proxy_host.NewGormRepository(ddmDb, log)
	//proxyHostService := proxy_host.NewService(proxyHostRepo, log)

	changeRepo := change_request2.NewGormRepository(ddmDb, log)
	changeService := change_request2.NewService(changeRepo, log)

	impGroupsRepo := implementor_groups2.NewGormRepository(ddmDb, log, cfg)
	impGroupsService := implementor_groups2.NewService(impGroupsRepo, log)

	ptySessionRepo := pty_sessions2.NewGormRepository(ddmDb, log)
	ptySessionService := pty_sessions2.NewService(ptySessionRepo, log)

	hostRepo := host2.NewGormRepository(ddmDb, log)
	hostService := host2.NewService(hostRepo, log)

	regexRepo := regex_filters2.NewGormRepository(ddmDb, log)
	regexSvc := regex_filters2.NewService(regexRepo, log)

	whiteblacklistSvc := whiteblacklist.NewService(regexSvc, log)

	connRepo := connection2.NewGormRepository(ddmDb, log)
	connectionService := connection2.NewService(connRepo, log)

	certNameRepo := certname2.NewGormRepository(ddmDb, log)
	certNameSvc := certname2.NewService(certNameRepo, log, cfg)

	osAdmUsersRepo := os_adm_users2.NewGormRepository(ddmDb, log)
	osAdmUsersService := os_adm_users2.NewService(osAdmUsersRepo, log)

	iexpressRepo := iexpress3.NewGormRepository(ddmDb, log)
	iexpressService := iexpress3.NewService(iexpressRepo, log)

	// validator for policy chains  ============================================================================================
	validator := validators.NewValidator(log, cfg)

	//// policy chains ============================================================================================
	var changeRequestPolicyChain policy.Policy[request.Ctx]
	var healthcheckPolicyChain policy.Policy[request.Ctx]

	// policies for change request
	ptySessionPolicyCR := policy.NewPtySessionPolicy[request.Ctx](ptySessionService, connectionService, log)
	implementorPolicyCR := policy.NewImplementorPolicy[request.Ctx](hostService, log, validator)
	changePolicy := policy.NewChangePolicy[request.Ctx](impGroupsService, osAdmUsersService, validator, log)

	ptySessionPolicyCR.SetNext(implementorPolicyCR)
	implementorPolicyCR.SetNext(changePolicy)
	changeRequestPolicyChain = ptySessionPolicyCR

	// policies for health check
	ptySessionPolicyHC := policy.NewPtySessionPolicy[request.Ctx](ptySessionService, connectionService, log)
	implementorPolicyHC := policy.NewImplementorPolicy[request.Ctx](hostService, log, validator)
	healthcheckPolicy := policy.NewHealthcheckPolicy[request.Ctx](hostService, impGroupsService, osAdmUsersService, log, validator)

	ptySessionPolicyHC.SetNext(healthcheckPolicy)
	healthcheckPolicy.SetNext(implementorPolicyHC)
	healthcheckPolicyChain = ptySessionPolicyHC

	// policies for iexpress
	ptySessionPolicyIExpress := policy.NewPtySessionPolicy[request.Ctx](ptySessionService, connectionService, log)
	implementorPolicyIExpress := policy.NewImplementorPolicy[request.Ctx](hostService, log, validator)
	iexpressPolicy := policy.NewIExpressPolicy[request.Ctx](impGroupsService, osAdmUsersService, validator, log)

	ptySessionPolicyIExpress.SetNext(iexpressPolicy)
	iexpressPolicy.SetNext(implementorPolicyIExpress)
	iexpressPolicyChain := ptySessionPolicyIExpress

	if cfg.Development.SkipPolicyChecks {
		log.Info("Skip policy checks enabled. Using noop policy")
		noopPolicy := policy.NewNoopPolicy[request.Ctx](log)
		changeRequestPolicyChain = noopPolicy
		healthcheckPolicyChain = noopPolicy
	}

	// pass policy chains to service
	ptyTokenService := pty_token.NewService(ptySessionService, proxyLbService, hostService, certNameSvc, issuer, changeService,
		osAdmUsersService, iexpressService, changeRequestPolicyChain, healthcheckPolicyChain, iexpressPolicyChain,
		log, cfg)

	makeChangeService := make_change.NewService(changeService, ptySessionService, hostService, impGroupsService, osAdmUsersService, log, cfg)
	healthcheckService := healthcheck.NewService(ptySessionService, hostService, osAdmUsersService, log, cfg)
	iexpresssvcEp := iexpress2.NewService(iexpressService, ptySessionService, hostService, impGroupsService, osAdmUsersService, log)
	return &Deps{
		Issuer:                   issuer,
		ProxyService:             proxyLbService,
		ChangeService:            changeService,
		ImplementorGroupsService: impGroupsService,
		PtySessionService:        ptySessionService,
		HostService:              hostService,
		PtyTokenService:          ptyTokenService,
		MakeChangeService:        makeChangeService,
		RegexService:             regexSvc,
		WhiteBlacklistService:    whiteblacklistSvc,
		HealthcheckService:       healthcheckService,
		IexpressService:          iexpresssvcEp,
	}
}
