package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hyphen/openfeature-provider-go/pkg/toggle"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	// Configuration for the Hyphen provider
	config := toggle.Config{
		Application: "app",
		// Using alternateId format for environment
		Environment: "production",
		// Alternatively, you can use a project environment ID:
		// Environment: "pevr_abc123",
		PublicKey:   "PUBLIC_KEY_HERE",
	}

	// Initialize the provider
	provider, err := toggle.NewProvider(config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Register the provider
	openfeature.SetProvider(provider)

	// Add a small delay to ensure provider initialization
	time.Sleep(100 * time.Millisecond)

	// Create an OpenFeature client
	client := openfeature.NewClient("basic-example")

	// Define evaluation context
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

	// Add context with logger
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", log.Default())

	type GammaStruct struct {
		Enabled     bool    `json:"enabled"`
		Probability float64 `json:"probability"`
		Parameters  struct {
			Alpha float64 `json:"alpha"`
			Beta  float64 `json:"beta"`
		} `json:"parameters"`
		Settings struct {
			MaxIterations int     `json:"maxIterations"`
			Tolerance     float64 `json:"tolerance"`
			Seed          int64   `json:"seed"`
		} `json:"settings"`
	}
	// Evaluate a feature flag
	flagKey := "gamma"
	defaultGamma := GammaStruct{}

	result, err := client.ObjectValueDetails(ctx, flagKey, defaultGamma, evalCtx)
	if err != nil {
		log.Printf("Evaluation context: %+v", evalCtx)
		log.Fatalf("Error evaluating flag: %v", err)
	}

	fmt.Printf("Feature flag '%s' is %s\n", flagKey, result)
}
