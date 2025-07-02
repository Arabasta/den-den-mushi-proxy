package pseudotty

import (
	"den-den-mushi-Go/internal/protocol"
)

type Purpose interface {
	HandleInput(s *Session, pkt protocol.Packet) (string, error)
	HandleOther(s *Session, pkt protocol.Packet) (string, error)
}
