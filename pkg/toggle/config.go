package toggle

import (
	"regexp"
	"strings"
)

const (
	DefaultHorizonURL = "https://dev-horizon.hyphen.ai"
	DefaultCacheTTL   = 30
)

type HorizonConfig struct {
	URL string
}

type HorizonEndpoints struct {
	Evaluate  string
	Telemetry string
}

var (
	horizon = HorizonConfig{
		URL: "https://horizon.hyphen.ai",
	}
	// Regex for validating alternateId format
	alternateIdRegex = regexp.MustCompile(`^[a-z0-9\-_]{1,25}$`)
)

func newEndpoints(urls []string) []HorizonEndpoints {
	endpoints := make([]HorizonEndpoints, len(urls))
	for i, url := range urls {
		endpoints[i] = HorizonEndpoints{
			Evaluate:  url + "/toggle/evaluate",
			Telemetry: url + "/toggle/telemetry",
		}
	}
	return endpoints
}

// validateEnvironmentFormat validates that the environment identifier follows one of these formats:
// - A project environment ID that starts with the prefix "pevr_" followed by alphanumeric characters
// - A valid alternateId that meets these criteria:
//   - Contains only lowercase letters, numbers, hyphens, and underscores
//   - Has a length between 1 and 25 characters
//   - Does not contain the word "environments"
func validateEnvironmentFormat(environment string) error {
	// Check if it's a project environment ID (starts with 'pevr_')
	isEnvironmentId := strings.HasPrefix(environment, "pevr_")
	
	// Check if it's a valid alternateId
	isValidAlternateId := alternateIdRegex.MatchString(environment) && 
		!strings.Contains(environment, "environments")
	
	if !isEnvironmentId && !isValidAlternateId {
		return ErrInvalidEnvironmentFormat
	}
	
	return nil
}

func validateConfig(config Config) error {
	if config.Application == "" {
		return ErrMissingApplication
	}
	if config.Environment == "" {
		return ErrMissingEnvironment
	}
	
	// Validate environment format
	if err := validateEnvironmentFormat(config.Environment); err != nil {
		return err
	}
	
	if config.PublicKey == "" {
		return ErrMissingPublicKey
	}
	return nil
}
