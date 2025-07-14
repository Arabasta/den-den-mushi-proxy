package policy

type ChainBuilder[T any] struct {
	policies []Policy[T]
}

func NewPolicyBuilder[T any]() *ChainBuilder[T] {
	return &ChainBuilder[T]{
		policies: make([]Policy[T], 0),
	}
}

func (b *ChainBuilder[T]) Add(p Policy[T]) *ChainBuilder[T] {
	if p == nil {
		return b
	}
	if len(b.policies) > 0 {
		last := b.policies[len(b.policies)-1]
		last.SetNext(p)
	}
	b.policies = append(b.policies, p)
	return b
}

func (b *ChainBuilder[T]) Build() Policy[T] {
	if len(b.policies) == 0 {
		return nil
	}
	return b.policies[0]
}
