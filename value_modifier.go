package valuemodifier

import (
	"sort"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/giantswarm/valuemodifier/path"
	"github.com/spf13/cast"
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
	var pathService *path.Service
	{
		pathConfig := path.DefaultConfig()
		pathConfig.InputBytes = input
		pathService, err = path.New(pathConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	{
		var fields []string
		fields = append(fields, s.ignoreFields...)
		fields = append(fields, s.selectFields...)

		err := pathService.Validate(fields)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var paths []string
	{
		paths, err = pathService.All()
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
		v, err := pathService.Get(p)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		b := []byte(cast.ToString(v))
		for _, m := range s.valueModifiers {
			b, err = m.Modify(b)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		err = pathService.Set(p, string(b))
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	return pathService.OutputBytes(), nil
}

func containsString(list []string, item string) bool {
	for _, l := range list {
		if l == item {
			return true
		}
	}

	return false
}
