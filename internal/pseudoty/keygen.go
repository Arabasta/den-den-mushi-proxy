package pseudoty

import (
	"fmt"
	"github.com/charmbracelet/keygen"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

func GenerateEphemeralKey() (string, string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "ephemeral-ssh")
	if err != nil {
		return "", "", nil, err
	}

	keyPath := filepath.Join(tmpDir, "id_ed25519")
	keyPair, err := keygen.New(keyPath, keygen.WithKeyType(keygen.Ed25519), keygen.WithWrite())
	if err != nil {
		return "", "", nil, err
	}

	pubKeyString := string(ssh.MarshalAuthorizedKey(keyPair.PublicKey()))

	cleanup := func() { _ = os.RemoveAll(keyPath) }

	// todo: if dev mode
	fmt.Print(pubKeyString)
	fmt.Println(keyPath)
	fmt.Println(string(keyPair.RawPrivateKey()))

	return keyPath, pubKeyString, cleanup, nil
}
