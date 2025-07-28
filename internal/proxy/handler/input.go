package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"io"
)

type Input struct{}

func (h *Input) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {
	_, err := pty.Write(pkt.Data)
	if err != nil {
		return "", err
	}
	return "", err
}
