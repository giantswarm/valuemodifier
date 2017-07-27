package path

import "github.com/giantswarm/microerror"

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFormatError = microerror.New("invalid format")

// IsInvalidFormat asserts invalidFormatError.
func IsInvalidFormat(err error) bool {
	return microerror.Cause(err) == invalidFormatError
}

var keyNotIndexError = microerror.New("key not index")

// IsKeyNotIndex asserts keyNotIndexError.
func IsKeyNotIndex(err error) bool {
	return microerror.Cause(err) == keyNotIndexError
}

var notFoundError = microerror.New("not found")

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}
