package toggle

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
