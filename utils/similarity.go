package utils

import (
	"errors"
	"math"
)

// CosineSimilarity computes the similarity between two vectors.
// Returns an error if vectors are different lengths.
func CosineSimilarity(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("vector length mismatch")
	}

	var dot, magA, magB float64

	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}

	denom := math.Sqrt(magA) * math.Sqrt(magB)
	if denom == 0 {
		return 0, errors.New("zero magnitude vector")
	}

	return dot / denom, nil
}
