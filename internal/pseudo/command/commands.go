package command

import (
	"den-den-mushi-Go/pkg/dto/connection"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

type Builder struct {
	log *zap.Logger
}

func NewBuilder(log *zap.Logger) *Builder {
	return &Builder{log: log}
}

func (b *Builder) BuildSshCmd(privateKeyPath string, c connection.Connection) *exec.Cmd {
	args := []string{
		"-i", privateKeyPath,
		"-p", c.Port,
		fmt.Sprintf("%s@%s", c.OSUser, c.ServerIP),
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=30",
		"-tt",
	}

	cmd := exec.Command("ssh", args...)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	b.log.Info("Created SSH command", zap.String("command", cmd.String()))
	return cmd
}

func (b *Builder) BuildBashCmd() *exec.Cmd {
	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	cmd.Dir = os.Getenv("HOME")

	b.log.Info("Created Bash command", zap.String("command", cmd.String()))
	return cmd
}
