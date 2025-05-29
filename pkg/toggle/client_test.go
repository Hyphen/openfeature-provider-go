package toggle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		endpoints []HorizonEndpoints
		wantErr   bool
	}{
		{
			name: "valid config",
			config: Config{
				PublicKey:   "test-key",
				HorizonUrls: []string{"http://test.com"},
			},
			endpoints: []HorizonEndpoints{
				{
					Evaluate:  "http://test.com/toggle/evaluate",
					Telemetry: "http://test.com/toggle/telemetry",
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with cache",
			config: Config{
				PublicKey: "test-key",
				Cache: &CacheConfig{
					TTL: time.Minute,
					KeyGen: func(ctx EvaluationContext) string {
						return ctx.TargetingKey
					},
				},
			},
			endpoints: []HorizonEndpoints{
				{
					Evaluate:  "http://test.com/toggle/evaluate",
					Telemetry: "http://test.com/toggle/telemetry",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newClient(tt.config, tt.endpoints)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, client)
			if tt.config.Cache != nil {
				assert.NotNil(t, client.cache)
				assert.NotNil(t, client.keyGen)
			}
		})
	}
}

func TestClientEvaluate(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))

		response := Response{
			Toggles: map[string]Evaluation{
				"test-flag": {
					Key:   "test-flag",
					Value: true,
					Type:  "boolean",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		PublicKey:   "test-key",
		HorizonUrls: []string{server.URL},
	}

	endpoints := []HorizonEndpoints{
		{
			Evaluate:  fmt.Sprintf("%s/toggle/evaluate", server.URL),
			Telemetry: fmt.Sprintf("%s/toggle/telemetry", server.URL),
		},
	}

	client, err := newClient(config, endpoints)
	assert.NoError(t, err)

	ctx := EvaluationContext{
		TargetingKey: "test-user",
		Application:  "test-app",
		Environment:  "test-env",
	}

	resp, err := client.Evaluate(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Toggles, "test-flag")
}

func TestClientEvaluateWithCache(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		response := Response{
			Toggles: map[string]Evaluation{
				"test-flag": {
					Key:   "test-flag",
					Value: true,
					Type:  "boolean",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := Config{
		PublicKey:   "test-key",
		HorizonUrls: []string{server.URL},
		Cache: &CacheConfig{
			TTL: time.Minute,
			KeyGen: func(ctx EvaluationContext) string {
				return ctx.TargetingKey
			},
		},
	}
	endpoints := []HorizonEndpoints{
		{
			Evaluate:  fmt.Sprintf("%s/toggle/evaluate", server.URL),
			Telemetry: fmt.Sprintf("%s/toggle/telemetry", server.URL),
		},
	}

	client, err := newClient(config, endpoints)
	assert.NoError(t, err)

	ctx := EvaluationContext{
		TargetingKey: "test-user",
		Application:  "test-app",
		Environment:  "test-env",
	}

	// First call should hit the server
	resp1, err := client.Evaluate(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)
	assert.Equal(t, 1, callCount)

	// Second call should use cache
	resp2, err := client.Evaluate(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, 1, callCount) // Call count should not increase
}

func TestClientSendTelemetry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		PublicKey:   "test-key",
		HorizonUrls: []string{server.URL},
	}
	endpoints := []HorizonEndpoints{
		{
			Evaluate:  fmt.Sprintf("%s/toggle/evaluate", server.URL),
			Telemetry: fmt.Sprintf("%s/toggle/telemetry", server.URL),
		},
	}

	client, err := newClient(config, endpoints)
	assert.NoError(t, err)

	payload := TelemetryPayload{
		Context: EvaluationContext{
			TargetingKey: "test-user",
			Application:  "test-app",
			Environment:  "test-env",
		},
		Data: struct {
			Toggle Evaluation `json:"toggle"`
		}{
			Toggle: Evaluation{
				Key:   "test-flag",
				Value: true,
				Type:  "boolean",
			},
		},
	}

	err = client.SendTelemetry(payload)
	assert.NoError(t, err)
}
