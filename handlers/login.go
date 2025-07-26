package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"semantic-auth/db"
	"semantic-auth/models"
	"semantic-auth/openai"
	"semantic-auth/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type LoginRequest struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Threshold float64 `json:"threshold"` // optional
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing username or password")
		return
	}

	// Set default threshold if not provided
	threshold := req.Threshold
	if threshold == 0 {
		threshold = 0.88
	}

	// Get user
	userColl := db.Client.Database("semantic_auth").Collection("users")
	var user models.User
	err = userColl.FindOne(r.Context(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	// Embed the guessed password
	guessVec, err := openai.Embed(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to embed password")
		log.Println("OpenAI error:", err)
		return
	}

	// Compare with stored
	similarity, err := utils.CosineSimilarity(user.Vector, guessVec)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Similarity calculation failed")
		return
	}

	// Log attempt
	attempt := models.LoginAttempt{
		Username:   req.Username,
		Input:      req.Password,
		Similarity: similarity,
		Timestamp:  time.Now(),
	}
	_, _ = db.Client.Database("semantic_auth").Collection("login_attempts").
		InsertOne(r.Context(), attempt)

	// Decide
	if similarity >= threshold {
		RespondWithSuccess(w, "Login successful", map[string]interface{}{
			"username": req.Username,
			"similarity": similarity,
			"threshold": threshold,
		})
	} else {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect password (not semantically similar enough)")
	}
}
