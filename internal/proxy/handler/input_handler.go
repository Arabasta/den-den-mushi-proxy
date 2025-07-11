package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"io"
)

type InputHandler struct{}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {
	_, err := pty.Write(pkt.Data)
	if err != nil {
		return "", err
	}
	return "", err
}
