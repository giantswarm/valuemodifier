package path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v1"

	yamltojson "github.com/ghodss/yaml"
	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"
)

// Config represents the configuration used to create a new path service.
type Config struct {
	// Settings.
	InputBytes []byte
	Separator  string
}

// DefaultConfig provides a default configuration to create a new path service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		InputBytes: nil,
		Separator:  ".",
	}
}

// New creates a new configured path service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.InputBytes == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.InputBytes must not be empty")
	}
	if config.Separator == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Separator must not be empty")
	}

	var err error

	var isJSON bool
	var jsonBytes []byte
	var jsonStructure interface{}
	{
		jsonBytes, isJSON, err = toJSON(config.InputBytes)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		err := json.Unmarshal(jsonBytes, &jsonStructure)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		// Internals.
		isJSON:        isJSON,
		jsonBytes:     jsonBytes,
		jsonStructure: jsonStructure,

		// Settings.
		separator: config.Separator,
	}

	return newService, nil
}

// Service implements the path service.
type Service struct {
	// Internals.
	isJSON        bool
	jsonBytes     []byte
	jsonStructure interface{}

	// Settings.
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
	value, err := s.getFromInterface(path, s.jsonStructure)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return value, nil
}

func (s *Service) OutputBytes() ([]byte, error) {
	b := s.jsonBytes
	if !s.isJSON {
		var err error
		b, err = yamltojson.JSONToYAML(b)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	return b, nil
}

// Set changes the value of the given path.
func (s *Service) Set(path string, value interface{}) error {
	var err error

	s.jsonStructure, err = s.setFromInterface(path, value, s.jsonStructure)
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

func (s *Service) Validate(paths []string) error {
	all, err := s.All()
	if err != nil {
		return microerror.MaskAny(err)
	}

	for _, p := range paths {
		if containsString(all, p) {
			continue
		}

		return microerror.MaskAnyf(notFoundError, "path '%s' not found", p)
	}

	return nil
}

func (s *Service) allFromInterface(value interface{}) ([]string, error) {
	// process map
	{
		stringMap, err := cast.ToStringMapE(value)
		if err != nil {
			// fall through
		} else {
			var paths []string

			for k, v := range stringMap {
				ps, err := s.allFromInterface(v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				paths = append(paths, pathWithKey(k, ps, s.separator))
			}

			return paths, nil
		}
	}

	// process slice
	{
		slice, err := cast.ToSliceE(value)
		if err != nil {
			// fall through
		} else {
			var paths []string

			for i, v := range slice {
				ps, err := s.allFromInterface(v)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}
				if ps != nil {
					paths = append(paths, pathWithKey(fmt.Sprintf("[%d]", i), ps, s.separator))
				}
			}

			return paths, nil
		}
	}

	// process string
	{
		str, err := cast.ToStringE(value)
		if err != nil {
			// fall through
		} else {
			jsonBytes, _, err := toJSON([]byte(str))
			if err != nil {
				// fall through
			} else {
				var jsonStructure interface{}
				err := json.Unmarshal(jsonBytes, &jsonStructure)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				ps, err := s.allFromInterface(jsonStructure)
				if err != nil {
					return nil, microerror.MaskAny(err)
				}

				return ps, nil
			}
		}
	}

	return nil, nil
}

func (s *Service) getFromInterface(path string, jsonStructure interface{}) (interface{}, error) {
	// process map
	{
		stringMap, err := cast.ToStringMapE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			split := strings.Split(path, s.separator)

			value, ok := stringMap[split[0]]
			if ok {
				if len(split) == 1 {
					return value, nil
				} else {
					recPath := strings.Join(split[1:], s.separator)

					value, err := s.getFromInterface(recPath, value)
					if err != nil {
						return nil, microerror.MaskAny(err)
					}

					return value, nil
				}
			} else {
				return nil, microerror.MaskAnyf(notFoundError, "key '%s' not found in path", path)
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
				return nil, microerror.MaskAnyf(notFoundError, "key '%s' not found in path", split[0])
			}
			recPath := strings.Join(split[1:], s.separator)

			value, err := s.getFromInterface(recPath, slice[index])
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			return value, nil
		}
	}

	return nil, nil
}

func (s *Service) setFromInterface(path string, value interface{}, jsonStructure interface{}) (interface{}, error) {
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
					return nil, microerror.MaskAnyf(notFoundError, "key '%s' not found in path", path)
				}
			} else {
				_, ok := stringMap[split[0]]
				if ok {
					recPath := strings.Join(split[1:], s.separator)

					modified, err := s.setFromInterface(recPath, value, stringMap[split[0]])
					if err != nil {
						return nil, microerror.MaskAny(err)
					}
					stringMap[split[0]] = modified

					return stringMap, nil
				} else {
					return nil, microerror.MaskAnyf(notFoundError, "key '%s' not found in path", path)
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
				return nil, microerror.MaskAnyf(notFoundError, "key '%s' not found in path", split[0])
			}
			recPath := strings.Join(split[1:], s.separator)

			modified, err := s.setFromInterface(recPath, value, slice[index])
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
			slice[index] = modified

			return slice, nil
		}
	}

	return nil, nil
}

func containsString(list []string, item string) bool {
	for _, l := range list {
		if l == item {
			return true
		}
	}

	return false
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

func isJSON(b []byte) bool {
	var l []interface{}
	isList := json.Unmarshal(b, &l) == nil

	var m map[string]interface{}
	isObject := json.Unmarshal(b, &m) == nil

	return isObject || isList
}

func isYAMLList(b []byte) bool {
	var l []interface{}
	return yaml.Unmarshal(b, &l) == nil && bytes.HasPrefix(b, []byte("-"))
}

func isYAMLObject(b []byte) bool {
	var m map[interface{}]interface{}
	return yaml.Unmarshal(b, &m) == nil && !bytes.HasPrefix(b, []byte("-"))
}

func pathWithKey(key string, paths []string, separator string) string {
	return strings.Join(append([]string{key}, paths...), separator)
}

func toJSON(b []byte) ([]byte, bool, error) {
	if isJSON(b) {
		return b, true, nil
	}

	isYAMLList := isYAMLList(b)
	isYAMLObject := isYAMLObject(b)

	var jsonBytes []byte
	if isYAMLList && !isYAMLObject {
		var jsonList []interface{}
		err := yamltojson.Unmarshal(b, &jsonList)
		if err != nil {
			return nil, false, microerror.MaskAny(err)
		}

		jsonBytes, err = json.Marshal(jsonList)
		if err != nil {
			return nil, false, microerror.MaskAny(err)
		}

		return jsonBytes, false, nil
	}

	if !isYAMLList && isYAMLObject {
		var jsonMap map[string]interface{}
		err := yamltojson.Unmarshal(b, &jsonMap)
		if err != nil {
			return nil, false, microerror.MaskAny(err)
		}

		jsonBytes, err = json.Marshal(jsonMap)
		if err != nil {
			return nil, false, microerror.MaskAny(err)
		}

		return jsonBytes, false, nil
	}

	return nil, false, microerror.MaskAny(invalidFormatError)
}
