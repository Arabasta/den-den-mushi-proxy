package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"io"
)

type OutputHandler struct{}

func (h *OutputHandler) Handle(pkt protocol.Packet, _ io.Writer, ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(pkt))
}
