package hyphen

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
			Error: err,
		}
	}

	eval, err := p.client.Evaluate(hyphenCtx)
	if err != nil {
		return openfeature.BooleanEvaluationDetails{
			Value: defaultValue,
			Error: err,
		}
	}

	if toggle, ok := eval.Toggles[flag]; ok && toggle.Type == "boolean" {
		if value, ok := toggle.Value.(bool); ok {
			return openfeature.BooleanEvaluationDetails{
				Value:  value,
				Reason: openfeature.TargetingMatchReason,
			}
		}
	}

	return openfeature.BooleanEvaluationDetails{
		Value: defaultValue,
		Error: ErrInvalidFlagType,
	}
}

// Similar implementations for StringEvaluation, NumberEvaluation, and ObjectEvaluation...

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
