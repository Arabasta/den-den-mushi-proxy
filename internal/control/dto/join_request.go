package dto

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
)

type JoinRequest struct {
	PtySessionId string          `json:"pty_session_id,required"`
	StartRole    types.StartRole `json:"start_role,required,oneof=implementor observer"`
	// populated by keycloak
	UserId   string   `json:"-"`
	OuGroups []string `json:"-"`
}

type JoinAdapter struct {
	Req      *JoinRequest
	Purpose  types.ConnectionPurpose
	ChangeID string
	Server   dto.ServerInfo
}

func (j *JoinAdapter) GetPurpose() types.ConnectionPurpose { return j.Purpose }
func (j *JoinAdapter) GetChangeId() string                 { return j.ChangeID }
func (j *JoinAdapter) GetServerInfo() dto.ServerInfo       { return j.Server }
func (j *JoinAdapter) GetUserId() string                   { return j.Req.UserId }
func (j *JoinAdapter) GetUserGroups() []string             { return j.Req.OuGroups }
