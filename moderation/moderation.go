package moderation

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

// ModerationResponse represents the response from the moderation service
type ModerationResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
	Version string `json:"version,omitempty"`
}

// HealthResponse represents the response from the moderation service health endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// Initialize checks the health of the moderation service and logs the result
func Initialize() {
	moderationBaseURL := os.Getenv("MODERATION_SERVICE_URL")
	if moderationBaseURL == "" {
		log.Println("WARNING: MODERATION_SERVICE_URL environment variable not set")
		return
	}

	// Ensure the URL doesn't end with a slash
	moderationBaseURL = strings.TrimSuffix(moderationBaseURL, "/")
	
	// Check health endpoint
	healthURL := fmt.Sprintf("%s/api/health", moderationBaseURL)
	resp, err := client.R().
		SetResult(&HealthResponse{}).
		Get(healthURL)

	if err != nil {
		log.Printf("WARNING: Moderation service health check failed: %v", err)
		return
	}

	if resp.StatusCode() >= 400 {
		log.Printf("WARNING: Moderation service health check returned status %d", resp.StatusCode())
		return
	}

	result := resp.Result().(*HealthResponse)
	log.Printf("Moderation service health check: Status=%s, Version=%s", result.Status, result.Version)
}

// CheckContent sends content to the moderation service and returns whether it's allowed
func CheckContent(content string) (*ModerationResponse, error) {
	// Get moderation service URL from environment variable
	moderationBaseURL := os.Getenv("MODERATION_SERVICE_URL")
	if moderationBaseURL == "" {
		return nil, fmt.Errorf("missing MODERATION_SERVICE_URL environment variable")
	}

	// Ensure the URL doesn't end with a slash
	moderationBaseURL = strings.TrimSuffix(moderationBaseURL, "/")
	
	// Construct the full URL with the /api/moderate endpoint
	moderationURL := fmt.Sprintf("%s/api/moderate", moderationBaseURL)

	// Prepare request body
	body := map[string]interface{}{
		"content":        content,
		"sending_system": "semantic-auth",
	}

	// Send request to moderation service
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&ModerationResponse{}).
		Post(moderationURL)

	if err != nil {
		return nil, fmt.Errorf("moderation service request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("moderation service returned error status: %d", resp.StatusCode())
	}

	result := resp.Result().(*ModerationResponse)
	return result, nil
}
