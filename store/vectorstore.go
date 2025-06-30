package store

import (
	//"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/Ashank007/docai/utils"
)

type MemoryVectorStore struct {
	mu      sync.RWMutex
	data []struct {
		ID      int64
		Vector  []float32
		DocName string // New field to store the document name
	}

}

func NewMemoryVectorStore() *MemoryVectorStore {
	return &MemoryVectorStore{
		data: make([]struct{ ID int64; Vector []float32; DocName string }, 0),
	}
}

func (m *MemoryVectorStore) AddVector(id int64, vec []float32, docName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if ID already exists and update, or append
	found := false
	for i, item := range m.data {
		if item.ID == id {
			m.data[i].Vector = vec
			m.data[i].DocName = docName
			found = true
			break
		}
	}
	if !found {
		m.data = append(m.data, struct {
			ID      int64
			Vector  []float32
			DocName string
		}{ID: id, Vector: vec, DocName: docName})
	}

	return nil
}

func (m *MemoryVectorStore) SearchSimilar(query []float32, topK int, docNameFilter string) ([]int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type result struct {
		id    int64
		score float64
	}
	var results []result

	fmt.Println("vectors: ",m.data) // Changed from m.vectors to m.data to show full stored info
	for _, item := range m.data { // Iterate over the new data slice
		// Apply docNameFilter: if filter is not empty, check for match
		if docNameFilter != "" && item.DocName != docNameFilter {
			continue // Skip this vector if it doesn't match the filter
		}

		sim, err := utils.CosineSimilarity(query, item.Vector)
		if err != nil {
			// Log error if needed, but continue processing other vectors
			// fmt.Printf("Error calculating cosine similarity for id %d: %v\n", item.ID, err)
			continue
		}
		results = append(results, result{id: item.ID, score: sim})
	}
	fmt.Println("results: ",results)

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	limit := topK
	if len(results) < topK {
		limit = len(results)
	}

	topIDs := make([]int64, limit)
	for i := 0; i < limit; i++ {
		topIDs[i] = results[i].id
	}

	return topIDs, nil
}

func (m *MemoryVectorStore) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = []struct{ ID int64; Vector []float32; DocName string }{} // Clear the new data slice
	return nil
}

func (m *MemoryVectorStore) DeleteVectorsByDoc(docName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var newData []struct {
		ID      int64
		Vector  []float32
		DocName string
	}
	for _, item := range m.data {
		if item.DocName != docName {
			newData = append(newData, item)
		}
	}
	m.data = newData
	return nil
}
