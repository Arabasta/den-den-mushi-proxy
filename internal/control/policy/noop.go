package policy

import (
	"den-den-mushi-Go/internal/control/ep/pty_token/request"
	"go.uber.org/zap"
)

type NoopPolicy[T request.Ctx] struct {
	next Policy[T]
	log  *zap.Logger
}

func NewNoopPolicy[T request.Ctx](log *zap.Logger) *NoopPolicy[T] {
	return &NoopPolicy[T]{
		log: log,
	}
}

func (p *NoopPolicy[T]) SetNext(n Policy[T]) {
	p.next = n
}

func (p *NoopPolicy[T]) Check(_ T) error {
	p.log.Debug("Noop Policy enabled. Skipping checks...")
	return nil
}
