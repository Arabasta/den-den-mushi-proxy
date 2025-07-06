package pty_sessions

type Repository interface {
	Save(session *Entity) error
	Get(id string) (*Entity, error)
	Delete(id string) error
	List() ([]*Entity, error)
}
