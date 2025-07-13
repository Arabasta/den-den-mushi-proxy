package client

import (
	"go.uber.org/zap"
)

// ObserverReadLoop is for reading close messages from Observers ONLY
func (c *Connection) ObserverReadLoop() {
	for {
		if c.Ctx.Err() != nil {
			c.Log.Info("ObserverReadLoop: context done")
			return
		}

		_, _, err := c.Sock.ReadMessage()
		if err != nil {
			c.logWsError(err)
			if c.Close != nil {
				c.Close()
			} else {
				c.Log.Error("Error reading from observer connection, failed to call Close()", zap.Error(err))
			}
			return
		}
	}
}
