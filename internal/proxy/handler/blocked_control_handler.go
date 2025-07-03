package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type BlockedControlHandler struct{}

func (h *BlockedControlHandler) Handle(pkt protocol.Packet, _ io.Writer, ws *websocket.Conn, _ *token.Claims) (string, error) {
	return "", ws.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(pkt))
}
