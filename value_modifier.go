package valuemodifier

import (
	"encoding/json"

	microerror "github.com/giantswarm/microkit/error"
	yaml "gopkg.in/yaml.v1"
)

func traverseJSON(input []byte, valueModifiers ...ValueModifier) ([]byte, error) {
	var m map[string]interface{}
	err := json.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	for k, v := range m {
		m[k] = toModifiedValueJSON(v, valueModifiers...)
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func traverseYAML(input []byte, valueModifiers ...ValueModifier) ([]byte, error) {
	var m map[interface{}]interface{}
	err := yaml.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	for k, v := range m {
		m[k] = toModifiedValueYAML(v, valueModifiers...)
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func toModifiedValueJSON(val interface{}, valueModifiers ...ValueModifier) interface{} {
	m, ok := val.(map[string]interface{})
	if ok {
		for k, v := range m {
			m[k] = toModifiedValueJSON(v, valueModifiers...)
		}

		return m
	}

	s := cast.ToString(val)
	if s != "" {
		for _, m := range valueModifiers {
			o, err := m.Modify([]byte(s))
			if err != nil {
				panic(err)
			}
			s = string(o)
		}
	}

	return s
}

func toModifiedValueYAML(val interface{}, valueModifiers ...ValueModifier) interface{} {
	m, ok := val.(map[interface{}]interface{})
	if ok {
		for k, v := range m {
			m[k] = toModifiedValueYAML(v, valueModifiers...)
		}

		return m
	}

	s := cast.ToString(val)
	if s != "" {
		for _, m := range valueModifiers {
			o, err := m.Modify([]byte(s))
			if err != nil {
				panic(err)
			}
			s = string(o)
		}
	}

	return s
}
