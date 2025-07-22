package change_request

import "time"

type Record struct {
	ChangeRequestId   string
	Country           []string
	Lob               string
	Summary           string
	Description       string
	ChangeStartTime   *time.Time
	ChangeEndTime     *time.Time
	ImplementorGroups []string
	State             string
	CyberArkObjects   []string
}
