package toggle

import "errors"

var (
	ErrMissingApplication       = errors.New("application is required")
	ErrMissingEnvironment       = errors.New("environment is required")
	ErrMissingPublicKey         = errors.New("public key is required")
	ErrMissingTargetKey         = errors.New("targeting key is required")
	ErrInvalidFlagType          = errors.New("invalid flag type")
	ErrFlagNotFound             = errors.New("flag not found")
	ErrInvalidEnvironmentFormat = errors.New("invalid environment format. Must be either a project environment ID (starting with \"pevr_\") or a valid alternateId (1-25 characters, lowercase letters, numbers, hyphens, and underscores, not containing the word \"environments\")")
)
