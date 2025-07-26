package main

import (
	"log"

	"semantic-auth/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect()
	// Your logic here: routes, inserts, etc.
}
