package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type OutputHandler struct{}

func (h *OutputHandler) Handle(pkt protocol.Packet, _ io.Writer, ws *websocket.Conn, _ *token.Claims) (string, error) {
	return string(pkt.Data), ws.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(pkt))
}
