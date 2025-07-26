package openai

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"semantic-auth/db"
	"semantic-auth/models"

	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client = resty.New()

func Embed(input string) ([]float64, error) {
	clean := strings.TrimSpace(strings.ToLower(input))
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(clean)))

	collection := db.Client.Database("semantic_auth").Collection("embeddings")

	// Check cache (Mongo)
	var cached models.Embedding
	err := collection.FindOne(context.TODO(), bson.M{"hash": hash}).Decode(&cached)
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

	// Cache it
	embedding := models.Embedding{
		Hash:   hash,
		Input:  clean,
		Vector: vector,
	}
	_, err = collection.InsertOne(context.TODO(), embedding)
	if err != nil {
		return nil, err
	}

	return vector, nil
}
