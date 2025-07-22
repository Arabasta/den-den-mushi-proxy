package wrapper

import "den-den-mushi-Go/pkg/middleware"

type WithAuth[T any] struct {
	Body    T
	AuthCtx middleware.AuthContext
}
