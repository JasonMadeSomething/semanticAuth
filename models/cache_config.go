package models

// CacheConfig represents the configuration for the semantic cache
type CacheConfig struct {
	Enabled           bool    `json:"enabled"`
	URL               string  `json:"url"`
	SimilarityThreshold float64 `json:"similarity_threshold"`
	AllowFallback     bool    `json:"allow_fallback"`
}

// DefaultCacheConfig returns the default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Enabled:           false,
		URL:               "http://localhost:8081",
		SimilarityThreshold: 0.88,
		AllowFallback:     true,
	}
}
