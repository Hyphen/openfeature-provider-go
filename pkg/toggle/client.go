package toggle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	cache      *cache.Cache
	config     Config
	publicKey  string
	keyGen     func(ctx EvaluationContext) string
	serverURLs []string
}

func newClient(config Config) (*Client, error) {
	urls := config.HorizonServerURLs
	if len(urls) == 0 {
		urls = []string{DefaultHorizonURL}
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		config:     config,
		publicKey:  config.PublicKey,
		serverURLs: urls,
	}

	if config.Cache != nil {
		c.cache = cache.New(config.Cache.TTL, 10*time.Minute)
		c.keyGen = config.Cache.KeyGen
	}

	return c, nil
}

func (c *Client) Evaluate(ctx EvaluationContext) (*Response, error) {
	if c.cache != nil && c.keyGen != nil {
		key := c.keyGen(ctx)
		if cached, found := c.cache.Get(key); found {
			return cached.(*Response), nil
		}
	}

	var lastErr error
	for _, url := range c.serverURLs {
		resp, err := c.fetchEvaluation(url, ctx)
		if err != nil {
			lastErr = err
			continue
		}

		if c.cache != nil && c.keyGen != nil {
			key := c.keyGen(ctx)
			c.cache.Set(key, resp, cache.DefaultExpiration)
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all evaluation attempts failed: %v", lastErr)
}

func (c *Client) fetchEvaluation(url string, ctx EvaluationContext) (*Response, error) {
	payload, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url+"/evaluate", bytes.NewBuffer(payload))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) SendTelemetry(payload TelemetryPayload) error {
	for _, url := range c.serverURLs {
		if err := c.postTelemetry(url, payload); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to send telemetry to all endpoints")
}

func (c *Client) postTelemetry(url string, payload TelemetryPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url+"/telemetry", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.publicKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}
