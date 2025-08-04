package pty_sessions

import (
	"go.uber.org/zap"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}
