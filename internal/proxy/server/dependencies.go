package server

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/connect"
	"den-den-mushi-Go/internal/proxy/control_server_tmp"
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"den-den-mushi-Go/internal/proxy/jwt_service"
	"den-den-mushi-Go/internal/proxy/jwt_service/jti"
	"den-den-mushi-Go/internal/proxy/orchestrator/puppet"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/internal/proxy/websocket"
	"go.uber.org/zap"
)

type Deps struct {
	WebsocketService *websocket.Service
	Validator        *jwt_service.Validator
	Issuer           *control_server_tmp.Issuer // todo: tmp move this to control server
	SessionManager   *session_manager.Service   // todo: tmp remove this
}

func initDependencies(cfg *config.Config, log *zap.Logger) *Deps {
	connectionMethodFactory := connect.NewConnectionMethodFactory(
		connect.NewDeps(
			puppet.NewClient(cfg, log),
			pty_util.NewBuilder(log),
			cfg,
			log))

	sessionManager := session_manager.New(log, cfg)
	websocketService := websocket.NewService(connectionMethodFactory, sessionManager, log, cfg)

	issuer := control_server_tmp.New(cfg, log)

	parser := jwt_service.NewParser(cfg, log)

	var jtiRepo jti.Repository

	if cfg.App.Environment == "dev" && cfg.Development.UseInMemoryRepository {
		jtiRepo = jti.NewInMemRepository(log)
	} else {
		log.Fatal("JTI repository not implemented for this environment")
	}

	jtiService := jti.New(jtiRepo, log)

	val := jwt_service.NewValidator(parser, jtiService, cfg.Token.Secret, cfg, log)

	return &Deps{
		WebsocketService: websocketService,
		Validator:        val,
		Issuer:           issuer,
		SessionManager:   sessionManager,
	}
}
