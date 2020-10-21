package encrypt

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var vaultResponseError = &microerror.Error{
	Kind: "vaultResponseError",
}

// IsVaultResponseError asserts vaultResponseError.
func IsVaultResponseError(err error) bool {
	return microerror.Cause(err) == vaultResponseError
}
