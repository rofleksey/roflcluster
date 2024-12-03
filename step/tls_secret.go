package step

import (
	"encoding/base64"
	"fmt"
	"github.com/melbahja/goph"
	"os"
	"roflcluster/config"
)

type TLSSecretStep struct {
}

type TLSData struct {
	Cert string
	Key  string
	CA   string
}

func (s *TLSSecretStep) Execute(client *goph.Client, cfg *config.Config) error {
	certBytes, err := os.ReadFile("certs/tls.crt")
	if err != nil {
		return fmt.Errorf("failed to read tls.crt: %w", err)
	}

	keyBytes, err := os.ReadFile("certs/tls.key")
	if err != nil {
		return fmt.Errorf("failed to read tls.key: %w", err)
	}

	caBytes, err := os.ReadFile("certs/ca.crt")
	if err != nil {
		return fmt.Errorf("failed to read ca.crt: %w", err)
	}

	data := TLSData{
		Cert: base64.StdEncoding.EncodeToString(certBytes),
		Key:  base64.StdEncoding.EncodeToString(keyBytes),
		CA:   base64.StdEncoding.EncodeToString(caBytes),
	}

	return FormatApplyTemplate(client, "templates/tls-secret.yaml.tmpl", &data)
}

func (s *TLSSecretStep) String() string {
	return "Install TLS Secret"
}
