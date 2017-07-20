package valuemodifier

import (
	"encoding/json"
	"strings"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v1"
)

// Config represents the configuration used to create a new value modifier
// traverser.
type Config struct {
	// Dependencies.
	ValueModifiers []ValueModifier

	// Settings.
	IgnoreFields []string
	SelectFields []string
}

// DefaultConfig provides a default configuration to create a new GPG decryption
// decoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		ValueModifiers: nil,

		// Settings.
		IgnoreFields: nil,
		SelectFields: nil,
	}
}

// New creates a new configured GPG decryption value modifier.
func New(config Config) (*Service, error) {
	// Dependencies.
	if len(config.ValueModifiers) == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.ValueModifiers must not be empty")
	}

	// Settings.
	if len(config.IgnoreFields) != 0 && len(config.SelectFields) != 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.IgnoreFields must be empty when config.SelectFields provided")
	}

	newService := &Service{
		// Dependencies.
		valueModifiers: config.ValueModifiers,

		// Settings.
		ignoreFields: config.IgnoreFields,
		selectFields: config.SelectFields,
	}

	return newService, nil
}

// Service implements the traversing value modifier.
type Service struct {
	// Dependencies.
	valueModifiers []ValueModifier

	// Settings.
	ignoreFields []string
	selectFields []string
}

func (s *Service) TraverseJSON(input []byte) ([]byte, error) {
	var m map[string]interface{}
	err := json.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var p []string
	for k, v := range m {
		m[k] = s.toModifiedValueJSON(k, v, p)
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func (s *Service) TraverseYAML(input []byte) ([]byte, error) {
	var m map[interface{}]interface{}
	err := yaml.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var p []string
	for k, v := range m {
		m[k] = s.toModifiedValueYAML(k, v, p)
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func (s *Service) toModifiedValueJSON(key string, val interface{}, path []string) interface{} {
	var p []string
	p = append(p, path...)

	m1, ok := val.(map[string]interface{})
	if ok {
		for k, v := range m1 {
			m1[k] = s.toModifiedValueJSON(k, v, p)
		}

		return m1
	}

	m2, ok := val.([]interface{})
	if ok {
		for i, v := range m2 {
			m2[i] = s.toModifiedValueJSON("", v, p)
		}

		return m2
	}

	str := cast.ToString(val)
	if str != "" {
		for _, f := range s.ignoreFields {
			if f == key {
				return str
			}
		}
		for _, m := range s.valueModifiers {
			o, err := m.Modify([]byte(str))
			if err != nil {
				panic(err)
			}
			str = string(o)
		}
	}

	return str
}

func (s *Service) toModifiedValueYAML(key interface{}, val interface{}, path []string) interface{} {
	var p []string
	p = append(p, path...)

	m1, ok := val.(map[interface{}]interface{})
	if ok {
		for k, v := range m1 {
			m1[k] = s.toModifiedValueYAML(k, v, p)
		}

		return m1
	}

	m2, ok := val.([]interface{})
	if ok {
		for i, v := range m2 {
			m2[i] = s.toModifiedValueYAML("", v, p)
		}

		return m2
	}

	str := cast.ToString(val)
	if str != "" {
		var m3 map[interface{}]interface{}
		err := yaml.Unmarshal([]byte(str), &m3)
		if err != nil || m3 == nil {
			if len(s.selectFields) == 0 {
				for _, f := range s.ignoreFields {
					if f == key {
						return str
					}
				}
			} else {
				for _, f := range s.selectFields {
					if f == strings.Join(append(p, key.(string)), ".") {
						for _, m := range s.valueModifiers {
							o, err := m.Modify([]byte(str))
							if err != nil {
								panic(err)
							}
							str = string(o)
						}
					}
				}
			}

			if len(s.selectFields) == 0 && len(s.selectFields) == 0 {
				for _, m := range s.valueModifiers {
					o, err := m.Modify([]byte(str))
					if err != nil {
						panic(err)
					}
					str = string(o)
				}
			}
		} else {
			var m4 map[string]interface{}
			err := json.Unmarshal([]byte(str), &m4)
			if err != nil || m4 == nil {
				for k, v := range m3 {
					m3[k] = s.toModifiedValueYAML(k, v, p)
				}

				b, err := yaml.Marshal(m3)
				if err != nil {
					panic(err)
				}

				return string(b)
			} else {
				for k, v := range m4 {
					m4[k] = s.toModifiedValueJSON(k, v, p)
				}

				b, err := json.MarshalIndent(m4, "", "  ")
				if err != nil {
					panic(err)
				}

				return string(b)
			}
		}
	}

	return str
}
