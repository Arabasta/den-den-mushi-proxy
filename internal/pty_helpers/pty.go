package pty_helpers

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
)

func Spawn(cmd *exec.Cmd) (*os.File, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return ptmx, nil
}

func Resize(ptmx *os.File, cols, rows uint16) error {
	return pty.Setsize(ptmx, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
}
