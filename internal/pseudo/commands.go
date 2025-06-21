package pseudo

import (
	"den-den-mushi-Go/pkg/connection"
	"fmt"
	"os"
	"os/exec"
)

func BuildSshCmd(privateKeyPath string, c connection.Connection) *exec.Cmd {
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
	fmt.Println(args)
	cmd := exec.Command("ssh", args...)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	return cmd
}

func BuildBashCmd() *exec.Cmd {
	cmd := exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	cmd.Dir = os.Getenv("HOME")
	return cmd
}
