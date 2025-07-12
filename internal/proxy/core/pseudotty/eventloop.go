package pseudotty

// connEventLoop should be the only goroutine mutating primary and observers, other code needs to lock if reading
func (s *Session) connEventLoop() {
	for {
		select {
		case c := <-s.connRegisterCh:
			s.addConn(c)
		case c := <-s.connDeregisterCh:
			s.removeConn(c)
		}
	}
}
