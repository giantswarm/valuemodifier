package decrypt

import (
	"encoding/base64"
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
)

// Config represents the configuration used to create a new vault decrypting
// value modifier.
type Config struct {
	VaultClient *vaultclient.Client
	Key         string
}

// DefaultConfig provides a default configuration to create a new vault
// decrypting value modifier by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured vault decrypting value modifier.
func New(config Config) (*Service, error) {
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultClient must be defined")
	}
	if config.Key == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Key must be defined")

	}
	newService := &Service{
		vaultClient: config.VaultClient,
		path:        fmt.Sprintf("/transit/decrypt/%s", config.Key),
	}

	return newService, nil
}

// Service implements the vault decrypting value modifier.
type Service struct {
	vaultClient *vaultclient.Client
	path        string
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	plainText, err := s.Decrypt(value)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	decrypted, err := base64.StdEncoding.DecodeString(plainText)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	return decrypted, nil
}

func (s *Service) Decrypt(cipherText []byte) (string, error) {
	secret, err := s.vaultClient.Logical().Write(s.path, map[string]interface{}{
		"ciphertext": string(cipherText),
	})

	if err != nil {
		return "", microerror.Mask(err)
	}

	return fmt.Sprintf("%v", secret.Data["plaintext"]), nil
}
