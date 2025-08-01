package openai

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"semantic-auth/cache"
	"semantic-auth/db"
	"semantic-auth/models"
	"semantic-auth/moderation"

	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client = resty.New()

func Embed(input string) ([]float64, error) {
	clean := strings.TrimSpace(strings.ToLower(input))
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(clean)))

	// Check content with moderation service
	modResp, err := moderation.CheckContent(clean)
	if err != nil {
		return nil, fmt.Errorf("moderation check failed: %w", err)
	}

	// If content is not allowed, return error
	if !modResp.Allowed {
		message := "Content not allowed by moderation service"
		if modResp.Message != "" {
			message = modResp.Message
		}
		return nil, fmt.Errorf("moderation error: %s", message)
	}

	// Try to get embedding from external semantic cache if enabled
	if cache.DefaultClient != nil && cache.DefaultClient.IsEnabled() {
		vector, err := cache.DefaultClient.GetEmbedding(context.Background(), clean)
		if err == nil {
			// Successfully retrieved from external cache
			log.Printf("Retrieved embedding from semantic cache for input: %s", clean)
			return vector, nil
		} else {
			// Log the error but continue with fallback
			log.Printf("Semantic cache retrieval failed: %v, falling back to local cache/OpenAI", err)
		}
	}

	collection := db.Client.Database("semantic_auth").Collection("embeddings")

	// Check local cache (MongoDB)
	var cached models.Embedding
	err = collection.FindOne(context.TODO(), bson.M{"hash": hash}).Decode(&cached)
	if err == nil {
		return cached.Vector, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Hit OpenAI
	openaiKey := os.Getenv("OPENAI_KEY")
	if openaiKey == "" {
		return nil, fmt.Errorf("missing OPENAI_KEY")
	}

	body := map[string]interface{}{
		"input": clean,
		"model": "text-embedding-3-small",
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openaiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("https://api.openai.com/v1/embeddings")

	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	vector := result.Data[0].Embedding

	// Cache it locally
	embedding := models.Embedding{
		Hash:   hash,
		Input:  clean,
		Vector: vector,
	}
	_, err = collection.InsertOne(context.TODO(), embedding)
	if err != nil {
		log.Printf("Warning: Failed to save embedding to local cache: %v", err)
		// Continue despite the error
	}

	// Store in external semantic cache if enabled
	if cache.DefaultClient != nil && cache.DefaultClient.IsEnabled() {
		go func() {
			cache.DefaultClient.StoreEmbedding(context.Background(), clean, vector)
		}()
	}

	return vector, nil
}
