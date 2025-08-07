package request

import (
	"den-den-mushi-Go/pkg/dto"
	changerequestpkg "den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/iexpress"
	"den-den-mushi-Go/pkg/dto/implementor_groups"
	"den-den-mushi-Go/pkg/types"
)

// todo: refactor this garbage
type AdapterFields struct {
	Purpose                types.ConnectionPurpose
	ChangeID               string
	Server                 dto.ServerInfo
	CR                     *changerequestpkg.Record
	Iexpress               *iexpress.Record
	UserOuGroup            string
	UsersImplementorGroups []*implementor_groups.Record
}

func (a *AdapterFields) GetPurpose() types.ConnectionPurpose        { return a.Purpose }
func (a *AdapterFields) GetChangeId() string                        { return a.ChangeID }
func (a *AdapterFields) GetServerInfo() dto.ServerInfo              { return a.Server }
func (a *AdapterFields) GetChangeRequest() *changerequestpkg.Record { return a.CR }
func (a *AdapterFields) GetIexpress() *iexpress.Record {
	return a.Iexpress
}
func (a *AdapterFields) GetUsersImplementorGroups() []*implementor_groups.Record {
	return a.UsersImplementorGroups
}
