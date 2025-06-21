package pseudo

import (
	"github.com/creack/pty"
	"log"
	"os"
	"os/exec"
)

func Spawn(cmd *exec.Cmd) (*os.File, error) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Failed to spawn pty. Error: ", err)
		return nil, err
	}

	return ptmx, nil
}
