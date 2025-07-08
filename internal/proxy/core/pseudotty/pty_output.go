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
				s.Log.Info("PTY session ended normally")
				s.closePty()
			} else {
				s.Log.Error("Error reading from pty", zap.Error(err))
				s.closePty()
			}
			close(s.outboundCh)
			return
		}

		data := append([]byte{}, buf[:n]...)
		pkt := protocol.Packet{
			Header: protocol.Output,
			Data:   data,
		}

		s.outboundCh <- pkt

		s.mu.Lock()
		s.ptyLastPackets = append(s.ptyLastPackets, pkt)
		maxLastPackets := 20 // todo: make configurable and maybe track line or something
		if len(s.ptyLastPackets) >= maxLastPackets {
			s.ptyLastPackets = s.ptyLastPackets[1:]
		}
		s.mu.Unlock()

		s.Log.Info("Pty Output", zap.ByteString("data", data))
	}
}
