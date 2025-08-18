package policy

import (
	"den-den-mushi-Go/internal/control/app/pty_token/request"
	"den-den-mushi-Go/internal/control/core/host"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type ObserverPolicy[T request.Ctx] struct {
	next Policy[T]

	hostService *host.Service
	log         *zap.Logger
	v           *validators.Validator
}

func NewObserverPolicy[T request.Ctx](log *zap.Logger, v *validators.Validator) *ObserverPolicy[T] {
	return &ObserverPolicy[T]{
		log: log,
		v:   v,
	}
}

func (p *ObserverPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *ObserverPolicy[T]) Check(r T) error {
	p.log.Debug("Checking observer policy...")

	if r.GetStartRole() != types.Observer {
		p.log.Debug("Request is not for Observer role, skipping Observer policy check")
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	p.log.Debug("Request is for Observer role, checking for OU presence")
	// user must have OU group
	if err := p.v.HasOuGroup(r.GetUserOuGroup()); err != nil {
		p.log.Warn("User does not have OU group", zap.String("userId", r.GetUserId()), zap.Error(err))
		return err
	}

	if p.next != nil {
		p.log.Debug("Observer Policy Check done, checking next policy")
		return p.next.Check(r)
	}
	return nil
}
