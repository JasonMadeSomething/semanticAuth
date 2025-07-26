package models

type Embedding struct {
	Hash   string    `bson:"hash"`
	Input  string    `bson:"input"`
	Vector []float64 `bson:"vector"`
}
