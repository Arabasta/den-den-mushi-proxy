package pty_util

import (
	"github.com/creack/pty"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

func Spawn(cmd *exec.Cmd) (*os.File, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Warn("SSH exited with error", zap.Error(err))
		} else {
			log.Info("SSH exited cleanly")
		}
	}()

	return ptmx, nil
}

func Resize(ptmx *os.File, cols, rows uint16) error {
	return pty.Setsize(ptmx, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
}
