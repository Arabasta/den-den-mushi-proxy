package types

type PtySessionState string

const (
	Created  PtySessionState = "created"
	Active   PtySessionState = "active"
	Inactive PtySessionState = "inactive"
	Closed   PtySessionState = "closed"
)
