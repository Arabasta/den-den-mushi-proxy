package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"io"
	"time"
)

type SudoPassword struct{}

func (h *SudoPassword) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {
	password := append(pkt.Data, constants.Enter...)
	if _, err := pty.Write(password); err != nil {
		return "", err
	}
	time.Sleep(1500 * time.Millisecond)

	return "", nil
}
