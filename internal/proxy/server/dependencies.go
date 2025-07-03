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
	"time"
)

type Deps struct {
	WebsocketService *websocket.Service
	Validator        *jwt_service.Validator
	Issuer           *control_server_tmp.Issuer // todo: tmp move this to control server
}

func initDependencies(cfg *config.Config, log *zap.Logger) *Deps {
	connectionMethodFactory := connect.NewConnectionMethodFactory(
		connect.NewDeps(
			puppet.NewClient(cfg, log),
			pty_util.NewBuilder(log),
			cfg,
			log))

	sessionManager := session_manager.New(log)

	websocketService := websocket.NewWebsocketService(connectionMethodFactory, sessionManager, log, cfg)

	ttl := 60 * time.Second // todo: tmp here for now

	issuer := control_server_tmp.NewIssuer(cfg.Token.Secret, cfg.Token.Issuer, cfg.Token.Audience, ttl)

	parser := jwt_service.NewParser(cfg, log)

	jtiInMemRepo := jti.NewInMemRepository(log)
	jtiService := jti.New(jtiInMemRepo, ttl, log)

	val := jwt_service.NewValidator(parser, jtiService, cfg.Token.Secret, cfg, log)

	return &Deps{
		WebsocketService: websocketService,
		Validator:        val,
		Issuer:           issuer,
	}
}
