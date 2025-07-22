package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/types"
	"time"
)

type Model struct {
	ID           string                `gorm:"primaryKey;column:id;size:191"`
	CreatedBy    string                `gorm:"column:created_by;size:191"`
	StartTime    *time.Time            `gorm:"column:start_time"`
	EndTime      *time.Time            `gorm:"column:end_time"`
	State        types.PtySessionState `gorm:"column:state;size:191"`
	LastActivity *time.Time            `gorm:"column:last_activity"`

	// One-to-one: session uses 1 proxy
	ProxyHostName string            `gorm:"column:proxy_host_name;size:191"`
	ProxyDetails  *proxy_host.Model `gorm:"foreignKey:ProxyHostName;references:HostName"`

	// start connection details
	StartConnServerIP        string                  `gorm:"column:start_conn_server_ip;size:191"`
	StartConnServerPort      string                  `gorm:"column:start_conn_server_port;size:191"`
	StartConnServerOSUser    string                  `gorm:"column:start_conn_server_os_user;size:191"`
	StartConnType            types.ConnectionMethod  `gorm:"column:start_conn_type;size:191"`
	StartConnPurpose         types.ConnectionPurpose `gorm:"column:start_conn_purpose;size:191"`
	StartConnUserSessionID   string                  `gorm:"column:start_conn_user_session_id;size:191"`
	StartConnChangeRequestID string                  `gorm:"column:start_conn_cr_id;size:191"`
	StartConnFilterType      types.Filter            `gorm:"column:start_conn_filter_type;size:191"`

	// One-to-many: session has many connections
	Connections []*connections.Model `gorm:"foreignKey:PtySessionID;references:ID"`
}

func (Model) TableName() string {
	return "ddm_pty_sessions"
}
