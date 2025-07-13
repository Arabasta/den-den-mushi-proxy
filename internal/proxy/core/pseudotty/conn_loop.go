package pseudotty

func (s *Session) connLoop() {
	for {
		select {
		case c := <-s.connRegisterCh:
			s.addConn(c)
		case c := <-s.connDeregisterCh:
			s.removeConn(c)
		case <-s.ctx.Done():
			s.log.Info("connLoop: context done")
			return
		}
	}
}
