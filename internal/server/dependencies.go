package server

import (
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/control_server_tmp"
	"den-den-mushi-Go/internal/orchestrator/puppet"
	"den-den-mushi-Go/internal/pseudo/connect"
	"den-den-mushi-Go/internal/validator"
	"den-den-mushi-Go/internal/websocket"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type Deps struct {
	PuppetClient     *puppet.PuppetClient
	WebsocketService *websocket.Service
	Validator        *validator.Validator
	Issuer           *control_server_tmp.Issuer // todo: tmp move this to control server
}

func initDependencies(cfg *config.Config, log *zap.Logger) *Deps {
	puppetClient := puppet.NewPuppetClient(cfg, log)

	connectionMethodFactory := connect.NewConnectionMethodFactory(connect.Deps{Puppet: puppetClient, Cfg: cfg, Log: log})
	websocketService := websocket.NewWebsocketService(connectionMethodFactory)

	secret, iss, aud := cfg.Token.Secret, "control", "proxy"
	ttl := 60 * time.Second
	issuer := control_server_tmp.NewIssuer(secret, iss, aud, ttl)
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience(aud),
		jwt.WithIssuer(iss),
	)

	val := validator.NewValidator(parser, cfg.Token.Secret, ttl)

	return &Deps{
		PuppetClient:     puppetClient,
		WebsocketService: websocketService,
		Validator:        val,
		Issuer:           issuer,
	}
}
