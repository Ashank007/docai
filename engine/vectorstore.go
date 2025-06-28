package engine

import (
	"math"
	"sync"
)

type VectorStore struct {
	mu      sync.Mutex
	vectors [][]float32
	ids     []int64
}

var store *VectorStore

func InitVectorStore() {
	store = &VectorStore{
		vectors: make([][]float32, 0),
		ids:     make([]int64, 0),
	}
}

func AddVector(id int64, vec []float32) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.ids = append(store.ids, id)
	store.vectors = append(store.vectors, vec)
}

func cosineSimilarity(a, b []float32) float64 {
	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func SearchVector(query []float32, topK int) []int64 {
	store.mu.Lock()
	defer store.mu.Unlock()

	type result struct {
		id    int64
		score float64
	}

	results := make([]result, 0, len(store.ids))
	for i, vec := range store.vectors {
		score := cosineSimilarity(query, vec)
		results = append(results, result{id: store.ids[i], score: score})
	}

	// Simple bubble sort by score descending
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	limit := topK
	if len(results) < topK {
		limit = len(results)
	}

	topIDs := make([]int64, limit)
	for i := 0; i < limit; i++ {
		topIDs[i] = results[i].id
	}

	return topIDs
}


