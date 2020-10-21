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
		keyring:     fmt.Sprintf("/v1/transit/encrypt/%s", config.Key),
	}

	return newService, nil
}

// Service implements the vault encrypting value modifier.
type Service struct {
	vaultClient *vaultclient.Client
	keyring     string
}

type vaultResponse struct {
	Data vaultResponseData `json:"data"`
}

type vaultResponseData struct {
	Ciphertext string `json:"ciphertext"`
}

type vaultRequest struct {
	Plaintext string `json:"plaintext"`
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	request := s.vaultClient.NewRequest("POST", s.keyring)

	base64Encoded := base64.StdEncoding.EncodeToString(value)
	err := request.SetJSONBody(vaultRequest{Plaintext: base64Encoded})
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

	return []byte(content.Data.Ciphertext), nil
}
