package encrypt

import (
	"bytes"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/giantswarm/microerror"
)

// Config represents the configuration used to create a new GPG encryption
// value modifier.
type Config struct {
	// Settings.

	// Pass is the passphrase used to encrypt GPG messages.
	Pass string
}

// DefaultConfig provides a default configuration to create a new GPG encryption
// encoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		Pass: "",
	}
}

// New creates a new configured GPG encryption value modifier.
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

// Service implements the GPG encryption value modifier.
type Service struct {
	// Settings.
	pass string
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	if len(value) == 0 {
		return value, nil
	}

	buf := bytes.NewBuffer(nil)
	encoder, err := armor.Encode(buf, openpgp.SignatureType, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	encrypter, err := openpgp.SymmetricallyEncrypt(encoder, []byte(s.pass), nil, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	_, err = encrypter.Write(value)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	encrypter.Close()
	encoder.Close()

	return buf.Bytes(), nil
}
