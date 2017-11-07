package decrypt

import "github.com/giantswarm/microerror"

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var wrongGPGPasswordError = microerror.New("wrong GPG password")

// IsWrongGPGPassword asserts wrongGPGPassword.
func IsWrongGPGPassword(err error) bool {
	return microerror.Cause(err) == wrongGPGPasswordError
}
