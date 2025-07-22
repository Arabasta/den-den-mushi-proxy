package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/types"
	"time"
)

type Record struct {
	ID                     string
	ProxyHostName          string
	ProxyDetails           *proxy_host.Model
	StartConnectionDetails dto.Connection
	CreatedBy              string
	StartTime              *time.Time
	EndTime                *time.Time
	State                  types.PtySessionState
	LastActivity           *time.Time

	// start connection details
	StartConnServerIP        string
	StartConnServerPort      string
	StartConnServerOSUser    string
	StartConnType            types.ConnectionMethod
	StartConnPurpose         types.ConnectionPurpose
	StartConnUserSessionID   string
	StartConnChangeRequestID string
	StartConnFilterType      types.Filter

	Connections []*connections.Record
}
