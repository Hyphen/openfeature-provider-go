package toggle

import "time"

type Config struct {
	PublicKey        string
	Application      string
	Environment      string
	HorizonServerURL string
	EnableUsage      bool
	Cache            *CacheConfig
}

type CacheConfig struct {
	TTL    time.Duration
	KeyGen func(ctx EvaluationContext) string
}

type EvaluationContext struct {
	TargetingKey string                 `json:"targetingKey"`
	IPAddress    string                 `json:"ipAddress,omitempty"`
	Application  string                 `json:"application"`
	Environment  string                 `json:"environment"`
	User         *User                  `json:"user,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

type User struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type Evaluation struct {
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
	Type   string      `json:"type"`
	Reason string      `json:"reason,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type Response struct {
	Toggles map[string]Evaluation `json:"toggles"`
}

type TelemetryPayload struct {
	Context EvaluationContext `json:"context"`
	Data    struct {
		Toggle Evaluation `json:"toggle"`
	} `json:"data"`
}
