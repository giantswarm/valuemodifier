package decrypt

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var wrongGPGPasswordError = &microerror.Error{
	Kind: "wrongGPGPasswordError",
}

// IsWrongGPGPassword asserts wrongGPGPassword.
func IsWrongGPGPassword(err error) bool {
	return microerror.Cause(err) == wrongGPGPasswordError
}
