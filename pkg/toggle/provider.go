package toggle

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

type Provider struct {
	config    Config
	client    *Client
	endpoints []HorizonEndpoints
	hooks     []openfeature.Hook
}

func NewProvider(config Config) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	if len(config.HorizonServerURLs) == 0 {
		config.HorizonServerURLs = []string{horizon.URL}
	}

	p := &Provider{
		config:    config,
		endpoints: newEndpoints(config.HorizonServerURLs),
	}

	client, err := newClient(config)
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

	ctx := EvaluationContext{
		TargetingKey:     targetingKey,
		Application:      p.config.Application,
		Environment:      p.config.Environment,
		CustomAttributes: make(map[string]interface{}),
	}

	// Copy additional attributes
	for k, v := range evalCtx {
		if k != "targetingKey" {
			ctx.CustomAttributes[k] = v
		}
	}

	return ctx, nil
}

func (p *Provider) Hooks() []openfeature.Hook {
	return p.hooks
}
