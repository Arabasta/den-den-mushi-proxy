package policy

type Policy[T any] interface {
	SetNext(n Policy[T])
	Check(req T) error
}
