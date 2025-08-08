package request

import (
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/iexpress"
	"den-den-mushi-Go/pkg/dto/implementor_groups"
	"den-den-mushi-Go/pkg/types"
)

//  todo refactor this garbage

type Ctx interface {
	GetPurpose() types.ConnectionPurpose
	GetChangeId() string
	GetServerInfo() dto.ServerInfo
	GetUserId() string
	GetUserOuGroup() string
	GetChangeRequest() *changerequestpkg.Record
	GetIexpress() *iexpress.Record
	GetUsersImplementorGroups() []*implementor_groups.Record
	GetStartRole() types.StartRole
}

type HasJoinRequestFields interface {
	GetPtySessionId() string
	GetStartRole() types.StartRole
}
