package path

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

	var jsonStructure interface{}
	{
		err := json.Unmarshal(config.JSONBytes, &jsonStructure)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		// Internals.
		jsonStructure: jsonStructure,

		// Settings.
		jsonBytes: config.JSONBytes,
		separator: config.Separator,
	}

	return newService, nil
}

// Service implements the path service.
type Service struct {
	// Internals.
	jsonStructure interface{}

	// Settings.
	jsonBytes []byte
	separator string
}

// All returns all paths found in the configured JSON structure.
func (s *Service) All() ([]string, error) {
	var paths []string
	{
		for k, v := range cast.ToStringMap(s.jsonStructure) {
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

	v := s.jsonStructure

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
	var err error

	s.jsonStructure, err = s.setFromInterface(s.jsonStructure, path, value)
	if err != nil {
		return microerror.MaskAny(err)
	}

	b, err := json.MarshalIndent(s.jsonStructure, "", "  ")
	if err != nil {
		return microerror.MaskAny(err)
	}
	s.jsonBytes = b

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

func (s *Service) setFromInterface(jsonStructure interface{}, path string, value interface{}) (interface{}, error) {
	// process map
	{
		stringMap, err := cast.ToStringMapE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			split := strings.Split(path, s.separator)

			if len(split) == 1 {
				_, ok := stringMap[path]
				if ok {
					stringMap[path] = value
					return stringMap, nil
				} else {
					return nil, microerror.MaskAnyf(pathNotFoundError, path)
				}
			} else {
				_, ok := stringMap[split[0]]
				if ok {
					recursedKey := strings.Join(split[1:], s.separator)

					modified, err := s.setFromInterface(stringMap[split[0]], recursedKey, value)
					if err != nil {
						return nil, microerror.MaskAny(err)
					}
					stringMap[split[0]] = modified

					return stringMap, nil
				} else {
					return nil, microerror.MaskAnyf(pathNotFoundError, path)
				}
			}
		}
	}

	// process slice
	{
		slice, err := cast.ToSliceE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			split := strings.Split(path, s.separator)

			index, err := indexFromKey(split[0])
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			if index >= len(slice) {
				return nil, microerror.MaskAnyf(pathNotFoundError, split[0])
			}
			recursedKey := strings.Join(split[1:], s.separator)

			modified, err := s.setFromInterface(slice[index], recursedKey, value)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
			slice[index] = modified

			return slice, nil
		}
	}

	return nil, nil
}

func indexFromKey(key string) (int, error) {
	re := regexp.MustCompile("\\[[0-9]\\]")
	ok := re.MatchString(key)
	if !ok {
		return 0, microerror.MaskAnyf(keyNotIndexError, key)
	}

	s := key[1 : len(key)-1]
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, microerror.MaskAny(err)
	}

	return i, nil
}

func pathWithKey(key string, paths []string, separator string) string {
	return strings.Join(append([]string{key}, paths...), separator)
}
