package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"semantic-auth/db"
	"semantic-auth/models"
	"semantic-auth/openai"

	"go.mongodb.org/mongo-driver/bson"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
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

	collection := db.Client.Database("semantic_auth").Collection("users")

	// Check if user exists
	count, err := collection.CountDocuments(r.Context(), bson.M{"username": req.Username})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Database error occurred")
		return
	}
	if count > 0 {
		RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	vec, err := openai.Embed(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to embed password")
		log.Println("OpenAI error:", err)
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ToLower(req.Password))))

	user := models.User{
		Username: req.Username,
		Hash:     hash,
		Vector:   vec,
		Raw:      req.Password, // optional, remove if you want to be pure
	}

	_, err = collection.InsertOne(r.Context(), user)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to store user")
		return
	}

	RespondWithSuccess(w, "User registered successfully", map[string]interface{}{
		"username": req.Username,
	})
}
