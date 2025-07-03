package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
)

type Purpose interface {
	HandleInput(s *Session, pkt protocol.Packet) (string, error)
	HandleOther(s *Session, pkt protocol.Packet) (string, error)
}
