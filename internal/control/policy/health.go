package policy

import (
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/implementor_groups"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type HealthcheckPolicy[T dto.RequestCtx] struct {
	next Policy[T]

	impGroupService *implementor_groups.Service
	hostService     *host.Service
	log             *zap.Logger
}

func NewHealthcheckPolicy[T dto.RequestCtx](hostService *host.Service, impGroupService *implementor_groups.Service,
	log *zap.Logger) *HealthcheckPolicy[T] {
	return &HealthcheckPolicy[T]{
		hostService:     hostService,
		impGroupService: impGroupService,
		log:             log,
	}
}

func (p *HealthcheckPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *HealthcheckPolicy[T]) Check(req T) error {
	// 1. skip healthcheck requests
	if req.GetPurpose() != types.Healthcheck {
		if p.next != nil {
			return p.next.Check(req)
		}
		return nil
	}

	// 2. get host type

	// 3. check if user is in host healthcheck group

	// 4. check is os user is a readonly user

	return nil
}
