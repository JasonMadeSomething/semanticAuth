package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"semantic-auth/db"
	"semantic-auth/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportRequest struct {
	Username  string  `json:"username"`
	Threshold float64 `json:"threshold,omitempty"`
}

type ReportResponse struct {
	Input      string    `json:"input"`
	Similarity float64   `json:"similarity"`
	Timestamp  time.Time `json:"timestamp"`
	Passed     bool      `json:"passed"`
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.Username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	if req.Threshold == 0 {
		req.Threshold = 0.88 // default if omitted
	}

	coll := db.Client.Database("semantic_auth").Collection("login_attempts")

	cursor, err := coll.Find(r.Context(),
		bson.M{"username": req.Username},
		options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(20),
	)
	if err != nil {
		http.Error(w, "Failed to query login attempts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var results []ReportResponse
	for cursor.Next(r.Context()) {
		var attempt models.LoginAttempt
		if err := cursor.Decode(&attempt); err != nil {
			continue
		}

		results = append(results, ReportResponse{
			Input:      attempt.Input,
			Similarity: attempt.Similarity,
			Timestamp:  attempt.Timestamp,
			Passed:     attempt.Similarity >= req.Threshold,
		})
	}

	json.NewEncoder(w).Encode(results)
}
