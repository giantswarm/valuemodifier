package path

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidFormatError = errgo.New("invalid format")

// IsInvalidFormat asserts invalidFormatError.
func IsInvalidFormat(err error) bool {
	return errgo.Cause(err) == invalidFormatError
}

var keyNotIndexError = errgo.New("key not index")

// IsKeyNotIndex asserts keyNotIndexError.
func IsKeyNotIndex(err error) bool {
	return errgo.Cause(err) == keyNotIndexError
}

var notFoundError = errgo.New("not found")

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return errgo.Cause(err) == notFoundError
}
