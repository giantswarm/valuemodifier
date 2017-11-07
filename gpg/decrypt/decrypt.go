package decrypt

import (
	"bytes"
	"io/ioutil"

	"github.com/giantswarm/microerror"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// Config represents the configuration used to create a new GPG decryption
// value modifier.
type Config struct {
	// Settings.

	// Pass is the passphrase used to decrypt GPG messages.
	Pass string
}

// DefaultConfig provides a default configuration to create a new GPG decryption
// decoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		Pass: "",
	}
}

// New creates a new configured GPG decryption value modifier.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Pass == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Pass must not be empty")
	}

	newService := &Service{
		pass: config.Pass,
	}

	return newService, nil
}

// Service implements the GPG decryption value modifier.
type Service struct {
	// Settings.
	pass string
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	if len(value) == 0 {
		return value, nil
	}

	buf := bytes.NewBuffer(value)
	decoder, err := armor.Decode(buf)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	promptFunc := func() func([]openpgp.Key, bool) ([]byte, error) {
		retried := false
		return func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			if !retried {
				retried = true
				return []byte(s.pass), nil
			}
			return nil, microerror.Maskf(wrongGPGPasswordError, "Decryption failed with given GPG password")
		}
	}()
	details, err := openpgp.ReadMessage(decoder.Body, nil, promptFunc, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	b, err := ioutil.ReadAll(details.UnverifiedBody)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return b, nil
}
