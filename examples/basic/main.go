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
		Application: "application_id",
		Environment: "production",
		PublicKey:   "your-key",
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

	// Evaluate a feature flag
	flagKey := "beta"
	defaultValue := "default value"

	result, err := client.StringValue(ctx, flagKey, defaultValue, evalCtx)
	if err != nil {
		log.Printf("Evaluation context: %+v", evalCtx)
		log.Fatalf("Error evaluating flag: %v", err)
	}

	fmt.Printf("Feature flag '%s' is %s\n", flagKey, result)
}
