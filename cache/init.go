package cache

import (
	"log"
	"os"
	"strconv"

	"semantic-auth/models"
)

var (
	// DefaultClient is the default semantic cache client
	DefaultClient *Client
)

// Initialize initializes the semantic cache client
func Initialize() {
	config := models.DefaultCacheConfig()

	// Check if cache is enabled from environment variable
	if enabledStr := os.Getenv("SEMANTIC_CACHE_ENABLED"); enabledStr != "" {
		enabled, err := strconv.ParseBool(enabledStr)
		if err != nil {
			log.Printf("Warning: Invalid SEMANTIC_CACHE_ENABLED value: %s, defaulting to %v", enabledStr, config.Enabled)
		} else {
			config.Enabled = enabled
		}
	}

	// Get cache URL from environment variable
	if url := os.Getenv("SEMANTIC_CACHE_URL"); url != "" {
		config.URL = url
	}

	// Get similarity threshold from environment variable
	if thresholdStr := os.Getenv("SEMANTIC_CACHE_THRESHOLD"); thresholdStr != "" {
		threshold, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			log.Printf("Warning: Invalid SEMANTIC_CACHE_THRESHOLD value: %s, defaulting to %v", thresholdStr, config.SimilarityThreshold)
		} else {
			config.SimilarityThreshold = threshold
		}
	}

	// Get allow fallback from environment variable
	if fallbackStr := os.Getenv("SEMANTIC_CACHE_ALLOW_FALLBACK"); fallbackStr != "" {
		allowFallback, err := strconv.ParseBool(fallbackStr)
		if err != nil {
			log.Printf("Warning: Invalid SEMANTIC_CACHE_ALLOW_FALLBACK value: %s, defaulting to %v", fallbackStr, config.AllowFallback)
		} else {
			config.AllowFallback = allowFallback
		}
	}

	// Create the client
	DefaultClient = NewClient(config)

	if DefaultClient.IsEnabled() {
		log.Printf("Semantic cache enabled at %s", config.URL)
		
		// Check if the cache is healthy
		if DefaultClient.HealthCheck(nil) {
			log.Println("Semantic cache is healthy")
		} else {
			log.Println("Warning: Semantic cache is not healthy, but will continue with fallback to direct OpenAI calls")
		}
	} else {
		log.Println("Semantic cache is disabled")
	}
}
