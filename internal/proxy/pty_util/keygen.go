package pty_util

import (
	"den-den-mushi-Go/internal/proxy/config"
	"github.com/charmbracelet/keygen"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

func GenerateEphemeralKey(cfg *config.Config, log *zap.Logger) (string, string, func(), error) {
	err := os.MkdirAll(cfg.Ssh.EphemeralKeyPath, 0700)
	if err != nil {
		log.Error("Failed to create directory for ephemeral SSH key", zap.String("path", cfg.Ssh.EphemeralKeyPath),
			zap.Error(err))
		return "", "", nil, err
	}

	var keyPath string
	var keyPair *keygen.KeyPair
	uniqueSuffix := uuid.NewString()

	if cfg.Ssh.IsRSAKeyPair {
		log.Debug("Generating RSA key pair for ephemeral SSH key")
		keyPath = filepath.Join(cfg.Ssh.EphemeralKeyPath, "id_rsa_"+uniqueSuffix)
		keyPair, err = keygen.New(keyPath, keygen.WithKeyType(keygen.RSA), keygen.WithWrite())
	} else {
		log.Debug("Generating Ed25519 key pair for ephemeral SSH key")
		keyPath = filepath.Join(cfg.Ssh.EphemeralKeyPath, "id_ed25519_"+uniqueSuffix)
		keyPair, err = keygen.New(keyPath, keygen.WithKeyType(keygen.Ed25519), keygen.WithWrite())
	}

	if err != nil {
		log.Error("Failed to generate ephemeral SSH key", zap.Error(err), zap.String("keyPath", keyPath))
		return "", "", nil, err
	}

	pubKeyString := strings.Trim(string(ssh.MarshalAuthorizedKey(keyPair.PublicKey())), "\n")
	pubKeyString += " " + cfg.Ssh.PubKeyHostnameSuffix

	cleanup := func() {
		if cfg.Ssh.IsCleanupEnabled {
			os.Remove(keyPath)
			os.Remove(keyPath + ".pub")
			log.Debug("Ephemeral SSH key pair removed", zap.String("keyPath", keyPath))
		}
	}

	if cfg.Ssh.IsLogPrivateKey {
		log.Debug("Ephemeral SSH key generated",
			zap.String("keyPath", keyPath),
			zap.String("publicKey", pubKeyString),
			zap.String("privateKey", string(keyPair.RawPrivateKey())))
	}

	log.Debug("Ephemeral SSH key generated")
	return keyPath, pubKeyString, cleanup, nil
}
