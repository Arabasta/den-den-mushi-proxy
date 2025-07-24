package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/types"
)

type Repository interface {
	FindById(id string) (*dto.Record, error)
	FindByStartConnChangeRequestIdsAndState(changeIDs []string, state types.PtySessionState) ([]*dto.Record, error)
	FindAllByChangeRequestIDAndServerIPs(changeRequestID string, ips []string) ([]*dto.Record, error)
	FindAllByStartConnServerIpsAndState(hostips []string, state *types.PtySessionState) ([]*dto.Record, error)
}
