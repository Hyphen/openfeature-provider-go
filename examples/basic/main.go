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
		Application:       "application-id",
		Environment:       "production",
		PublicKey:         "your-public-key",
		HorizonServerURLs: []string{"https://horizon.hyphen.ai"},
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
			"region": "us-east",
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
