package policy

import (
	"den-den-mushi-Go/internal/control/pty_sessions"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type PtySessionPolicy[T request.Ctx] struct {
	next Policy[T]

	ptySessionService *pty_sessions.Service
	log               *zap.Logger
}

func NewPtySessionPolicy[T request.Ctx](ptySessionService *pty_sessions.Service, log *zap.Logger) *PtySessionPolicy[T] {
	return &PtySessionPolicy[T]{
		ptySessionService: ptySessionService,
		log:               log,
	}
}

func (p *PtySessionPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *PtySessionPolicy[T]) Check(r T) error {
	sessionAware, isJoin := any(r).(request.HasJoinRequestFields)
	if !isJoin {
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// 1. check if pty session is active
	//ptyID := sessionAware.GetPtySessionId()

	if sessionAware.GetStartRole() == types.Implementor {
		// 2. check if pty session as active implementor

	}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
