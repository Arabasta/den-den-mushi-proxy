package server

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/connect"
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"den-den-mushi-Go/internal/proxy/core/session_manager/connections"
	"den-den-mushi-Go/internal/proxy/core/session_manager/pty_sessions"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/internal/proxy/regex_filters"
	"den-den-mushi-Go/internal/proxy/tmp/control_server_tmp"
	"den-den-mushi-Go/internal/proxy/websocket"
	"den-den-mushi-Go/internal/proxy/websocket_jwt"
	"den-den-mushi-Go/internal/proxy/websocket_jwt/jti"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Deps struct {
	WebsocketService *websocket.Service
	Validator        *websocket_jwt.Validator
	Issuer           *control_server_tmp.Issuer // todo: tmp move this to control server
	SessionManager   *session_manager.Service   // todo: tmp remove this
}

func initDependencies(db *gorm.DB, redis *redis.Client, cfg *config.Config, log *zap.Logger) *Deps {
	puppetClient := puppet.NewClient(cfg.Puppet, cfg, log)
	connMethodStrategy := connect.NewRegistry(
		connect.NewDeps(
			puppetClient,
			pty_util.NewBuilder(log, cfg.Ssh),
			cfg,
			log))

	ptySessionsRepo := pty_sessions.NewGormRepository(db, log)
	ptySessionsSvc := pty_sessions.NewService(ptySessionsRepo, log)

	connectionRepo := connections.NewGormRepository(db, log)
	connectionSvc := connections.NewService(connectionRepo, log)

	issuer := control_server_tmp.New(cfg.JwtAudience, log) // todo: remove

	parser := websocket_jwt.NewParser(cfg.JwtAudience, log)

	var jtiRepo jti.Repository

	if cfg.Development.UseSqlJtiRepo {
		jtiRepo = jti.NewGormRepository(db, log)
	} else {
		jtiRepo = jti.NewRedisRepository(redis, log)
	}
	jtiService := jti.New(jtiRepo, log, cfg.JwtAudience)

	val := websocket_jwt.NewValidator(parser, jtiService, cfg.JwtAudience, log)

	regexRepo := regex_filters.NewGormRepository(db, log)
	regexFiltersSvc := regex_filters.NewService(regexRepo, log, cfg)
	filterSvc := filter.NewService(regexFiltersSvc, log, cfg)

	sessionManager := session_manager.New(ptySessionsSvc, connectionSvc, log, cfg, puppetClient, filterSvc)
	websocketService := websocket.NewService(connMethodStrategy, sessionManager, log, cfg)

	return &Deps{
		WebsocketService: websocketService,
		Validator:        val,
		Issuer:           issuer,         // todo: remove
		SessionManager:   sessionManager, // todo: remove
	}
}
