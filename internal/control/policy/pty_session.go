package policy

import (
	"den-den-mushi-Go/internal/control/app/pty_token/request"
	"den-den-mushi-Go/internal/control/core/connection"
	"den-den-mushi-Go/internal/control/core/pty_sessions"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

type PtySessionPolicy[T request.Ctx] struct {
	next Policy[T]

	ptySessionService *pty_sessions.Service
	connectionService *connection.Service
	log               *zap.Logger
}

func NewPtySessionPolicy[T request.Ctx](ptySessionService *pty_sessions.Service, connectionSvc *connection.Service,
	log *zap.Logger) *PtySessionPolicy[T] {
	return &PtySessionPolicy[T]{
		ptySessionService: ptySessionService,
		connectionService: connectionSvc,
		log:               log,
	}
}

func (p *PtySessionPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *PtySessionPolicy[T]) Check(r T) error {
	p.log.Debug("Checking PtySession Policy...")
	sessionAware, isJoin := any(r).(request.HasJoinRequestFields)
	if !isJoin {
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	ptyID := sessionAware.GetPtySessionId()

	// todo move validation logic to validators package

	// check if pty session exists
	p.log.Debug("Checking if pty session exists", zap.String("ptyId", ptyID))
	ptySession, err := p.ptySessionService.FindById(ptyID)
	if err != nil || ptySession == nil {
		p.log.Warn("Failed to find pty session", zap.String("ptyId", ptyID), zap.Error(err))
		return errors.New("pty session not found")
	}

	// check if pty session is active
	p.log.Debug("Checking if pty session is active", zap.String("ptyId", ptyID), zap.String("state", string(ptySession.State)))
	if ptySession.State != types.Active {
		p.log.Warn("Pty session is not active", zap.String("ptyId", ptyID), zap.String("state", string(ptySession.State)))
		return errors.New("pty session is not active")
	}

	// check if pty session has an active implementor if joining as implementor
	if sessionAware.GetStartRole() == types.Implementor {
		p.log.Debug("Starting as Implementor... Checking if pty session has an active implementor", zap.String("ptyId", ptyID))
		activeImplementor, err := p.connectionService.FindActiveImplementorByPtySessionId(ptyID)
		if err != nil {
			p.log.Warn("Error when finding active implementor for pty session", zap.String("ptyId", ptyID), zap.Error(err))
			return errors.New("error when finding active implementor for pty session")
		}
		if activeImplementor != nil {
			p.log.Warn("Pty session already has an active implementor", zap.String("ptyId", ptyID))
			return errors.New("pty session already has an active implementor")
		}
	}

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
