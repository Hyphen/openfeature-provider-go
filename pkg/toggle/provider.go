package toggle

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
)

type Provider struct {
	config    Config
	client    *Client
	endpoints endpoints
}

func NewProvider(config Config) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	if config.HorizonServerURL == "" {
		config.HorizonServerURL = DefaultHorizonURL
	}

	p := &Provider{
		config:    config,
		endpoints: newEndpoints(config.HorizonServerURL),
	}

	client, err := newClient(config)
	if err != nil {
		return nil, err
	}
	p.client = client

	return p, nil
}

func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "hyphen-provider",
	}
}

func (p *Provider) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BooleanEvaluationDetails {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.BooleanEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Boolean,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.ParseErrorCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.BooleanEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Boolean,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.GeneralCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "boolean" {
		if value, ok := toggle.Value.(bool); ok {
			return openfeature.BooleanEvaluationDetails{
				Value: value,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  flag,
					FlagType: openfeature.Boolean,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason: openfeature.TargetingMatchReason,
					},
				},
			}
		}
	}

	return openfeature.BooleanEvaluationDetails{
		Value: defaultValue,
		EvaluationDetails: openfeature.EvaluationDetails{
			FlagKey:  flag,
			FlagType: openfeature.Boolean,
			ResolutionDetail: openfeature.ResolutionDetail{
				Reason:       openfeature.ErrorReason,
				ErrorCode:    openfeature.TypeMismatchCode,
				ErrorMessage: ErrInvalidFlagType.Error(),
			},
		},
	}
}

func (p *Provider) StringEvaluation(ctx context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringEvaluationDetails {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.StringEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.String,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.ParseErrorCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.StringEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.String,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.GeneralCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "string" {
		if value, ok := toggle.Value.(string); ok {
			return openfeature.StringEvaluationDetails{
				Value: value,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  flag,
					FlagType: openfeature.String,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason: openfeature.TargetingMatchReason,
					},
				},
			}
		}
	}

	return openfeature.StringEvaluationDetails{
		Value: defaultValue,
		EvaluationDetails: openfeature.EvaluationDetails{
			FlagKey:  flag,
			FlagType: openfeature.String,
			ResolutionDetail: openfeature.ResolutionDetail{
				Reason:       openfeature.ErrorReason,
				ErrorCode:    openfeature.TypeMismatchCode,
				ErrorMessage: ErrInvalidFlagType.Error(),
			},
		},
	}
}

func (p *Provider) FloatEvaluation(ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatEvaluationDetails {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.FloatEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Float,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.ParseErrorCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.FloatEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Float,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.GeneralCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "float" {
		if value, ok := toggle.Value.(float64); ok {
			return openfeature.FloatEvaluationDetails{
				Value: value,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  flag,
					FlagType: openfeature.Float,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason: openfeature.TargetingMatchReason,
					},
				},
			}
		}
	}

	return openfeature.FloatEvaluationDetails{
		Value: defaultValue,
		EvaluationDetails: openfeature.EvaluationDetails{
			FlagKey:  flag,
			FlagType: openfeature.Float,
			ResolutionDetail: openfeature.ResolutionDetail{
				Reason:       openfeature.ErrorReason,
				ErrorCode:    openfeature.TypeMismatchCode,
				ErrorMessage: ErrInvalidFlagType.Error(),
			},
		},
	}
}

func (p *Provider) IntEvaluation(ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntEvaluationDetails {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.IntEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Int,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.ParseErrorCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.IntEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Int,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.GeneralCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "int" {
		if value, ok := toggle.Value.(int64); ok {
			return openfeature.IntEvaluationDetails{
				Value: value,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  flag,
					FlagType: openfeature.Int,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason: openfeature.TargetingMatchReason,
					},
				},
			}
		}
	}

	return openfeature.IntEvaluationDetails{
		Value: defaultValue,
		EvaluationDetails: openfeature.EvaluationDetails{
			FlagKey:  flag,
			FlagType: openfeature.Int,
			ResolutionDetail: openfeature.ResolutionDetail{
				Reason:       openfeature.ErrorReason,
				ErrorCode:    openfeature.TypeMismatchCode,
				ErrorMessage: ErrInvalidFlagType.Error(),
			},
		},
	}
}

func (p *Provider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceEvaluationDetails {
	hyphenCtx, err := p.buildContext(evalCtx)
	if err != nil {
		return openfeature.InterfaceEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Object,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.ParseErrorCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.InterfaceEvaluationDetails{
			Value: defaultValue,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Object,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason:       openfeature.ErrorReason,
					ErrorCode:    openfeature.GeneralCode,
					ErrorMessage: err.Error(),
				},
			},
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "object" {
		return openfeature.InterfaceEvaluationDetails{
			Value: toggle.Value,
			EvaluationDetails: openfeature.EvaluationDetails{
				FlagKey:  flag,
				FlagType: openfeature.Object,
				ResolutionDetail: openfeature.ResolutionDetail{
					Reason: openfeature.TargetingMatchReason,
				},
			},
		}
	}

	return openfeature.InterfaceEvaluationDetails{
		Value: defaultValue,
		EvaluationDetails: openfeature.EvaluationDetails{
			FlagKey:  flag,
			FlagType: openfeature.Object,
			ResolutionDetail: openfeature.ResolutionDetail{
				Reason:       openfeature.ErrorReason,
				ErrorCode:    openfeature.TypeMismatchCode,
				ErrorMessage: ErrInvalidFlagType.Error(),
			},
		},
	}
}

func (p *Provider) buildContext(evalCtx openfeature.FlattenedContext) (EvaluationContext, error) {
	targetingKey, ok := evalCtx["targetingKey"].(string)
	if !ok {
		return EvaluationContext{}, ErrMissingTargetKey
	}

	ctx := EvaluationContext{
		TargetingKey: targetingKey,
		Application:  p.config.Application,
		Environment:  p.config.Environment,
		Attributes:   make(map[string]interface{}),
	}

	// Copy additional attributes
	for k, v := range evalCtx {
		if k != "targetingKey" {
			ctx.Attributes[k] = v
		}
	}

	return ctx, nil
}
