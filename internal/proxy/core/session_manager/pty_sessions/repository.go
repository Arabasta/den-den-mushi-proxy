package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
)

type Repository interface {
	FindById(id string) (*dto.Record, error)
	Save(*dto.Record) error
	FindActiveByProxyHostWithConnections(proxyHost string) ([]*dto.Record, error)
	CloseSessionsAndConnections(sessionIDs []string, connectionIDs []string) error
}
