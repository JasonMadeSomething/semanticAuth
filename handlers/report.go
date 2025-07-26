package handlers

import (
	"net/http"
	"strconv"
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
	// Get threshold from query parameter
	thresholdStr := r.URL.Query().Get("threshold")
	threshold := 0.88 // default threshold

	if thresholdStr != "" {
		parsedThreshold, err := strconv.ParseFloat(thresholdStr, 64)
		if err == nil && parsedThreshold >= 0.5 && parsedThreshold <= 1.0 {
			threshold = parsedThreshold
		}
	}

	// Get username from query parameter (optional)
	username := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("username")))

	coll := db.Client.Database("semantic_auth").Collection("login_attempts")

	// Create filter based on whether username is provided
	filter := bson.M{}
	if username != "" {
		filter["username"] = username
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
			Passed:     attempt.Similarity >= threshold,
		})
	}

	RespondWithSuccess(w, "Login attempts retrieved successfully", results)
}
