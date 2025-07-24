package handler

import "C"
import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"io"
	"time"
)

type SudoHandler struct{}

func (h *SudoHandler) Handle(pkt protocol.Packet, pty io.Writer) (string, error) {

	user := pkt.Data
	password := "opassword1234"
	// password := cyberark.DrawPassword(user)

	if _, err := pty.Write(constants.CtrlC); err != nil {
		return "", err
	}

	command := "su " + string(user) + "\n"
	if _, err := pty.Write([]byte(command)); err != nil {
		return "", err
	}

	// todo: make interactive and wait for the prompt
	// or maybe safer to just throw pty output away and
	// login then clear and resume
	time.Sleep(5 * time.Second)

	if _, err := pty.Write([]byte(password + "\n")); err != nil {
		return "", err
	}

	if _, err := pty.Write(constants.CtrlC); err != nil {
		return "", err
	}

	return "", nil
}
