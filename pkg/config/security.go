package config

type Security struct {
	IsEnabled           bool     `json:"isEnabled"`
	EnforceCsp          bool     `json:"enforceCsp"`
	AllowStyleAttribute bool     `json:"allowStyleAttribute"`
	ConnectSrc          []string `json:"connectSrc"`
}
