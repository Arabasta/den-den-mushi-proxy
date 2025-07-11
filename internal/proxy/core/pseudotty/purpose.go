package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/dto"
	"errors"
	"go.uber.org/zap"
)

type Purpose interface {
	HandleInput(s *Session, pkt protocol.Packet) (string, error)
	HandleOther(s *Session, pkt protocol.Packet) (string, error)
}

func setPurpose(s *Session, purpose dto.ConnectionPurpose) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch purpose {
	case dto.Change:
		s.log.Info("Setting session purpose to Change")
		s.purpose = &ChangeRequestPurpose{}
	case dto.Healthcheck:
		s.log.Info("Setting session purpose to Healthcheck")
		s.purpose = &HealthcheckPurpose{}
	default:
		s.log.Error("Unknown purpose for new connection", zap.String("purpose", string(purpose)))
		return errors.New("unknown purpose")
	}

	return nil
}
