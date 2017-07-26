package inspect

// Config represents the configuration used to create a new inspect value
// modifier.
type Config struct {
}

// DefaultConfig provides a default configuration to create a new inspect value
// modifier by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured inspect value modifier.
func New(config Config) (*Service, error) {
	newService := &Service{
		// Internals.
		inspectedValue: nil,
	}

	return newService, nil
}

// Service implements the inspect value modifier.
type Service struct {
	// Internals.
	inspectedValue []byte
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	s.inspectedValue = value
	return value, nil
}

func (s *Service) InspectedValue() []byte {
	return s.inspectedValue
}
