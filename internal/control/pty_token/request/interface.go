package request

import (
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/implementor_groups"
	"den-den-mushi-Go/pkg/types"
)

type Ctx interface {
	GetPurpose() types.ConnectionPurpose
	GetChangeId() string
	GetServerInfo() dto.ServerInfo
	GetUserId() string
	GetChangeRequest() *changerequestpkg.Record
	GetUsersImplementorGroups() []*implementor_groups.Record
}

type HasJoinRequestFields interface {
	GetPtySessionId() string
	GetStartRole() types.StartRole
}
