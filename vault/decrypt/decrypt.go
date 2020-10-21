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
		keyring:     fmt.Sprintf("/v1/transit/decrypt/%s", config.Key),
	}

	return newService, nil
}

// Service implements the vault decrypting value modifier.
type Service struct {
	vaultClient *vaultclient.Client
	keyring     string
}

type vaultResponse struct {
	Data vaultResponseData `json:"data"`
}

type vaultResponseData struct {
	Plaintext string `json:"plaintext"`
}

type vaultRequest struct {
	Ciphertext string `json:"ciphertext"`
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	request := s.vaultClient.NewRequest("POST", s.keyring)
	err := request.SetJSONBody(vaultRequest{Ciphertext: string(value)})
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	response, err := s.vaultClient.RawRequest(request)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	if response.StatusCode != 200 {
		return []byte{}, microerror.Maskf(vaultResponseError, "expected 200 response, got %d", response.StatusCode)
	}

	content := vaultResponse{}
	err = response.DecodeJSON(&content)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	decrypted, err := base64.StdEncoding.DecodeString(content.Data.Plaintext)
	if err != nil {
		return []byte{}, microerror.Mask(err)
	}

	return decrypted, nil
}
