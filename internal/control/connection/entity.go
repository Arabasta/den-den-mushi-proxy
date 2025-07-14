package connection

import "den-den-mushi-Go/pkg/types"

type Entity struct {
	Id        string          `json:"id"`
	UserID    string          `json:"user_id"`
	StartRole types.StartRole `json:"start_role,required"`
	JoinTime  string          `json:"join_time,omitempty"`
	LeaveTime string          `json:"leave_time,omitempty"`
}
