package handler

import (
	"den-den-mushi-Go/internal/protocol"
	"den-den-mushi-Go/internal/pty_util"
	"den-den-mushi-Go/pkg/token"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"io"
	"os"
)

type ResizeHandler struct{}

func (h *ResizeHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn, _ *token.Claims) (string, error) {
	if len(pkt.Data) != 4 {
		return "", nil
	}

	cols := binary.BigEndian.Uint16(pkt.Data[0:2])
	rows := binary.BigEndian.Uint16(pkt.Data[2:4])

	if f, ok := pty.(*os.File); ok {
		return "", pty_util.Resize(f, cols, rows)
	}
	return "", nil
}
