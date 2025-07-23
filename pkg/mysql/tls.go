package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
)

func registerMySQLTLSCA(caFile string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caFile)
	if err != nil {
		return fmt.Errorf("failed to read CA file: %w", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append CA certs")
	}

	return mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: false,
	})
}
