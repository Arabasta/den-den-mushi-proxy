package connect

import (
	"context"
	"den-den-mushi-Go/pkg/token"
	"os"
	"os/exec"
)

type ConnectionMethod interface {
	Connect(ctx context.Context, claims *token.Claims) (*os.File, *exec.Cmd, error)
}
