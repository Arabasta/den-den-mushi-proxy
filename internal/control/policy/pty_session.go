package policy

import (
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/pty_sessions"
	"go.uber.org/zap"
)

type PtySessionPolicy[T dto.RequestCtx] struct {
	next Policy[T]

	ptySessionService *pty_sessions.Service
	log               *zap.Logger
}

func NewPtySessionPolicy[T dto.RequestCtx](ptySessionService *pty_sessions.Service, log *zap.Logger) *PtySessionPolicy[T] {
	return &PtySessionPolicy[T]{
		ptySessionService: ptySessionService,
		log:               log,
	}
}

func (p *PtySessionPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *PtySessionPolicy[T]) Check(req T) error {
	// 1. check server max connections

	// 2. check if have active implementor (if joining as implementor)

	return nil
}
