package pty_util

import (
	"github.com/creack/pty"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

func Spawn(cmd *exec.Cmd, log *zap.Logger) (*os.File, *exec.Cmd, error) {
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	Setpgid: true,
	//}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Warn("SSH exited with error", zap.Error(err))
		} else {
			log.Info("SSH exited cleanly")
		}
	}()

	return ptmx, cmd, nil
}

func Resize(ptmx *os.File, cols, rows uint16) error {
	return pty.Setsize(ptmx, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
}
