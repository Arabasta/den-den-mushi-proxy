package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"go.uber.org/zap"
	"io"
)

// readPtyLoop and add data to outbound channel
func (s *Session) readPtyLoop() {
	//maxBufSize := s.cfg.Proxy.Pty.MaxBufferSize
	buf := make([]byte, 512)
	for {
		n, err := s.pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.log.Info("PTY session ended normally")
			} else {
				s.log.Error("Error reading from pty", zap.Error(err))
				s.logL("Error reading from pty, shutting down session")
				pkt := protocol.Packet{
					Header: protocol.Error,
					Data:   []byte("Error reading from pty, shutting down session: " + err.Error()),
				}
				s.ptyOutput.Add(pkt)
			}
			s.EndSession()
			return
		}

		data := append([]byte{}, buf[:n]...)
		pkt := protocol.Packet{
			Header: protocol.Output,
			Data:   data,
		}

		s.ptyOutput.Add(pkt)
		s.fanout(pkt, nil)
		s.logPacket(pkt)

		s.log.Debug("Pty Output", zap.ByteString("data", data))
	}
}
