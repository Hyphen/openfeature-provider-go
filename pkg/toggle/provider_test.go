package toggle

import (
	"context"
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

// MockClient implements ClientInterface for testing
type MockClient struct {
	EvaluateFunc      func(ctx EvaluationContext) (*Response, error)
	SendTelemetryFunc func(payload TelemetryPayload) error
}

func (m *MockClient) Evaluate(ctx EvaluationContext) (*Response, error) {
	if m.EvaluateFunc != nil {
		return m.EvaluateFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) SendTelemetry(payload TelemetryPayload) error {
	if m.SendTelemetryFunc != nil {
		return m.SendTelemetryFunc(payload)
	}
	return nil
}

func TestNewProvider(t *testing.T) {

	tests := []struct {
		name       string
		config     Config
		wantErr    bool
		errMessage string
		wantConfig *Config
	}{
		{
			name: "valid config",
			config: Config{
				PublicKey:   "test-key",
				Application: "test-app",
				Environment: "test-env",
			},
			wantErr: false,
			wantConfig: &Config{
				PublicKey:   "test-key",
				Application: "test-app",
				Environment: "test-env",
				HorizonUrls: []string{"https://horizon.hyphen.ai"}, // Include default URL
				Cache:       nil,
			},
		},
		{
			name:       "empty config",
			config:     Config{},
			wantErr:    true,
			errMessage: "application is required", // Updated error message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProvider(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
				assert.Equal(t, *tt.wantConfig, p.config) // Compare with expected config
				assert.NotNil(t, p.client)
				assert.NotEmpty(t, p.endpoints)
				assert.NotEmpty(t, p.hooks)
			}
		})
	}
}

func TestProvider_BooleanEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		flag         string
		defaultValue bool
		evalCtx      openfeature.FlattenedContext
		mockResponse *Response
		mockErr      error
		expected     openfeature.BoolResolutionDetail
	}{
		{
			name:         "successful boolean evaluation",
			flag:         "test-flag",
			defaultValue: false,
			evalCtx: openfeature.FlattenedContext{
				"targetingKey": "user-123",
			},
			mockResponse: &Response{
				Toggles: map[string]Evaluation{
					"test-flag": {
						Type:  "boolean",
						Value: true,
					},
				},
			},
			expected: openfeature.BoolResolutionDetail{
				Value: true,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			},
		},
		{
			name:         "missing targeting key",
			flag:         "test-flag",
			defaultValue: false,
			evalCtx:      openfeature.FlattenedContext{},
			expected: openfeature.BoolResolutionDetail{
				Value: false,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason:          openfeature.ErrorReason,
					ResolutionError: openfeature.NewParseErrorResolutionError(ErrMissingTargetKey.Error()),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClient{
				EvaluateFunc: func(ctx EvaluationContext) (*Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockResponse, nil
				},
			}

			p := &Provider{
				client: mockClient,
				config: Config{
					Application: "test-app",
					Environment: "test-env",
				},
			}

			result := p.BooleanEvaluation(context.Background(), tt.flag, tt.defaultValue, tt.evalCtx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProvider_StringEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		flag         string
		defaultValue string
		evalCtx      openfeature.FlattenedContext
		mockResponse *Response
		mockErr      error
		expected     openfeature.StringResolutionDetail
	}{
		{
			name:         "successful string evaluation",
			flag:         "test-flag",
			defaultValue: "default",
			evalCtx: openfeature.FlattenedContext{
				"targetingKey": "user-123",
			},
			mockResponse: &Response{
				Toggles: map[string]Evaluation{
					"test-flag": {
						Type:  "string",
						Value: "test-value",
					},
				},
			},
			expected: openfeature.StringResolutionDetail{
				Value: "test-value",
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockClient{
				EvaluateFunc: func(ctx EvaluationContext) (*Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockResponse, nil
				},
			}

			p := &Provider{
				client: mockClient,
				config: Config{
					Application: "test-app",
					Environment: "test-env",
				},
			}

			result := p.StringEvaluation(context.Background(), tt.flag, tt.defaultValue, tt.evalCtx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Similar patterns for IntEvaluation, FloatEvaluation, and ObjectEvaluation tests...

func TestProvider_Hooks(t *testing.T) {
	p := &Provider{
		hooks: []openfeature.Hook{
			NewProviderHook(&Provider{}),
		},
	}

	hooks := p.Hooks()
	assert.Len(t, hooks, 1)
}

func TestProvider_buildContext(t *testing.T) {
	tests := []struct {
		name    string
		evalCtx openfeature.FlattenedContext
		want    EvaluationContext
		wantErr bool
	}{
		{
			name: "valid context",
			evalCtx: openfeature.FlattenedContext{
				"targetingKey": "user-123",
				"custom":       "value",
			},
			want: EvaluationContext{
				TargetingKey: "user-123",
				Application:  "test-app",
				Environment:  "test-env",
				CustomAttributes: map[string]interface{}{
					"custom": "value",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing targeting key",
			evalCtx: openfeature.FlattenedContext{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				config: Config{
					Application: "test-app",
					Environment: "test-env",
				},
			}

			got, err := p.buildContext(tt.evalCtx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
