package dto

import "time"

type MintRequest struct {
	Purpose      ConnectionPurpose `json:"purpose" binding:"required,oneof=change_request health_check"`
	ChangeID     string            `json:"change_id,omitempty"`
	PtySessionId string            `json:"pty_session_id,omitempty"`           // if not provided, a new session will be created
	StartRole    StartRole         `json:"start_role,omitempty"`               // required if PtySessionId is provided
	Type         ConnectionType    `json:"connection_type" binding:"required"` // temporarily from client, should be set based on server details
	Server       ServerInfo        `json:"server" binding:"required"`
	FilterType   FilterType        `json:"filter_type"` // should be set based on server details
}

type ServerInfo struct {
	IP     string `json:"ip" binding:"required"`
	Port   string `json:"port" binding:"required"`
	OSUser string `json:"os_user" binding:"required"`
}

type Connection struct {
	Server        ServerInfo        `json:"server" binding:"required"`
	Type          ConnectionType    `json:"type" binding:"required"` // todo: should be set based on server details
	Purpose       ConnectionPurpose `json:"purpose" binding:"required"`
	UserSession   UserSession       `json:"user_session"`
	PtySession    PtySession        `json:"pty_session"`
	ChangeRequest ChangeRequest     `json:"change_request,omitempty"`
	FilterType    FilterType        `json:"filter_type" binding:"required"` // should be set based on server details
}

type UserSession struct {
	Id        string    `json:"id" binding:"required"`
	StartRole StartRole `json:"start_role,required"`
}

type StartRole string

const (
	Implementor StartRole = "implementor"
	Observer    StartRole = "observer"
)

type PtySession struct {
	Id                string `json:"id,omitempty"`
	IsNew             bool   `json:"is_new,omitempty"` // true if creating new session, false if joining existing
	IsObserverEnabled bool   `json:"is_observer_enabled,omitempty"`
	MaxObservers      int    `json:"max_observers,omitempty"`
	//MaxHeadlessMinutes        time.Duration `json:"max_headless_minutes,omitempty"`         // in minutes
	MaxSessionDurationMinutes time.Duration `json:"max_session_duration_minutes,omitempty"` // in minutes
}

type ChangeRequest struct {
	Id                       string        `json:"id" binding:"required"`
	ImplementorGroup         string        `json:"implementor_group" binding:"required"`
	EndTime                  string        `json:"end_time" binding:"required"` // ISO 8601 format
	ChangeGracePeriodMinutes time.Duration `json:"change_grace_period_minutes" binding:"required"`
}

type ConnectionPurpose string

const (
	Change      ConnectionPurpose = "change_request"
	Healthcheck ConnectionPurpose = "health_check"
)

type ConnectionType string

const (
	/* For development purposes only */
	LocalShell ConnectionType = "local_shell"
	SshTestKey ConnectionType = "ssh_test_key"
	/* For development purposes only End*/

	/* New Connection */
	SshOrchestratorKey ConnectionType = "ssh_orchestrator_key"
	SshPassword        ConnectionType = "ssh_password"
)

type FilterType string

const (
	Blacklist FilterType = "blacklist"
	Whitelist FilterType = "whitelist"
)
