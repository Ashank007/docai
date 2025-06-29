package utils

import (
	"errors"
	"math"
)

// CosineSimilarity returns similarity between two vectors (0 to 1 scale).
func CosineSimilarity(vec1, vec2 []float32) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, errors.New("vectors must be of equal length")
	}

	dot := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for i := 0; i < len(vec1); i++ {
		v1 := float64(vec1[i])
		v2 := float64(vec2[i])
		dot += v1 * v2
		norm1 += v1 * v1
		norm2 += v2 * v2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0, errors.New("zero-vector passed")
	}

	return dot / (math.Sqrt(norm1) * math.Sqrt(norm2)), nil
}

// DotProduct returns the dot product of two vectors.
func DotProduct(vec1, vec2 []float32) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, errors.New("vectors must be of equal length")
	}
	dot := 0.0
	for i := 0; i < len(vec1); i++ {
		dot += float64(vec1[i]) * float64(vec2[i])
	}
	return dot, nil
}

// L2Norm returns the L2 norm (magnitude) of a vector.
func L2Norm(vec []float32) float64 {
	sum := 0.0
	for _, v := range vec {
		sum += float64(v) * float64(v)
	}
	return math.Sqrt(sum)
}

// Normalize returns unit vector (optional utility)
func Normalize(vec []float32) ([]float64, error) {
	norm := L2Norm(vec)
	if norm == 0 {
		return nil, errors.New("cannot normalize zero vector")
	}
	result := make([]float64, len(vec))
	for i, v := range vec {
		result[i] = float64(v) / norm
	}
	return result, nil
}
