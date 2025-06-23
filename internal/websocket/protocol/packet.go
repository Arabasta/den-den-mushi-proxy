package protocol

type Packet struct {
	Header Header
	Data   []byte
}

func Parse(msg []byte) Packet {
	if len(msg) < 1 {
		return Packet{ParseError, nil}
	}

	return Packet{Header(msg[0]), msg[1:]}
}

func PacketToByte(pkt Packet) []byte {
	return append([]byte{byte(pkt.Header)}, pkt.Data...)
}
