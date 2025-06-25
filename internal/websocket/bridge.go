package websocket

import (
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"os"
	"time"
)

func (s *Service) bridge(ws *websocket.Conn, pty io.ReadWriteCloser, claims *token.Claims) {
	s.closeAll(ws, pty)

	// temp for demo
	logPath := "./log/" + claims.Connection.ChangeID + "_" + claims.Connection.ServerIP + "_" + claims.Connection.OSUser + "_" + string(claims.Connection.Purpose) + ".log"
	logFile, err := os.OpenFile(logPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.log.Error("Failed to open log file", zap.Error(err))
		return
	}
	s.writeLogHeader(logFile, claims)
	defer logFile.Close()

	//var missedPongs int32
	//
	//go s.handlePing(ws, &missedPongs)
	//s.handlePong(ws, &missedPongs)

	go s.handleOutput(ws, pty, logFile)
	s.handleInput(ws, pty, claims, logFile)
}

func (s *Service) writeLogHeader(logFile *os.File, claims *token.Claims) {
	header := "User: " + claims.Subject + "\n" +
		"ChangeID: " + claims.Connection.ChangeID + "\n" +
		"ServerIP: " + claims.Connection.ServerIP + "\n" +
		"OSUser: " + claims.Connection.OSUser + "\n" +
		"Purpose: " + string(claims.Connection.Purpose) + "\n\n" +
		"Log Start Time: " + time.Now().Format(time.RFC3339) + "\n\n"
	if _, err := logFile.WriteString(header); err != nil {
		s.log.Error("Failed to write log header", zap.Error(err))
	}
}

func (s *Service) tempLogFunction(logFile *os.File, data []byte) {
	if logFile == nil {
		return
	}

	_, err := logFile.Write(data)
	if err != nil {
		s.log.Error("Failed to write to log file", zap.Error(err))
	}
}
