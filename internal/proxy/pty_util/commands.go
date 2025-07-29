package pty_util

import (
	"den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/dto"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

type Builder struct {
	log *zap.Logger
	cfg *config.Ssh
}

func NewBuilder(log *zap.Logger, cfg *config.Ssh) *Builder {
	return &Builder{log: log, cfg: cfg}
}

func (b *Builder) BuildSshCmd(privateKeyPath string, s dto.ServerInfo) *exec.Cmd {
	args := []string{
		"-i", privateKeyPath,
		"-p", s.Port,
		fmt.Sprintf("%s@%s", s.OSUser, s.IP),
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=30",
		"-tt",
	}

	cmd := exec.Command(b.cfg.SshCommand, args...)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	b.log.Debug("Created SSH command", zap.String("command", cmd.String()))
	return cmd
}

func (b *Builder) BuildBashCmd() *exec.Cmd {
	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	cmd.Dir = os.Getenv("HOME")

	b.log.Info("Created Bash command", zap.String("command", cmd.String()))
	return cmd
}
