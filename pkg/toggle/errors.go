package toggle

import "errors"

var (
	ErrMissingApplication = errors.New("application is required")
	ErrMissingEnvironment = errors.New("environment is required")
	ErrMissingPublicKey   = errors.New("public key is required")
	ErrMissingTargetKey   = errors.New("targeting key is required")
	ErrInvalidFlagType    = errors.New("invalid flag type")
	ErrFlagNotFound       = errors.New("flag not found")
)
