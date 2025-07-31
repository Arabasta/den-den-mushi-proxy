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
	p.log.Debug("Checking OU policy.... This doesnt do anythign for now")

	if p.next != nil {
		p.log.Debug("OuPolicy Check called, checking next policy")
		return p.next.Check(r)
	}
	p.log.Debug("OUPolicy Check called, but no next policy set")
	return nil
}
