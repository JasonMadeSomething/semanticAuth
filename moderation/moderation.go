package moderation

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

// ModerationResponse represents the response from the moderation service
type ModerationResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
	Version string `json:"version,omitempty"`
}

// CheckContent sends content to the moderation service and returns whether it's allowed
func CheckContent(content string) (*ModerationResponse, error) {
	// Get moderation service URL from environment variable
	moderationURL := os.Getenv("MODERATION_SERVICE_URL")
	if moderationURL == "" {
		return nil, fmt.Errorf("missing MODERATION_SERVICE_URL environment variable")
	}

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
