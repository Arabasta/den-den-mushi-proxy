package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

type SshVaultKeyConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

func getFromVault() (string, error) {
	return "", nil
}

func (c *SshVaultKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, *exec.Cmd, error) {
	// todo: implement vault key retrieval
	return nil, nil, errors.New("not implemented")
	//keyPath, err := getFromVault()
	//if err != nil {
	//	c.log.Error("failed to generate ephemeral key", zap.Error(err))
	//	return nil, nil, err
	//}
	//
	//var pty *os.File
	//var cmd *exec.Cmd
	//for i := 0; i <= c.cfg.Pty.SpawnRetryCount; i++ {
	//	c.log.Debug("Spawning pseudo terminal for SSH connection", zap.Int("attempt", i+1))
	//	cmd := c.commandBuilder.BuildSshCmd(keyPath, claims.Connection.Server)
	//
	//	pty, cmd, err = pty_util.Spawn(cmd, c.log)
	//	if err == nil {
	//		c.log.Info("Pseudo terminal spawned", zap.String("keyPath", keyPath))
	//		break
	//	}
	//
	//	time.Sleep(c.cfg.Pty.SpawnRetryIntervalSeconds * time.Second)
	//}
	//if err != nil {
	//	// cleanup
	//	return nil, nil, err
	//}
	//
	//go func() {
	//	time.Sleep(60 * time.Second)
	//	// cleanup()
	//}()
	//
	//return pty, cmd, nil
}
