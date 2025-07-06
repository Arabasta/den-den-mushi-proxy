package pty_sessions

import (
	"errors"
	"sync"
)

type InMemRepository struct {
	mu       sync.RWMutex
	sessions map[string]*Entity
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		sessions: make(map[string]*Entity),
	}
}

func (r *InMemRepository) Save(session *Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID] = session
	return nil
}

func (r *InMemRepository) Get(id string) (*Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

func (r *InMemRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	return nil
}

func (r *InMemRepository) List() ([]*Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Entity, 0, len(r.sessions))
	for _, s := range r.sessions {
		result = append(result, s)
	}
	return result, nil
}

// err constant
var ErrNotFound = errors.New("session not found")
