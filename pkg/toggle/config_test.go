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

func TestValidateEnvironmentFormat(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantErr     error
	}{
		{
			name:        "valid project environment ID",
			environment: "pevr_abc123",
			wantErr:     nil,
		},
		{
			name:        "valid alternateId - simple",
			environment: "production",
			wantErr:     nil,
		},
		{
			name:        "valid alternateId - with hyphens and underscores",
			environment: "prod-us_east",
			wantErr:     nil,
		},
		{
			name:        "valid alternateId - with numbers",
			environment: "prod123",
			wantErr:     nil,
		},
		{
			name:        "invalid - contains 'environments'",
			environment: "test-environments-prod",
			wantErr:     ErrInvalidEnvironmentFormat,
		},
		{
			name:        "invalid - uppercase letters",
			environment: "Production",
			wantErr:     ErrInvalidEnvironmentFormat,
		},
		{
			name:        "invalid - too long (26 characters)",
			environment: "abcdefghijklmnopqrstuvwxyz",
			wantErr:     ErrInvalidEnvironmentFormat,
		},
		{
			name:        "invalid - special characters",
			environment: "prod@test",
			wantErr:     ErrInvalidEnvironmentFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnvironmentFormat(tt.environment)
			assert.Equal(t, tt.wantErr, err)
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
			name: "valid config with project environment ID",
			config: Config{
				Application: "test-app",
				Environment: "pevr_abc123",
				PublicKey:   "test-key",
			},
			wantErr: nil,
		},
		{
			name: "valid config with alternateId",
			config: Config{
				Application: "test-app",
				Environment: "production",
				PublicKey:   "test-key",
			},
			wantErr: nil,
		},
		{
			name: "missing application",
			config: Config{
				Environment: "production",
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
			name: "invalid environment format",
			config: Config{
				Application: "test-app",
				Environment: "Production", // Invalid: uppercase
				PublicKey:   "test-key",
			},
			wantErr: ErrInvalidEnvironmentFormat,
		},
		{
			name: "missing public key",
			config: Config{
				Application: "test-app",
				Environment: "production",
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
