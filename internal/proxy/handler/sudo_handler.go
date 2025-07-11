package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"io"
)

type SudoHandler struct{}

func (h *SudoHandler) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {

	if true { // todo: this should check user's claims for sudo
		user := pkt.Data
		command := "su " + string(user) + " -i\n"

		_, err := pty.Write([]byte(command))
		return "", err
	}
	return "", nil
}
