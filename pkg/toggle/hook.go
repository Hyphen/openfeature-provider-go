package toggle

import (
	"context"
	"strconv"

	"github.com/open-feature/go-sdk/openfeature"
	"golang.org/x/exp/rand"
)

type ProviderHook struct {
	openfeature.UnimplementedHook
	provider *Provider
}

func NewProviderHook(provider *Provider) *ProviderHook {
	return &ProviderHook{
		provider: provider,
	}
}

func (h *ProviderHook) Before(ctx context.Context, hookContext openfeature.HookContext, hookHints openfeature.HookHints) (*openfeature.EvaluationContext, error) {
	attributes := make(map[string]interface{})

	for k, v := range hookContext.EvaluationContext().Attributes() {
		attributes[k] = v
	}

	attributes["application"] = h.provider.config.Application
	attributes["environment"] = h.provider.config.Environment

	targetingKey := hookContext.EvaluationContext().TargetingKey()
	if targetingKey == "" {
		if userID, ok := getUserID(attributes); ok {
			targetingKey = userID
		} else {
			targetingKey = generateTargetingKey(h.provider.config.Application, h.provider.config.Environment)
		}
	}

	newCtx := openfeature.NewEvaluationContext(
		targetingKey,
		attributes,
	)

	return &newCtx, nil
}

func (h *ProviderHook) After(ctx context.Context, hookContext openfeature.HookContext, details openfeature.InterfaceEvaluationDetails, hookHints openfeature.HookHints) error {
	if h.provider.config.EnableUsage != nil && !*h.provider.config.EnableUsage {
		return nil
	}

	evalCtx := hookContext.EvaluationContext()

	attributes := evalCtx.Attributes()

	flattenedAttributes := make(map[string]interface{})

	if region, ok := attributes["region"].(string); ok {
		flattenedAttributes["region"] = region
	}
	if subscriptionLevel, ok := attributes["subscriptionLevel"].(string); ok {
		flattenedAttributes["subscriptionLevel"] = subscriptionLevel
	}
	if ipAddress, ok := attributes["ipAddress"].(string); ok {
		flattenedAttributes["ipAddress"] = ipAddress
	}

	if user, ok := attributes["user"].(map[string]interface{}); ok {
		if email, ok := user["email"].(string); ok {
			flattenedAttributes["user.email"] = email
		}
		if id, ok := user["id"].(string); ok {
			flattenedAttributes["user.id"] = id
		}
		if name, ok := user["name"].(string); ok {
			flattenedAttributes["user.name"] = name
		}
		if userCustomAttrs, ok := user["customAttributes"].(map[string]interface{}); ok {
			if role, ok := userCustomAttrs["role"].(string); ok {
				flattenedAttributes["user.role"] = role
			}
		}
	}

	hyphenCtx := EvaluationContext{
		TargetingKey:     evalCtx.TargetingKey(),
		Application:      h.provider.config.Application,
		Environment:      h.provider.config.Environment,
		CustomAttributes: flattenedAttributes,
	}

	payload := TelemetryPayload{
		Context: hyphenCtx,
		Data: struct {
			Toggle Evaluation `json:"toggle"`
		}{
			Toggle: Evaluation{
				Key:    details.FlagKey,
				Value:  details.Value,
				Type:   strconv.FormatInt(int64(details.FlagType), 10),
				Reason: string(details.ResolutionDetail.Reason),
			},
		},
	}

	return h.provider.client.SendTelemetry(payload)
}

func (h *ProviderHook) Error(ctx context.Context, hookContext openfeature.HookContext, err error, hookHints openfeature.HookHints) {
	if logger, ok := ctx.Value("logger").(interface{ Error(args ...interface{}) }); ok {
		logger.Error("Error in hook:", err)
	}
}

func getUserID(attributes map[string]interface{}) (string, bool) {
	if user, ok := attributes["user"].(map[string]interface{}); ok {
		if id, ok := user["id"].(string); ok {
			return id, true
		}
	}
	return "", false
}

func generateTargetingKey(application, environment string) string {
	return application + "-" + environment + "-" + getRandomString(7)
}

func getRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
