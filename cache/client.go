package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"semantic-auth/models"

	"github.com/go-resty/resty/v2"
)

// Client represents a client for the semantic cache service
type Client struct {
	config models.CacheConfig
	client *resty.Client
}

// NewClient creates a new semantic cache client
func NewClient(config models.CacheConfig) *Client {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(1)

	return &Client{
		config: config,
		client: client,
	}
}

// GetEmbedding attempts to get an embedding from the cache
// If the cache is not enabled or fails, it returns nil and an error
// The error should be logged but can be ignored if fallback is allowed
func (c *Client) GetEmbedding(ctx context.Context, input string) ([]float64, error) {
	if !c.config.Enabled {
		return nil, fmt.Errorf("cache is disabled")
	}

	req := models.CacheRequest{
		Text:          input,
		SourceSystem:  "semanticAuth",
		AllowFallback: c.config.AllowFallback,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("%s/cache", c.config.URL))

	if err != nil {
		return nil, fmt.Errorf("cache request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("cache returned non-OK status: %d - %s", resp.StatusCode(), resp.String())
	}

	var cacheResp models.CacheResponse
	if err := json.Unmarshal(resp.Body(), &cacheResp); err != nil {
		return nil, fmt.Errorf("failed to parse cache response: %w", err)
	}

	if !cacheResp.Cached {
		return nil, fmt.Errorf("no cached embedding found")
	}

	// Parse the response string as a JSON array of float64
	var vector []float64
	if err := json.Unmarshal([]byte(cacheResp.Response), &vector); err != nil {
		return nil, fmt.Errorf("failed to parse cached vector: %w", err)
	}

	return vector, nil
}

// StoreEmbedding stores an embedding in the cache
// This is a best-effort operation - errors are logged but not returned
func (c *Client) StoreEmbedding(ctx context.Context, input string, vector []float64) {
	if !c.config.Enabled {
		return
	}

	// Convert the vector to a JSON string
	vectorJSON, err := json.Marshal(vector)
	if err != nil {
		log.Printf("Failed to marshal vector for cache storage: %v", err)
		return
	}

	req := models.CacheRequest{
		Text:         input,
		SourceSystem: "semanticAuth",
		Response:     string(vectorJSON),
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("%s/cache", c.config.URL))

	if err != nil {
		log.Printf("Failed to store embedding in cache: %v", err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Cache returned non-OK status when storing: %d - %s", resp.StatusCode(), resp.String())
	}
}

// IsEnabled returns whether the cache is enabled
func (c *Client) IsEnabled() bool {
	return c.config.Enabled
}

// HealthCheck checks if the cache service is healthy
func (c *Client) HealthCheck(ctx context.Context) bool {
	if !c.config.Enabled {
		return false
	}

	resp, err := c.client.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/cache/health", c.config.URL))

	if err != nil {
		log.Printf("Cache health check failed: %v", err)
		return false
	}

	return resp.StatusCode() == http.StatusOK
}
