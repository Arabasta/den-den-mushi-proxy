package jti

type Repository interface {
	// consumeIfNotExists returns true if the JTI was consumed, false if it already existed
	consumeIfNotExists(jti *Record) (bool, error)
}
