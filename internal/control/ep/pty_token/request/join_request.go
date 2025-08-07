package request

import (
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
)

type JoinRequest struct {
	PtySessionId string          `json:"pty_session_id,required"`
	StartRole    types.StartRole `json:"start_role,required,oneof=implementor observer"`
}

type JoinAdapter struct {
	AdapterFields
	Req wrapper.WithAuth[JoinRequest]
}

func (j *JoinAdapter) GetPtySessionId() string {
	return j.Req.Body.PtySessionId
}

func (j *JoinAdapter) GetStartRole() types.StartRole {
	return j.Req.Body.StartRole
}

func (j *JoinAdapter) GetUserId() string      { return j.Req.AuthCtx.UserID }
func (j *JoinAdapter) GetUserOuGroup() string { return j.Req.AuthCtx.OuGroup }
