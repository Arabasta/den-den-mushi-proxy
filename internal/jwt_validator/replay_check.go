package jwt_validator

import (
	"sync"
	"time"
)

type jtiStore struct {
	seen sync.Map
	ttl  time.Duration
}

func (s *jtiStore) Check(id string) bool {
	if _, found := s.seen.Load(id); found {
		return false
	}
	s.seen.Store(id, struct{}{})
	go func() { <-time.After(s.ttl); s.seen.Delete(id) }()
	return true
}
