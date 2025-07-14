package pty_sessions

import "den-den-mushi-Go/pkg/types"

type Repository interface {
	FindById(id string) (*Entity, error)
	FindAllByState(st types.PtySessionState) ([]*Entity, error)
	FindAllByIp(ip string) ([]*Entity, error)
}
