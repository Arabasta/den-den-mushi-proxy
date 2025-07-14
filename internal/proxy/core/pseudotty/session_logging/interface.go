package session_logging

type SessionLogger interface {
	WriteLine(line string) error
	Close() error
}
