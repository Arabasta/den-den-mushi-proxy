package dto

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/types"
)

type RequestCtx interface {
	GetPurpose() types.ConnectionPurpose
	GetChangeId() string
	GetServerInfo() dto.ServerInfo
	GetUserId() string
	GetUserGroups() []string
}
