package config

import (
	"den-den-mushi-Go/pkg/config"
)

type Config struct {
	App *config.App

	Development struct {
		IsUsingInvDb bool
		IsSMX        bool
	}

	Ssl *config.Ssl

	Cors *config.Cors

	Logger *config.Logger

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Api *config.Api
}
