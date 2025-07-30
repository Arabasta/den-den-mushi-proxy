package server

import (
	"den-den-mushi-Go/internal/llm_external/config"
	"den-den-mushi-Go/internal/llm_external/pty_sessions"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Deps struct {
	PtySessionService *pty_sessions.Service
}

func initDependencies(ddmDb *gorm.DB, cfg *config.Config, log *zap.Logger) *Deps {
	ptySessionRepo := pty_sessions.NewGormRepository(ddmDb, log)
	ptySessionService := pty_sessions.NewService(ptySessionRepo, log)

	return &Deps{
		PtySessionService: ptySessionService,
	}
}
