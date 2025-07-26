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
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	collection := db.Client.Database("semantic_auth").Collection("users")

	// Check if user exists
	count, err := collection.CountDocuments(r.Context(), bson.M{"username": req.Username})
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	vec, err := openai.Embed(req.Password)
	if err != nil {
		http.Error(w, "Failed to embed password", http.StatusInternalServerError)
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
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
