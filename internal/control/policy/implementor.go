package policy

import (
	"den-den-mushi-Go/internal/control/ep/pty_token/request"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type ImplementorPolicy[T request.Ctx] struct {
	next Policy[T]

	hostService *host.Service
	log         *zap.Logger
	v           *validators.Validator
}

func NewImplementorPolicy[T request.Ctx](hostService *host.Service, log *zap.Logger, v *validators.Validator) *ImplementorPolicy[T] {
	return &ImplementorPolicy[T]{
		hostService: hostService,
		log:         log,
		v:           v,
	}
}

func (p *ImplementorPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *ImplementorPolicy[T]) Check(r T) error {
	p.log.Debug("Checking implementor policy...")

	if r.GetStartRole() != types.Implementor {
		p.log.Debug("Request is not for implementor role, skipping implementor policy check")
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	p.log.Debug("Request is for implementor role, checking OU Level")
	// user must be in L1 OU group
	if err := p.v.IsL1OuGroup(r.GetUserOuGroup()); err != nil {
		p.log.Warn("User OU group validation failed", zap.String("userId", r.GetUserId()), zap.Error(err))
		return err
	}

	if p.next != nil {
		p.log.Debug("Implementor Policy Check done, checking next policy")
		return p.next.Check(r)
	}
	return nil
}
