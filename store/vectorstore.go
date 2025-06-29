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
	vectors [][]float32
	ids     []int64
}

func NewMemoryVectorStore() *MemoryVectorStore {
	return &MemoryVectorStore{
		vectors: make([][]float32, 0),
		ids:     make([]int64, 0),
	}
}

func (m *MemoryVectorStore) AddVector(id int64, vec []float32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ids = append(m.ids, id)
	m.vectors = append(m.vectors, vec)
	return nil
}

func (m *MemoryVectorStore) SearchSimilar(query []float32, topK int) ([]int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type result struct {
		id    int64
		score float64
	}
	var results []result
	fmt.Println("vectors: ",m.vectors)
	for i, vec := range m.vectors {
		sim, err := utils.CosineSimilarity(query, vec)
		if err != nil {
			continue
		}
		results = append(results, result{id: m.ids[i], score: sim})
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
	m.ids = []int64{}
	m.vectors = [][]float32{}
	return nil
}
