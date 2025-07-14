package dto

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
)

type StartRequest struct {
	Purpose  types.ConnectionPurpose `json:"purpose" binding:"required,oneof=change_request health_check"`
	ChangeID string                  `json:"change_id,omitempty"`
	Server   dto.ServerInfo          `json:"server" binding:"required"`
	// populated by keycloak
	UserId   string   `json:"-"`
	OuGroups []string `json:"-"`
}

func (s *StartRequest) GetPurpose() types.ConnectionPurpose { return s.Purpose }
func (s *StartRequest) GetChangeId() string                 { return s.ChangeID }
func (s *StartRequest) GetServerInfo() dto.ServerInfo       { return s.Server }
func (s *StartRequest) GetUserId() string                   { return s.UserId }
func (s *StartRequest) GetUserGroups() []string             { return s.OuGroups }
