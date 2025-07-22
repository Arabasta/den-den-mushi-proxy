package types

type ConnectionStatus string

const (
	ConnectionStatusActive ConnectionStatus = "active"
	ConnectionStatusClosed ConnectionStatus = "closed"
)
