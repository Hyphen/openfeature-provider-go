package toggle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	cache      *Cache
	config     Config
	publicKey  string
}

func newClient(config Config) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		config:    config,
		publicKey: config.PublicKey,
	}

	if config.Cache != nil {
		c.cache = newCache(config.Cache)
	}

	return c, nil
}

func (c *Client) Evaluate(ctx EvaluationContext) (*Response, error) {
	if c.cache != nil {
		if cached := c.cache.Get(ctx); cached != nil {
			return cached.(*Response), nil
		}
	}

	payload, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.config.HorizonServerURL+"/evaluate", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.publicKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if c.cache != nil {
		c.cache.Set(ctx, &result)
	}

	return &result, nil
}

func (c *Client) SendTelemetry(payload TelemetryPayload) error {
	// Implementation for telemetry...
	return nil
}
