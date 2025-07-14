package change_request

type Entity struct {
	ChangeRequestId   string   `json:"change_request_id"`
	ChangeStartTime   string   `json:"change_start_time"`
	ChangeEndTime     string   `json:"change_end_time"`
	ImplementorGroups []string `json:"implementor_groups"`
	CyberArkObjects   []string `json:"cyber_ark_objects"`
	State             string   `json:"state"`
}
