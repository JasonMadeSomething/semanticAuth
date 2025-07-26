package models

import "time"

type LoginAttempt struct {
	Username   string    `bson:"username"`
	Input      string    `bson:"input"`
	Similarity float64   `bson:"similarity"`
	Timestamp  time.Time `bson:"timestamp"`
}
