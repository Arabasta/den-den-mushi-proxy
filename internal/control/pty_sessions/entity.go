package pty_sessions

import (
	"den-den-mushi-Go/internal/control/connection"
	"den-den-mushi-Go/internal/control/proxy_lb"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
)

type Entity struct {
	Id                     string                `json:"id"`
	ProxyDetails           proxy_lb.Entity       `json:"proxy_details"`
	StartConnectionDetails dto.Connection        `json:"start_connection_details"`
	CreatedBy              string                `json:"created_by"`
	StartTime              string                `json:"start_time"`
	EndTime                string                `json:"end_time,omitempty"`
	State                  types.PtySessionState `json:"state,omitempty"`
	LastActivity           string                `json:"last_activity,omitempty"`
	ActiveConnections      []connection.Entity   `json:"active_connections"`
	LivetimeConnections    []connection.Entity   `json:"livetime_connections"`
}
