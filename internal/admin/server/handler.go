package server

import (
	"den-den-mushi-Go/internal/llm_external/pty_sessions"
)

type MasterHandler struct {
	PtySessions *pty_sessions.Handler
}

// Forwarding methods (required by oapi.ServerInterface)
