package jwt

import (
	request2 "den-den-mushi-Go/internal/control/ep/pty_token/request"
	dtopkg "den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/dto/change_request"
	iexpress2 "den-den-mushi-Go/pkg/dto/iexpress"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"github.com/google/uuid"
)

func BuildConnForStart(t types.ConnectionMethod, r wrapper.WithAuth[request2.StartRequest], cr *change_request.Record,
	exp *iexpress2.Record, f types.Filter, port string, allowedSuOsUsers []string, serverFQDNTmpTillRefactor string) *dtopkg.Connection {
	userSessionId := r.AuthCtx.UserID + "/" + uuid.NewString()

	return &dtopkg.Connection{
		Server: dtopkg.ServerInfo{
			OSUser: r.Body.Server.OSUser,
			IP:     r.Body.Server.IP,
			Port:   port,
		},
		Type:    t,
		Purpose: r.Body.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        userSessionId,
			StartRole: types.Implementor, // always start as implementor
		},
		PtySession: dtopkg.PtySession{
			IsNew:                true,
			InitialUserSessionId: userSessionId,
		},
		ChangeRequest: func() dtopkg.ChangeRequest {
			if r.Body.Purpose == types.Change {
				return dtopkg.ChangeRequest{
					Id:                r.Body.ChangeID,
					ImplementorGroups: cr.ImplementorGroups,
					EndTime:           *cr.ChangeEndTime,
				}
			} else if r.Body.Purpose == types.IExpress {
				return dtopkg.ChangeRequest{
					Id: r.Body.ChangeID,
					// todo will be empty string if not set
					ImplementorGroups: []string{exp.ApproverGroup1, exp.ApproverGroup2, exp.MDApproverGroup},
					EndTime:           *exp.EndTime,
				}
			}
			return dtopkg.ChangeRequest{}
		}(),
		FilterType:                f,
		AllowedSuOsUsers:          allowedSuOsUsers,
		ServerFQDNTmpTillRefactor: serverFQDNTmpTillRefactor,
	}
}

func BuildConnForJoin(p *pty_sessions.Record, r wrapper.WithAuth[request2.JoinRequest]) *dtopkg.Connection {
	return &dtopkg.Connection{
		Server:  p.StartConnectionDetails.Server,
		Type:    p.StartConnectionDetails.Type,
		Purpose: p.StartConnectionDetails.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        r.AuthCtx.UserID + "/" + uuid.NewString(),
			StartRole: r.Body.StartRole,
		},
		PtySession: dtopkg.PtySession{
			Id:                   p.ID,
			IsNew:                false, // joining existing session
			InitialUserSessionId: p.StartConnUserSessionID,
		},
		ChangeRequest: func() dtopkg.ChangeRequest {
			if p.StartConnectionDetails.Purpose == types.Change {
				return dtopkg.ChangeRequest{
					Id:                p.StartConnectionDetails.ChangeRequest.Id,
					ImplementorGroups: p.StartConnectionDetails.ChangeRequest.ImplementorGroups,
					EndTime:           p.StartConnectionDetails.ChangeRequest.EndTime,
				}
			} else if p.StartConnectionDetails.Purpose == types.IExpress {
				return dtopkg.ChangeRequest{
					Id: p.StartConnectionDetails.ChangeRequest.Id,
					// todo will be empty string if not set
					ImplementorGroups: p.StartConnectionDetails.ChangeRequest.ImplementorGroups,
					EndTime:           p.StartConnectionDetails.ChangeRequest.EndTime,
				}
			}
			return dtopkg.ChangeRequest{}
		}(),
		FilterType: p.StartConnectionDetails.FilterType}
}
