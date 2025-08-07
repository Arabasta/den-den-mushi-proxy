package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

// todo: move to purpose package

type Purpose interface {
	HandleInput(s *Session, pkt protocol.Packet) (string, error)
	HandleOther(s *Session, pkt protocol.Packet) (string, error)
}

func setPurpose(s *Session, purpose types.ConnectionPurpose) error {
	s.log.Info("Setting purpose for connection", zap.String("purpose", string(purpose)))

	switch purpose {
	case types.Change:
		s.log.Debug("Setting session purpose to Change")
		s.purpose = &ChangeRequestPurpose{}
	case types.Healthcheck:
		s.log.Debug("Setting session purpose to Healthcheck")
		s.purpose = &HealthcheckPurpose{}
	case types.IExpress:
		// todo: implement IExpressPurpose
		s.log.Debug("Iexpress purpose...Setting session purpose to Healthcheck for now")
		s.purpose = &HealthcheckPurpose{}
	default:
		s.log.Error("Unknown purpose for new connection", zap.String("purpose", string(purpose)))
		return errors.New("unknown purpose")
	}

	return nil
}
