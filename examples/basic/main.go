package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hyphen-ai/openfeature-provider-go/pkg/toggle"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	// Configuration for the Hyphen provider
	config := toggle.Config{
		Application: "application-id",
		Environment: "production",
		PublicKey:   "your-public-key",
	}

	// Initialize the provider
	provider, err := toggle.NewProvider(config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Register the provider
	openfeature.SetProvider(provider)

	// Create an OpenFeature client
	client := openfeature.NewClient("basic-example")

	// Define evaluation context
	evalCtx := openfeature.NewEvaluationContext(
		"user-123",
		map[string]interface{}{
			"targetingKey": "user-123",
			"ipAddress":    "203.0.113.42",
			"customAttributes": map[string]interface{}{
				"subscriptionLevel": "premium",
				"region":            "us-east",
			},
			"user": map[string]interface{}{
				"id":    "user-123",
				"email": "user@example.com",
				"name":  "John Doe",
				"customAttributes": map[string]interface{}{
					"role": "admin",
				},
			},
		},
	)

	// Evaluate a feature flag
	ctx := context.Background()
	flagKey := "my-bool-toggle"
	defaultValue := false
	result, err := client.BooleanValue(ctx, flagKey, defaultValue, evalCtx)
	if err != nil {
		// Handle the error appropriately
		log.Fatalf("Error evaluating flag: %v", err)
	}
	fmt.Printf("Feature flag '%s' is %t\n", flagKey, result)
}
