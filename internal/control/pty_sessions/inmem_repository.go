package pty_sessions

import "den-den-mushi-Go/pkg/types"

type InMemRepository struct {
	sessions map[string]*Entity
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		sessions: make(map[string]*Entity),
	}
}

func (r *InMemRepository) FindById(id string) (*Entity, error) {
	session, exists := r.sessions[id]
	if !exists {
		return nil, nil // or return an error if preferred
	}
	return session, nil
}

func (r *InMemRepository) FindAllByState(st types.PtySessionState) ([]*Entity, error) {
	var result []*Entity
	for _, session := range r.sessions {
		if session.State == st {
			result = append(result, session)
		}
	}
	return result, nil
}

func (r *InMemRepository) FindAllByIp(ip string) ([]*Entity, error) {
	var result []*Entity
	for _, session := range r.sessions {
		if session.StartConnectionDetails.Server.IP == ip {
			result = append(result, session)
		}
	}
	return result, nil
}
