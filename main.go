package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"semantic-auth/cache"
	"semantic-auth/db"
	"semantic-auth/handlers"
	"semantic-auth/moderation"
)

func main() {
	// Load environment variables from .env file only in local development
	// In production, environment variables should be set in the environment
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	// Connect to MongoDB
	db.Connect()

	// Initialize moderation service and check health
	moderation.Initialize()

	// Initialize semantic cache client
	cache.Initialize()

	// Setup router
	r := chi.NewRouter()

	// Get CORS allowed origins from environment variable
	// Format: comma-separated list of domains (e.g., "https://example.com,https://app.example.com")
	// Default to "*" (all origins) if not specified
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedOrigins := []string{"*"} // Default to all origins

	if corsOrigins != "" {
		// Split the comma-separated list into a slice
		allowedOrigins = strings.Split(corsOrigins, ",")
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))
	r.Use(middleware.Logger)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Register route
	r.Post("/register", handlers.RegisterHandler)

	// Login route
	r.Post("/login", handlers.LoginHandler)

	// Report route
	r.Get("/report", handlers.ReportHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start server
	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
