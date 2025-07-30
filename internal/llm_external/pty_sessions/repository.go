package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
)

type Repository interface {
	FindAllByChangeRequestID(changeRequestID string) ([]*dto.Record, error)
}
