package jti

type repository interface {
	tryGetJti(jti string) (bool, error)
	addJti(jti string) error
}
