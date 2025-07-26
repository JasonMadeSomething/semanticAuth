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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
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
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Embed the guessed password
	guessVec, err := openai.Embed(req.Password)
	if err != nil {
		http.Error(w, "Failed to embed password", http.StatusInternalServerError)
		log.Println("OpenAI error:", err)
		return
	}

	// Compare with stored
	similarity, err := utils.CosineSimilarity(user.Vector, guessVec)
	if err != nil {
		http.Error(w, "Similarity calculation failed", http.StatusInternalServerError)
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
		w.Write([]byte("Login successful"))
	} else {
		http.Error(w, "Incorrect password (not semantically similar enough)", http.StatusUnauthorized)
	}
}
