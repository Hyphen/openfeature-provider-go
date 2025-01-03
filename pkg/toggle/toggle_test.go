package toggle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluationContext(t *testing.T) {
	ctx := EvaluationContext{
		TargetingKey: "test-key",
		IPAddress:    "127.0.0.1",
		Application:  "test-app",
		Environment:  "test-env",
		User: &User{
			ID:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
			CustomAttributes: map[string]interface{}{
				"role": "admin",
			},
		},
		CustomAttributes: map[string]interface{}{
			"version": "1.0.0",
		},
	}

	assert.Equal(t, "test-key", ctx.TargetingKey)
	assert.Equal(t, "127.0.0.1", ctx.IPAddress)
	assert.Equal(t, "test-app", ctx.Application)
	assert.Equal(t, "test-env", ctx.Environment)
	assert.NotNil(t, ctx.User)
	assert.Equal(t, "user-123", ctx.User.ID)
	assert.Equal(t, "test@example.com", ctx.User.Email)
	assert.Equal(t, "Test User", ctx.User.Name)
	assert.Equal(t, "admin", ctx.User.CustomAttributes["role"])
	assert.Equal(t, "1.0.0", ctx.CustomAttributes["version"])
}

func TestResponse(t *testing.T) {
	resp := Response{
		Toggles: map[string]Evaluation{
			"test-flag": {
				Key:    "test-flag",
				Value:  true,
				Type:   "boolean",
				Reason: "targeting_match",
			},
		},
	}

	assert.Contains(t, resp.Toggles, "test-flag")
	assert.Equal(t, true, resp.Toggles["test-flag"].Value)
	assert.Equal(t, "boolean", resp.Toggles["test-flag"].Type)
	assert.Equal(t, "targeting_match", resp.Toggles["test-flag"].Reason)
}

func TestTelemetryPayload(t *testing.T) {
	payload := TelemetryPayload{
		Context: EvaluationContext{
			TargetingKey: "test-key",
			Application:  "test-app",
			Environment:  "test-env",
		},
		Data: struct {
			Toggle Evaluation `json:"toggle"`
		}{
			Toggle: Evaluation{
				Key:    "test-flag",
				Value:  true,
				Type:   "boolean",
				Reason: "targeting_match",
			},
		},
	}

	assert.Equal(t, "test-key", payload.Context.TargetingKey)
	assert.Equal(t, "test-flag", payload.Data.Toggle.Key)
	assert.Equal(t, true, payload.Data.Toggle.Value)
}
