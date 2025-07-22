package policy

type Builder[T any] struct {
	policies []Policy[T]
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		policies: make([]Policy[T], 0),
	}
}

func (b *Builder[T]) Add(p Policy[T]) *Builder[T] {
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

func (b *Builder[T]) Build() Policy[T] {
	if len(b.policies) == 0 {
		return nil
	}
	return b.policies[0]
}
