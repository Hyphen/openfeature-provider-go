# Hyphen OpenFeature Provider for Go

[![Latest Release](https://img.shields.io/github/v/release/hyphen/openfeature-provider-go)](https://github.com/hyphen/openfeature-provider-go/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/hyphen/openfeature-provider-go.svg)](https://pkg.go.dev/github.com/hyphen/openfeature-provider-go)
[![License](https://img.shields.io/github/license/hyphen/openfeature-provider-go)](https://github.com/hyphen/openfeature-provider-go/blob/main/LICENSE)

This repository contains the Hyphen Provider implementation for the [OpenFeature](https://openfeature.dev) Go SDK.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Configuration](#configuration)
4. [Development](#development)
5. [Contributing](#contributing)
6. [License](#license)

## Installation

```bash
go get github.com/hyphen/openfeature-provider-go
```

## Usage

```go
package main

import (
    "context"
    "github.com/open-feature/go-sdk/openfeature"
    "github.com/hyphen/openfeature-provider-go/pkg/toggle"
)

func main() {
    // Configure the provider
    config := toggle.Config{
        PublicKey:   "your-public-key",
        Application: "your-app",
        Environment: "development",
    }

    // Create a new provider
    provider, err := toggle.NewProvider(config)
    if err != nil {
        panic(err)
    }

    // Set as global provider
    openfeature.SetProvider(provider)

    // Get a client
    client := openfeature.NewClient("my-app")

    // Evaluate flags
    value, err := client.BooleanValue(
        context.Background(),
        "my-flag",
        false,
        openfeature.NewEvaluationContext("user-123", nil),
    )
}
```

## Configuration

### Provider Options

| Option              | Type     | Description                                                                           |
|--------------------|----------|---------------------------------------------------------------------------------------|
| `PublicKey`        | string   | Your Hyphen API public key                                                            |
| `Application`      | string   | The application id or alternate id                                                    |
| `Environment`      | string   | The environment in which your application is running (e.g., `production`, `staging`)  |
| `EnableUsage`      | bool     | Enable or disable the logging of toggle usage (telemetry)                            |
| `Cache`            | object   | Configuration for caching feature flag evaluations                                    |

### Caching

The provider supports caching of evaluation results:

```go
config := toggle.Config{
    PublicKey:   "your-public-key",
    Application: "your-app",
    Environment: "development",
    Cache: &toggle.CacheConfig{
        TTL: time.Minute * 5,
        KeyGen: func(ctx toggle.EvaluationContext) string {
            return ctx.TargetingKey
        },
    },
}
```

## Development

### Requirements

- Go 1.19 or higher

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build ./...
```

## Contributing

We welcome contributions to this project! If you'd like to contribute, please follow the guidelines outlined in [CONTRIBUTING.md](CONTRIBUTING.md). Whether it's reporting issues, suggesting new features, or submitting pull requests, your help is greatly appreciated!

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for full details.
