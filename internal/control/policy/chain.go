package policy

// todo:  CoR not good, refactor entire thing
type Policy[T any] interface {
	SetNext(n Policy[T])
	Check(r T) error
}
