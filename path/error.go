package path

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var keyNotIndexError = errgo.New("key not index")

// IsKeyNotIndex asserts keyNotIndexError.
func IsKeyNotIndex(err error) bool {
	return errgo.Cause(err) == keyNotIndexError
}

var pathNotFoundError = errgo.New("path not index")

// IsPathNotFound asserts pathNotFoundError.
func IsPathNotFound(err error) bool {
	return errgo.Cause(err) == pathNotFoundError
}
