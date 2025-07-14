package host

type Entity struct {
	IpAddress   string `json:"ip_address"`
	HostName    string `json:"host_name"`
	OsType      string `json:"os_type"`
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Country     string `json:"country"`
}
