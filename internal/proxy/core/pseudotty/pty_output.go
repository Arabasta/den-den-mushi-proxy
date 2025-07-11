package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"go.uber.org/zap"
	"io"
)

// readPty and add data to outbound channel
func (s *Session) readPty() {
	buf := make([]byte, 4096)
	for {
		n, err := s.Pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.log.Info("PTY session ended normally")
			} else {
				s.log.Error("Error reading from pty", zap.Error(err))
				s.logL("Error reading from pty, shutting down session")
				s.outboundCh <- protocol.Packet{
					Header: protocol.Error,
					Data:   []byte("Error reading from pty, shutting down session: " + err.Error()),
				}
			}
			close(s.outboundCh)
			s.EndSession()
			return
		}

		data := append([]byte{}, buf[:n]...)
		pkt := protocol.Packet{
			Header: protocol.Output,
			Data:   data,
		}

		s.outboundCh <- pkt
		s.addPtyLastPacket(pkt)

		s.log.Info("Pty Output", zap.ByteString("data", data))
	}
}

func (s *Session) addPtyLastPacket(pkt protocol.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ptyLastPackets = append(s.ptyLastPackets, pkt)
	maxLastPackets := 100 // todo: make configurable and maybe track line or something
	if len(s.ptyLastPackets) >= maxLastPackets {
		s.ptyLastPackets = s.ptyLastPackets[1:]
	}
	s.log.Info("Added pty last packet", zap.Any("packet", pkt))
}
