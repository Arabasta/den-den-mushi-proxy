package pty_sessions

import "den-den-mushi-Go/pkg/dto"

type Entity struct {
	ID                       string         `json:"id"`
	ProxyServer              string         `json:"proxy_server"`
	InitialConnectionDetails dto.Connection `json:"initial_connection_details"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`

	Status         string `json:"status"`
	LastActivityAt int64  `json:"last_activity_at"`

	CreatedBy          string            `json:"created_by"`
	CurrentImplementor string            `json:"current_implementor"`
	Observers          map[string]string `json:"observers"`
}
