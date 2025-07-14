package policy

import (
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/host"
	"go.uber.org/zap"
)

type OUPolicy[T dto.RequestCtx] struct {
	next Policy[T]

	hostService *host.Service
	log         *zap.Logger
}

func NewOUPolicy[T dto.RequestCtx](hostService *host.Service, log *zap.Logger) *OUPolicy[T] {
	return &OUPolicy[T]{
		hostService: hostService,
		log:         log,
	}
}

func (p *OUPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *OUPolicy[T]) Check(req T) error {
	// 1. check if server is OU? what the hell this?

	return nil
}
