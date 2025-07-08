package pseudotty

// eventLoop should be the only goroutine mutating primary and observers, other code needs to lock if reading
func (s *Session) eventLoop() {
	for {
		select {
		case pkt, ok := <-s.outboundCh:
			if !ok {
				return
			}
			s.logLine(pkt.Header, string(pkt.Data))
			s.fanout(pkt)
		case c := <-s.connRegisterCh:
			s.addConn(c)
		case c := <-s.connDeregisterCh:
			s.removeConn(c)
		}
	}
}
