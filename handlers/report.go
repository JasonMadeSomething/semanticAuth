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
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	// Username is now optional - if empty, we'll return all login attempts

	if req.Threshold == 0 {
		req.Threshold = 0.88 // default if omitted
	}

	coll := db.Client.Database("semantic_auth").Collection("login_attempts")

	// Create filter based on whether username is provided
	filter := bson.M{}
	if req.Username != "" {
		filter["username"] = req.Username
	}

	cursor, err := coll.Find(r.Context(),
		filter,
		options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(50),
	)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to query login attempts")
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

	RespondWithSuccess(w, "Login attempts retrieved successfully", results)
}
