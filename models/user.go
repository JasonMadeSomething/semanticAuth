package models

type User struct {
	Username string    `bson:"username"`
	Hash     string    `bson:"hash"`
	Vector   []float64 `bson:"vector"`
	Raw      string    `bson:"raw,omitempty"`
}
