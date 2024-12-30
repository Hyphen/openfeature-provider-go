package hyphen

const (
	DefaultHorizonURL = "https://dev-horizon.hyphen.ai"
	DefaultCacheTTL   = 30
)

type endpoints struct {
	Evaluate  string
	Telemetry string
}

func newEndpoints(baseURL string) endpoints {
	return endpoints{
		Evaluate:  baseURL + "/toggle/evaluate",
		Telemetry: baseURL + "/toggle/telemetry",
	}
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
