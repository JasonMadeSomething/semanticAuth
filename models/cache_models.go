package models

// CacheRequest represents the incoming request to the cache API
type CacheRequest struct {
	Text          string `json:"text" binding:"required"`
	SourceSystem  string `json:"source_system" binding:"required"`
	AllowFallback bool   `json:"allow_fallback,omitempty"`
	Response      string `json:"response,omitempty"` // Optional response to store with this input
}

// CacheResponse represents the response from the cache API
type CacheResponse struct {
	Cached    bool    `json:"cached"`
	Response  string  `json:"response,omitempty"`
	SectorKey string  `json:"sector_key,omitempty"`
	Similarity float64 `json:"similarity,omitempty"`
}
