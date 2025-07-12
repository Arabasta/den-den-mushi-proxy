package types

type Filter string

const (
	Blacklist Filter = "blacklist"
	Whitelist Filter = "whitelist"
)
