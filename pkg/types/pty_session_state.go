package types

type PtySessionState string

const (
	Created PtySessionState = "created"
	Running PtySessionState = "active"
	Closed  PtySessionState = "closed"
)
