package jti

type Repository interface {
	tryGet(jti string) (bool, error)
	save(jti string) error
}
