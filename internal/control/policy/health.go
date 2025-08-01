package policy

import (
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/internal/control/policy/validators"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type HealthcheckPolicy[T request.Ctx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	hostService     *host.Service
	log             *zap.Logger
	v               *validators.Validator
}

func NewHealthcheckPolicy[T request.Ctx](hostService *host.Service, impGroupService *implementor_groups.Service,
	log *zap.Logger, v *validators.Validator) *HealthcheckPolicy[T] {
	return &HealthcheckPolicy[T]{
		hostService:     hostService,
		impGroupService: impGroupService,
		log:             log,
		v:               v,
	}
}

func (p *HealthcheckPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *HealthcheckPolicy[T]) Check(r T) error {
	p.log.Debug("Checking healthcheck policy...")

	// skip healthcheck requests
	if r.GetPurpose() != types.Healthcheck {
		if p.next != nil {
			return p.next.Check(r)
		}
		return nil
	}

	// 2. get host type todo how to map?
	// hostType, err := p.hostService.GetHostType(r.GetServerInfo().IP)

	// 3. check if user is in host healthcheck group todo uncomment this when host type mapping is up
	//if err := p.v.IsOuGroupTypeMatch(r.GetUserOuGroup(), hostType); err != nil {
	//	p.log.Warn("User OU group type validation failed", zap.String("userId", r.GetUserId()), zap.Error(err))
	//	return err
	//}

	// 4. check is os user is a valid readonly user for that user?, todo implement this once we have the mapping in DB

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
