package overwrite

import (
	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a new overwrite
// value modifier.
type Config struct {
	// Settings.

	// NewValue is the new piece of information used to overwrite the traversed
	// existing value.
	NewValue string
}

// DefaultConfig provides a default configuration to create a new overwrite
// decoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		NewValue: "",
	}
}

// New creates a new configured overwrite value modifier.
func New(config Config) (*Service, error) {
	// Settings.
	if config.NewValue == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.NewValue must not be empty")
	}

	newService := &Service{
		newValue: config.NewValue,
	}

	return newService, nil
}

// Service implements the overwrite value modifier.
type Service struct {
	// Settings.
	newValue string
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	return []byte(s.newValue), nil
}
