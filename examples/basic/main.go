package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hyphen-ai/openfeature-provider-go/pkg/toggle"
	"github.com/open-feature/go-sdk/openfeature"
)

type GammaStruct struct {
	Field1 string `json:"field1,omitempty"`
	Field2 int    `json:"field2,omitempty"`
}

func main() {
	config := toggle.Config{
		Application:       "application-id",
		Environment:       "production",
		PublicKey:         "your-apikey",
		HorizonServerURLs: []string{"https://horizon.hyphen.ai"},
	}

	provider, err := toggle.NewProvider(config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	openfeature.SetProvider(provider)
	time.Sleep(100 * time.Millisecond)

	client := openfeature.NewClient("basic-example")

	evalCtx := openfeature.NewEvaluationContext(
		"user-123",
		map[string]interface{}{
			"ipAddress": "203.0.113.42",
			"user": map[string]interface{}{
				"id":    "user-123",
				"email": "user@example.com",
				"name":  "John Doe",
				"customAttributes": map[string]interface{}{
					"role": "admin",
				},
			},
			"customAttributes": map[string]interface{}{
				"subscriptionLevel": "premium",
				"region":            "us-east",
			},
		},
	)

	ctx := context.Background()

	evaluateStringFlags(ctx, client, evalCtx)
	evaluateNumberFlags(ctx, client, evalCtx)
	evaluateBooleanFlag(ctx, client, evalCtx)
	evaluateObjectFlag(ctx, client, evalCtx)
}

func evaluateStringFlags(ctx context.Context, client *openfeature.Client, evalCtx openfeature.EvaluationContext) {
	alphaDetails, err := client.StringValueDetails(ctx, "alpha", "default string", evalCtx)
	if err != nil {
		log.Printf("Error evaluating alpha flag: %v", err)
	}
	fmt.Printf("Feature flag 'alpha' is %s (variant: %s, reason: %s)\n",
		alphaDetails.Value, alphaDetails.Variant, alphaDetails.Reason)

	betaDetails, err := client.StringValueDetails(ctx, "beta", "default string", evalCtx)
	if err != nil {
		log.Printf("Error evaluating beta flag: %v", err)
	}
	fmt.Printf("Feature flag 'beta' is %s (variant: %s, reason: %s)\n",
		betaDetails.Value, betaDetails.Variant, betaDetails.Reason)
}

func evaluateNumberFlags(ctx context.Context, client *openfeature.Client, evalCtx openfeature.EvaluationContext) {
	deltaDetails, err := client.FloatValueDetails(ctx, "delta", 1.0, evalCtx)
	if err != nil {
		log.Printf("Error evaluating delta flag: %v", err)
	}
	fmt.Printf("Feature flag 'delta' is %v (variant: %s, reason: %s)\n",
		deltaDetails.Value, deltaDetails.Variant, deltaDetails.Reason)

	periDetails, err := client.FloatValueDetails(ctx, "peri", 0.0, evalCtx)
	if err != nil {
		log.Printf("Error evaluating peri flag: %v", err)
	}
	fmt.Printf("Feature flag 'peri' is %v (variant: %s, reason: %s)\n",
		periDetails.Value, periDetails.Variant, periDetails.Reason)
}

func evaluateBooleanFlag(ctx context.Context, client *openfeature.Client, evalCtx openfeature.EvaluationContext) {
	tetaDetails, err := client.BooleanValueDetails(ctx, "teta", false, evalCtx)
	if err != nil {
		log.Printf("Error evaluating teta flag: %v", err)
	}
	fmt.Printf("Feature flag 'teta' is %v (variant: %s, reason: %s)\n",
		tetaDetails.Value, tetaDetails.Variant, tetaDetails.Reason)
}

func evaluateObjectFlag(ctx context.Context, client *openfeature.Client, evalCtx openfeature.EvaluationContext) {
	defaultGamma := GammaStruct{}
	gammaDetails, err := client.ObjectValueDetails(ctx, "gamma", defaultGamma, evalCtx)
	if err != nil {
		log.Printf("Error evaluating gamma flag: %v", err)
	}
	fmt.Printf("Feature flag 'gamma' is %+v (variant: %s, reason: %s)\n",
		gammaDetails.Value, gammaDetails.Variant, gammaDetails.Reason)
}
