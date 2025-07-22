package connections

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

type Model struct {
	ID           string                 `gorm:"primaryKey;column:id;size:191"`
	UserID       string                 `gorm:"user_id;size:191"`
	PtySessionID string                 `gorm:"column:pty_session_id;size:191"`
	StartRole    types.StartRole        `gorm:"start_role,not null;size:191"`
	Status       types.ConnectionStatus `gorm:"column:status,not null;size:191"`
	JoinTime     *time.Time             `gorm:"column:join_time"`
	LeaveTime    *time.Time             `gorm:"column:leave_time"`
}

func (Model) TableName() string {
	return "ddm_connections"
}
