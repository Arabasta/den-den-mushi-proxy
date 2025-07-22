package policy

type Policy[T any] interface {
	SetNext(n Policy[T])
	Check(r T) error
}
