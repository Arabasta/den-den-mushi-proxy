package request

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
)

type StartRequest struct {
	Purpose  types.ConnectionPurpose `json:"purpose" binding:"required,oneof=change_request health_check"`
	ChangeID string                  `json:"change_id,omitempty"`
	Server   dto.ServerInfo          `json:"server" binding:"required"`
}

func (s *StartRequest) GetPurpose() types.ConnectionPurpose { return s.Purpose }
func (s *StartRequest) GetChangeId() string                 { return s.ChangeID }
func (s *StartRequest) GetServerInfo() dto.ServerInfo       { return s.Server }

type StartAdapter struct {
	AdapterFields
	Req wrapper.WithAuth[StartRequest]
}

func (s *StartAdapter) GetUserId() string      { return s.Req.AuthCtx.UserID }
func (s *StartAdapter) GetUserOuGroup() string { return s.Req.AuthCtx.OuGroup }
