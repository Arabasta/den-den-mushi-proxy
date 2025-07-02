package protocol

import "github.com/gorilla/websocket"

type Packet struct {
	Header Header
	Data   []byte
}

func Parse(msgType int, msg []byte) Packet {
	if msgType != websocket.BinaryMessage {
		return Packet{ParseError, nil}
	}
	if len(msg) < 1 {
		return Packet{ParseError, nil}
	}

	// todo: validate header

	return Packet{Header(msg[0]), msg[1:]}
}

func PacketToByte(pkt Packet) []byte {
	return append([]byte{byte(pkt.Header)}, pkt.Data...)
}
