package dto

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

type MintRequest struct {
	Purpose      types.ConnectionPurpose `json:"purpose" binding:"required,oneof=change_request health_check"`
	ChangeID     string                  `json:"change_id,omitempty"`
	PtySessionId string                  `json:"pty_session_id,omitempty"`           // if not provided, a new session will be created
	StartRole    types.StartRole         `json:"start_role,omitempty"`               // required if PtySessionId is provided
	Type         types.ConnectionMethod  `json:"connection_type" binding:"required"` // temporarily from client, should be set based on server details
	Server       ServerInfo              `json:"server" binding:"required"`
	FilterType   types.Filter            `json:"filter_type"` // should be set based on server details
	UserId       string                  `json:"user_id"`     // temp, should be set with keycloak user id
}

type ServerInfo struct {
	IP     string `json:"ip" binding:"required"`
	Port   string `json:"port" binding:"required"`
	OSUser string `json:"os_user" binding:"required"`
}

type Connection struct {
	Server        ServerInfo              `json:"server" binding:"required"`
	Type          types.ConnectionMethod  `json:"type" binding:"required"` // todo: should be set based on server details
	Purpose       types.ConnectionPurpose `json:"purpose" binding:"required"`
	UserSession   UserSession             `json:"user_session"`
	PtySession    PtySession              `json:"pty_session"`
	ChangeRequest ChangeRequest           `json:"change_request,omitempty"`
	FilterType    types.Filter            `json:"filter_type" binding:"required"` // should be set based on server details
}

type UserSession struct {
	Id        string          `json:"id" binding:"required"`
	StartRole types.StartRole `json:"start_role,required"`
}

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
