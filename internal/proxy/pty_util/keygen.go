package pty_util

import (
	"den-den-mushi-Go/internal/proxy/config"
	"github.com/charmbracelet/keygen"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

// t odo: put in /var/tmp
func GenerateEphemeralKey(cfg *config.Config, log *zap.Logger) (string, string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "ephemeral-ssh")
	if err != nil {
		log.Error("Failed to create temporary directory for ephemeral SSH key", zap.Error(err))
		return "", "", nil, err
	}

	keyPath := filepath.Join(tmpDir, "id_ed25519")
	keyPair, err := keygen.New(keyPath, keygen.WithKeyType(keygen.Ed25519), keygen.WithWrite())
	if err != nil {
		log.Error("Failed to generate ephemeral SSH key", zap.Error(err), zap.String("keyPath", keyPath))
		return "", "", nil, err
	}

	pubKeyString := strings.Trim(string(ssh.MarshalAuthorizedKey(keyPair.PublicKey())), "\n")
	pubKeyString += " " + cfg.Development.SshPubKeyHostnameSuffix

	cleanup := func() {
		err = os.RemoveAll(keyPath)
		if err != nil {
			log.Error("Failed to remove ephemeral SSH key", zap.Error(err), zap.String("keyPath", keyPath))
			return
		}
		log.Debug("Ephemeral SSH key removed", zap.String("keyPath", keyPath))
	}

	// todo: remove after testing
	if cfg.App.Environment == "dev" {
		log.Debug("Ephemeral SSH key generated",
			zap.String("keyPath", keyPath),
			zap.String("publicKey", pubKeyString),
			zap.String("privateKey", string(keyPair.RawPrivateKey())))
	}

	log.Debug("Ephemeral SSH key generated")
	return keyPath, pubKeyString, cleanup, nil
}
