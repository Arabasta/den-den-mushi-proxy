package policy

import (
	"den-den-mushi-Go/internal/control/host"
	"den-den-mushi-Go/internal/control/pty_token/request"
	"go.uber.org/zap"
)

type OUPolicy[T request.Ctx] struct {
	next Policy[T]

	hostService *host.Service
	log         *zap.Logger
}

func NewOUPolicy[T request.Ctx](hostService *host.Service, log *zap.Logger) *OUPolicy[T] {
	return &OUPolicy[T]{
		hostService: hostService,
		log:         log,
	}
}

func (p *OUPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *OUPolicy[T]) Check(r T) error {
	// 1. check if server is OU? what the hell this?

	if p.next != nil {
		return p.next.Check(r)
	}
	return nil
}
