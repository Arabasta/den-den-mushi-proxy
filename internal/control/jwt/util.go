package jwt

import (
	"den-den-mushi-Go/internal/control/change_request"
	"den-den-mushi-Go/internal/control/dto"
	"den-den-mushi-Go/internal/control/pty_sessions"
	dtopkg "den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
	"github.com/google/uuid"
)

func BuildConnForJoin(p *pty_sessions.Entity, r *dto.JoinRequest) *dtopkg.Connection {
	return &dtopkg.Connection{
		Server:  p.StartConnectionDetails.Server,
		Type:    p.StartConnectionDetails.Type,
		Purpose: p.StartConnectionDetails.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        r.UserId + "/" + uuid.NewString(),
			StartRole: r.StartRole,
		},
		PtySession: dtopkg.PtySession{
			Id:    p.Id,
			IsNew: false, // joining existing session
		},
		ChangeRequest: dtopkg.ChangeRequest{
			Id:                       p.StartConnectionDetails.ChangeRequest.Id,
			ImplementorGroup:         p.StartConnectionDetails.ChangeRequest.ImplementorGroup,
			EndTime:                  p.StartConnectionDetails.ChangeRequest.EndTime,
			ChangeGracePeriodMinutes: 30,
		},
		FilterType: p.StartConnectionDetails.FilterType,
	}
}

func BuildConnForStart(t types.ConnectionMethod, r *dto.StartRequest, cr *change_request.Entity, f types.Filter) *dtopkg.Connection {
	return &dtopkg.Connection{
		Server:  r.Server,
		Type:    t,
		Purpose: r.Purpose,
		UserSession: dtopkg.UserSession{
			Id:        r.UserId + "/" + uuid.NewString(),
			StartRole: types.Implementor, // always start as implementor
		},
		PtySession: dtopkg.PtySession{
			IsNew: true,
		},
		ChangeRequest: dtopkg.ChangeRequest{
			Id:                       r.ChangeID,
			ImplementorGroup:         "myimplgroup", // todo
			EndTime:                  cr.ChangeEndTime,
			ChangeGracePeriodMinutes: 30,
		},
		FilterType: f,
	}
}
