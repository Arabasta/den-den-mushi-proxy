package jwt_service

import (
	"sync"
	"time"
)

type jtiStore struct {
	seen sync.Map
	ttl  time.Duration
}

func (s *jtiStore) isConsumed(id string) bool {
	_, found := s.seen.Load(id)
	if found {
		return true
	}

	return false
}

func (s *jtiStore) consume(id string) bool {
	s.seen.Store(id, struct{}{})

	// for now just delete after ttl, in prod maybe don't cleanup?
	// todo: discuss this
	go func() { <-time.After(s.ttl); s.seen.Delete(id) }()
	return true
}
