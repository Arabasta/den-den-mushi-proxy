package policy

import (
	"den-den-mushi-Go/internal/control/pty_token/request"
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

func (p *NoopPolicy[T]) Check(_ T) error { return nil }
