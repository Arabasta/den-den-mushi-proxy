package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"io"
	"time"
)

type SudoUsername struct{}

func (h *SudoUsername) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {
	// clear the current
	if _, err := pty.Write(append(constants.CtrlC, constants.Enter...)); err != nil {
		return "", err
	}
	time.Sleep(250 * time.Millisecond)

	suCommand := "su - " + string(pkt.Data)
	if _, err := pty.Write(append([]byte(suCommand), constants.Enter...)); err != nil {
		return "", err
	}
	time.Sleep(500 * time.Millisecond)

	return "", nil
}
