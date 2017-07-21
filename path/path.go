package path

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"
)

// Config represents the configuration used to create a new path service.
type Config struct {
	// Settings.
	JSONBytes []byte
	Separator string
}

// DefaultConfig provides a default configuration to create a new path service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		JSONBytes: nil,
		Separator: ".",
	}
}

// New creates a new configured path service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.JSONBytes == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.JSONBytes must not be empty")
	}
	if config.Separator == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Separator must not be empty")
	}

	var jsonMap map[string]interface{}
	{
		err := json.Unmarshal(config.JSONBytes, &jsonMap)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		// Internals.
		jsonMap: jsonMap,

		// Settings.
		jsonBytes: config.JSONBytes,
		separator: config.Separator,
	}

	return newService, nil
}

// Service implements the path service.
type Service struct {
	// Internals.
	jsonMap map[string]interface{}

	// Settings.
	jsonBytes []byte
	separator string
}

// All returns all paths found in the configured JSON structure.
func (s *Service) All() ([]string, error) {
	var paths []string
	{
		for k, v := range s.jsonMap {
			ps, err := s.allFromInterface(v)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
			if ps != nil {
				for _, p := range ps {
					paths = append(paths, pathWithKey(k, []string{p}, s.separator))
				}
			} else {
				paths = append(paths, pathWithKey(k, ps, s.separator))
			}
		}

		sort.Strings(paths)
	}

	return paths, nil
}

// Get returns the value found under the given path, if any.
func (s *Service) Get(path string) (interface{}, error) {
	var err error

	v := interface{}(s.jsonMap)

	for _, k := range strings.Split(path, s.separator) {
		v, err = s.getFromInterface(k, v)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	return v, nil
}

func (s *Service) JSONBytes() []byte {
	return s.jsonBytes
}

// Set changes the value of the given path.
func (s *Service) Set(path string, value interface{}) error {
	s.jsonMap = cast.ToStringMap(s.setFromInterface(s.jsonMap, path, value))

	b, err := json.MarshalIndent(s.jsonMap, "", "	")
	if err != nil {
		return err // TODO: Microerror
	}
	s.jsonBytes = b

	return nil
}

func (s *Service) setFromInterface(j interface{}, k string, v interface{}) interface{} {
	stringMap := cast.ToStringMap(j)
	keyPath := strings.Split(k, s.separator)

	if len(keyPath) == 1 {
		if _, ok := stringMap[k]; ok {
			stringMap[k] = v
			return stringMap
		} else {
			fmt.Println("key not found")
		}
	} else {
		if _, ok := stringMap[keyPath[0]]; ok {
			recursedKey := strings.Join(keyPath[1:], s.separator)
			modified := s.setFromInterface(stringMap[keyPath[0]], recursedKey, v)
			stringMap[keyPath[0]] = modified

			return stringMap
		} else {
			fmt.Println("key not found")
		}
	}

	return nil
}

func (s *Service) allFromInterface(value interface{}) ([]string, error) {
	var paths []string

	// process map
	{
		stringMap, err := cast.ToStringMapE(value)
		if err != nil {
			// fall through
		} else {
			for k, v := range stringMap {
				ps, err := s.allFromInterface(v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				paths = append(paths, pathWithKey(k, ps, s.separator))
			}
		}
	}

	// process slice
	if len(paths) == 0 {
		slice, err := cast.ToSliceE(value)
		if err != nil {
			// fall through
		} else {
			for i, v := range slice {
				ps, err := s.allFromInterface(v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				if ps != nil {
					paths = append(paths, pathWithKey(fmt.Sprintf("[%d]", i), ps, s.separator))
				}
			}
		}
	}

	return paths, nil
}

func (s *Service) getFromInterface(key string, value interface{}) (interface{}, error) {
	var newValue interface{}

	// process map
	{
		stringMap, err := cast.ToStringMapE(value)
		if err != nil {
			// fall through
		} else {
			for k, v := range stringMap {
				if k != key {
					continue
				}

				newValue, err = s.getFromInterface(k, v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				break
			}
		}
	}

	// process slice
	if newValue == nil {
		slice, err := cast.ToSliceE(value)
		if err != nil {
			// fall through
		} else {
			for i, v := range slice {
				k := fmt.Sprintf("[%d]", i)

				if k != key {
					continue
				}

				newValue, err = s.getFromInterface(k, v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				break
			}
		}
	}

	// value is neither map nor slice
	if newValue == nil {
		newValue = value
	}

	return newValue, nil
}

func pathWithKey(key string, paths []string, separator string) string {
	return strings.Join(append([]string{key}, paths...), separator)
}
