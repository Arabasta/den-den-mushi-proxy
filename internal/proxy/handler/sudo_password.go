package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"fmt"
	"go.uber.org/zap"
	"io"
	"time"
)

type SudoPassword struct{}

func (h *SudoPassword) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {
	password := append(pkt.Data, constants.Enter...)
	fmt.Println("sudo password received", zap.String("password", string(password)))
	if _, err := pty.Write(password); err != nil {
		return "", err
	}
	time.Sleep(1500 * time.Millisecond)

	if _, err := pty.Write(append(constants.CtrlC, constants.Enter...)); err != nil {
		return "", err
	}
	time.Sleep(1000 * time.Millisecond)

	// clear the screen just in case password leaks
	if _, err := pty.Write(append([]byte("clear"), constants.Enter...)); err != nil {
		return "", err
	}

	return "", nil
}
