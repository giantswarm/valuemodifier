package valuemodifier

import (
	"encoding/json"
	"sort"

	"github.com/Jeffail/gabs"
	yamltojson "github.com/ghodss/yaml"
	microerror "github.com/giantswarm/microkit/error"
	yaml "gopkg.in/yaml.v2"
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

// DefaultConfig provides a default configuration to create a new value modifier
// traverser by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		ValueModifiers: nil,

		// Settings.
		IgnoreFields: nil,
		SelectFields: nil,
	}
}

// New creates a new configured value modifier traverser.
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

// TODO indentation as provided
// TODO without indentation as provided
func (s *Service) Traverse(input []byte) ([]byte, error) {
	jsonBytes, _, _, err := toJSON(input)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	container, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	{
		var fields []string
		fields = append(fields, s.ignoreFields...)
		fields = append(fields, s.selectFields...)
		err := validateFields(jsonBytes, fields)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var paths []string
	{
		paths, err = containerPaths(container)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		if len(s.ignoreFields) != 0 {
			var newPaths []string

			for _, p := range paths {
				if containsString(s.ignoreFields, p) {
					continue
				}
				newPaths = append(newPaths, p)
			}

			paths = newPaths
		} else if len(s.selectFields) != 0 {
			paths = s.selectFields
		}

		sort.Strings(paths)
	}

	for _, p := range paths {

		// TODO get value by path

		v := []byte(str)
		for _, m := range s.valueModifiers {
			v, err = m.Modify(v)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}
		str = string(v)

		// TODO set value by path
	}

	return []byte(container.StringIndent("", "  ")), nil
}

func containsString(list []string, item string) bool {
	for _, l := range list {
		if l == item {
			return true
		}
	}

	return false
}

func isJSON(b []byte) bool {
	var m map[string]interface{}
	return json.Unmarshal(b, &m) == nil
}

func isYAML(b []byte) bool {
	var m map[interface{}]interface{}
	return yaml.Unmarshal(b, &m) == nil
}

func toJSON(b []byte) ([]byte, map[string]interface{}, bool, error) {
	if isJSON(b) {
		var m map[string]interface{}
		err := json.Unmarshal(b, &m)
		if err != nil {
			return nil, nil, false, microerror.MaskAny(err)
		}

		return b, m, true, nil
	}

	var jsonMap map[string]interface{}
	var jsonBytes []byte
	{
		var m map[interface{}]interface{}
		err := yaml.Unmarshal(b, &m)
		if err != nil {
			return nil, nil, false, microerror.MaskAny(err)
		}

		yamlBytes, err := yaml.Marshal(m)
		if err != nil {
			return nil, nil, false, microerror.MaskAny(err)
		}

		jsonBytes, err = yamltojson.YAMLToJSON(yamlBytes)
		if err != nil {
			return nil, nil, false, microerror.MaskAny(err)
		}
		err = json.Unmarshal(b, &jsonMap)
		if err != nil {
			return nil, nil, false, microerror.MaskAny(err)
		}
	}

	return jsonBytes, jsonMap, false, nil
}

func validateFields(jsonBytes []byte, fields []string) error {
	gabsJSON, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		return microerror.MaskAny(err)
	}

	for _, f := range fields {
		exists := gabsJSON.ExistsP(f)
		if !exists {
			return microerror.MaskAnyf(fieldNotFoundError, f)
		}
	}

	return nil
}
