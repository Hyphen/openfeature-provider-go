package toggle

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/open-feature/go-sdk/openfeature"
)

var (
	defaultHorizonURL = "toggle.hyphen.cloud"
	orgIDRegex        = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type Provider struct {
	config    Config
	client    ClientInterface
	endpoints []HorizonEndpoints
	hooks     []openfeature.Hook
}

func extractOrgID(publicKey string) (string, error) {
	key := strings.TrimPrefix(publicKey, "public_")

	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid key format: insufficient parts")
	}

	orgID := parts[0]
	if !orgIDRegex.MatchString(orgID) {
		return "", fmt.Errorf("invalid orgID format")
	}

	return orgID, nil
}

func NewProvider(config Config) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	horizonUrls := config.HorizonUrls
	if len(horizonUrls) == 0 {
		url := fmt.Sprintf("https://%s", defaultHorizonURL)
		if orgID, err := extractOrgID(config.PublicKey); err == nil && orgID != "" {
			url = fmt.Sprintf("https://%s.%s", orgID, defaultHorizonURL)
		}
		horizonUrls = []string{url}
	}

	p := &Provider{
		config:    config,
		endpoints: newEndpoints(horizonUrls),
	}

	client, err := newClient(config, p.endpoints)
	if err != nil {
		return nil, err
	}
	p.client = client

	hook := NewProviderHook(p)
	p.hooks = []openfeature.Hook{hook}

	return p, nil
}

func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "hyphen-provider",
	}
}
func (p *Provider) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:          openfeature.ErrorReason,
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:          openfeature.ErrorReason,
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "boolean" {
		if value, ok := toggle.Value.(bool); ok {
			return openfeature.BoolResolutionDetail{
				Value: value,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		}
	}

	return openfeature.BoolResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:          openfeature.ErrorReason,
			ResolutionError: openfeature.NewTypeMismatchResolutionError("invalid flag type"),
		},
	}
}

func (p *Provider) StringEvaluation(
	ctx context.Context,
	flag string,
	defaultValue string,
	evalCtx openfeature.FlattenedContext,
) openfeature.StringResolutionDetail {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "string" {
		if value, ok := toggle.Value.(string); ok {
			return openfeature.StringResolutionDetail{
				Value: value,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		}
	}

	return openfeature.StringResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError("invalid flag type"),
			Reason:          openfeature.ErrorReason,
		},
	}
}
func (p *Provider) FloatEvaluation(ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "number" {
		switch v := toggle.Value.(type) {
		case float64:
			return openfeature.FloatResolutionDetail{
				Value: v,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		case int:
			return openfeature.FloatResolutionDetail{
				Value: float64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		case int64:
			return openfeature.FloatResolutionDetail{
				Value: float64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		}
	}

	return openfeature.FloatResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError(ErrInvalidFlagType.Error()),
			Reason:          openfeature.ErrorReason,
		},
	}
}

func (p *Provider) IntEvaluation(ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "number" {
		switch v := toggle.Value.(type) {
		case int:
			return openfeature.IntResolutionDetail{
				Value: int64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		case int64:
			return openfeature.IntResolutionDetail{
				Value: v,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		case float64:
			return openfeature.IntResolutionDetail{
				Value: int64(v),
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			}
		}
	}

	return openfeature.IntResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewTypeMismatchResolutionError(ErrInvalidFlagType.Error()),
			Reason:          openfeature.ErrorReason,
		},
	}
}

func (p *Provider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:          openfeature.ErrorReason,
				ResolutionError: openfeature.NewParseErrorResolutionError(err.Error()),
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason:          openfeature.ErrorReason,
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "object" {
		return openfeature.InterfaceResolutionDetail{
			Value: toggle.Value,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				Reason: openfeature.TargetingMatchReason,
			},
		}
	}

	return openfeature.InterfaceResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason:          openfeature.ErrorReason,
			ResolutionError: openfeature.NewTypeMismatchResolutionError(ErrInvalidFlagType.Error()),
		},
	}
}

func (p *Provider) buildContext(evalCtx openfeature.FlattenedContext) (EvaluationContext, error) {
	targetingKey, ok := evalCtx["targetingKey"].(string)
	if !ok {
		return EvaluationContext{}, ErrMissingTargetKey
	}

	// Extract user data from the context
	userData, ok := evalCtx["user"].(map[string]interface{})
	if !ok {
		userData = make(map[string]interface{})
	}

	// Create User struct with data from the nested structure
	user := &User{
		ID:    getString(userData, "id", targetingKey),
		Email: getString(userData, "email", ""),
		Name:  getString(userData, "name", ""),
	}

	// Handle user custom attributes
	if customAttrs, ok := userData["customAttributes"].(map[string]interface{}); ok {
		user.CustomAttributes = customAttrs
	} else {
		user.CustomAttributes = make(map[string]interface{})
	}

	// Build the evaluation context
	ctx := EvaluationContext{
		TargetingKey:     targetingKey,
		Application:      p.config.Application,
		Environment:      p.config.Environment,
		User:             user,
		CustomAttributes: make(map[string]interface{}),
	}

	// Add remaining top-level attributes to customAttributes
	for k, v := range evalCtx {
		switch k {
		case "targetingKey", "user":
			continue
		default:
			ctx.CustomAttributes[k] = v
		}
	}

	return ctx, nil
}

// Helper function to safely get string values from map
func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

// Helper function to determine if an attribute belongs to user
func isUserAttribute(key string) bool {
	userAttributes := map[string]bool{
		"role":    true,
		"group":   true,
		"company": true,
		// Add other user-specific attributes as needed
	}
	return userAttributes[key]
}

func (p *Provider) Hooks() []openfeature.Hook {
	return p.hooks
}
