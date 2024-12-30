package toggle

const (
	DefaultHorizonURL = "https://dev-horizon.hyphen.ai"
	DefaultCacheTTL   = 30
)

type endpoints struct {
	Evaluate  string
	Telemetry string
}

func newEndpoints(urls []string) []endpoints {
	result := make([]endpoints, len(urls))
	for i, url := range urls {
		result[i] = endpoints{
			Evaluate:  url + "/toggle/evaluate",
			Telemetry: url + "/toggle/telemetry",
		}
	}
	return result
}

func validateConfig(config Config) error {
	if config.Application == "" {
		return ErrMissingApplication
	}
	if config.Environment == "" {
		return ErrMissingEnvironment
	}
	if config.PublicKey == "" {
		return ErrMissingPublicKey
	}
	return nil
}
