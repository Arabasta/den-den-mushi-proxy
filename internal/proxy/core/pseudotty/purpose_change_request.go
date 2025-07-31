package pseudotty

import (
	"bytes"
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type ChangeRequestPurpose struct{}

func (p *ChangeRequestPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	if pkt.Header != protocol.Input {
		s.log.Error("Invalid input packet", zap.Any("pkt", pkt))
		return "", fmt.Errorf("expected Input header, got %s", string(pkt.Header))
	}

	for _, b := range pkt.Data {
		if bytes.Equal([]byte{b}, constants.Enter) {
			return p.handleEnter(s, pkt)
		}

		return handler.Get[protocol.Input].Handle(pkt, s.pty)
	}
	return "", nil
}

func (p *ChangeRequestPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.pty)
}

func (p *ChangeRequestPurpose) handleEnter(s *Session, pkt protocol.Packet) (string, error) {
	s.log.Debug("Handling Enter key press. Checking command: ", zap.String("line", s.line.String()))

	s.log.Debug("Calling filter.IsValid for command", zap.String("line", s.line.String()))
	msg, allowed := s.filter.IsValid(s.line.String(), "")
	s.log.Debug("Filter result", zap.String("line", s.line.String()), zap.Bool("allowed", allowed), zap.String("message", msg))
	if !allowed {
		errPkt := protocol.Packet{Header: protocol.BlockedCommand, Data: []byte(s.line.String())}
		s.ptyOutput.Add(errPkt)
		s.fanout(errPkt, nil)

		// send Ctrl +C to clear pty
		ctrlCPkt := protocol.Packet{Header: protocol.Input, Data: constants.CtrlC}
		_, err := handler.Get[protocol.Input].Handle(ctrlCPkt, s.pty)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("\n%s [Command Blocked] %s", time.Now().Format(time.TimeOnly), msg), nil
	}

	s.log.Debug("Checked command", zap.String("line", s.line.String()), zap.Bool("allowed", allowed))
	return handler.Get[protocol.Input].Handle(pkt, s.pty)
}
