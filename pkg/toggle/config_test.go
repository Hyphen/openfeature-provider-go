package toggle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		expected []HorizonEndpoints
	}{
		{
			name: "single url",
			urls: []string{"https://test.com"},
			expected: []HorizonEndpoints{
				{
					Evaluate:  "https://test.com/toggle/evaluate",
					Telemetry: "https://test.com/toggle/telemetry",
				},
			},
		},
		{
			name: "multiple urls",
			urls: []string{"https://test1.com", "https://test2.com"},
			expected: []HorizonEndpoints{
				{
					Evaluate:  "https://test1.com/toggle/evaluate",
					Telemetry: "https://test1.com/toggle/telemetry",
				},
				{
					Evaluate:  "https://test2.com/toggle/evaluate",
					Telemetry: "https://test2.com/toggle/telemetry",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newEndpoints(tt.urls)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr error
	}{
		{
			name: "valid config",
			config: Config{
				Application: "test-app",
				Environment: "test-env",
				PublicKey:   "test-key",
			},
			wantErr: nil,
		},
		{
			name: "missing application",
			config: Config{
				Environment: "test-env",
				PublicKey:   "test-key",
			},
			wantErr: ErrMissingApplication,
		},
		{
			name: "missing environment",
			config: Config{
				Application: "test-app",
				PublicKey:   "test-key",
			},
			wantErr: ErrMissingEnvironment,
		},
		{
			name: "missing public key",
			config: Config{
				Application: "test-app",
				Environment: "test-env",
			},
			wantErr: ErrMissingPublicKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
