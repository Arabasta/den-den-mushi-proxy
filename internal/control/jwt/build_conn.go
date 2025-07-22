package jwt

import (
	"den-den-mushi-Go/internal/control/pty_token/request"
	dtopkg "den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"github.com/google/uuid"
)

func BuildConnForStart(t types.ConnectionMethod, r wrapper.WithAuth[request.StartRequest], cr *change_request.Record,
	f types.Filter) *dtopkg.Connection {
	return &dtopkg.Connection{
		Server: dtopkg.ServerInfo{
			OSUser: r.Body.Server.OSUser,
			IP:     r.Body.Server.IP,
			Port:   "22", // todo: demo
		},
		Type:    t,
		Purpose: r.Body.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        r.AuthCtx.UserID + "/" + uuid.NewString(),
			StartRole: types.Implementor, // always start as implementor
		},
		PtySession: dtopkg.PtySession{
			IsNew: true,
		},
		ChangeRequest: func() dtopkg.ChangeRequest {
			if r.Body.Purpose == types.Change {
				return dtopkg.ChangeRequest{
					Id:                r.Body.ChangeID,
					ImplementorGroups: cr.ImplementorGroups,
					EndTime:           *cr.ChangeEndTime,
				}
			}
			return dtopkg.ChangeRequest{}
		}(),
		FilterType: f,
	}
}

func BuildConnForJoin(p *pty_sessions.Record, r wrapper.WithAuth[request.JoinRequest]) *dtopkg.Connection {
	return &dtopkg.Connection{
		Server:  p.StartConnectionDetails.Server,
		Type:    p.StartConnectionDetails.Type,
		Purpose: p.StartConnectionDetails.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        r.AuthCtx.UserID + "/" + uuid.NewString(),
			StartRole: r.Body.StartRole,
		},
		PtySession: dtopkg.PtySession{
			Id:    p.ID,
			IsNew: false, // joining existing session
		},
		ChangeRequest: func() dtopkg.ChangeRequest {
			if p.StartConnectionDetails.Purpose == types.Change {
				return dtopkg.ChangeRequest{
					Id:                p.StartConnectionDetails.ChangeRequest.Id,
					ImplementorGroups: p.StartConnectionDetails.ChangeRequest.ImplementorGroups,
					EndTime:           p.StartConnectionDetails.ChangeRequest.EndTime,
				}
			}
			return dtopkg.ChangeRequest{}
		}(),
		FilterType: p.StartConnectionDetails.FilterType,
	}
}
