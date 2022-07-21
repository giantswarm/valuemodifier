package encrypt

import (
	"encoding/base64"
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
)

// Config represents the configuration used to create a new vault encrypting
// value modifier.
type Config struct {
	VaultClient *vaultclient.Client
	Key         string
}

// DefaultConfig provides a default configuration to create a new vault
// encrypting value modifier by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured vault encrypting value modifier.
func New(config Config) (*Service, error) {
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultClient must be defined")
	}
	if config.Key == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Key must be defined")

	}
	newService := &Service{
		vaultClient: config.VaultClient,
		path:        fmt.Sprintf("/transit/encrypt/%s", config.Key),
	}

	return newService, nil
}

// Service implements the vault encrypting value modifier.
type Service struct {
	vaultClient *vaultclient.Client
	path        string
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	base64Encoded := base64.StdEncoding.EncodeToString(value)

	cipherText, err := s.encrypt(base64Encoded)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	return []byte(cipherText), nil
}

func (s *Service) encrypt(plainText string) (string, error) {
	secret, err := s.vaultClient.Logical().Write(s.path, map[string]interface{}{
		"plaintext": plainText,
	})

	if err != nil {
		return "", microerror.Mask(err)
	}

	return fmt.Sprintf("%v", secret.Data["ciphertext"]), nil
}
