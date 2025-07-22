package types

type PtySessionState string

const (
	Created PtySessionState = "created"
	Active  PtySessionState = "active"
	Closed  PtySessionState = "closed"
)
