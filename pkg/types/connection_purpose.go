package types

type ConnectionPurpose string

const (
	Change      ConnectionPurpose = "change_request"
	Healthcheck ConnectionPurpose = "health_check"
)
