package filters

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

type ListCR struct {
	TicketIDs         *[]string
	ImplementorGroups *[]string
	LOB               *string
	Country           *string
	StartTime         *time.Time
	EndTime           *time.Time
	PtySessionState   *types.PtySessionState
	Page              int
	PageSize          int
}

type CrHostPtySession struct {
	CRID  string
	IPs   []string
	State *types.PtySessionState // optional
}
