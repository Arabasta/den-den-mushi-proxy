package connections

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

type Record struct {
	ID           string
	UserID       string
	PtySessionID string
	StartRole    types.StartRole
	Status       types.ConnectionStatus
	JoinTime     *time.Time
	LeaveTime    *time.Time
}
