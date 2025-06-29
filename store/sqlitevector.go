package store

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"sort"
	"sync"

	"github.com/Ashank007/docai/utils"
)

type SQLiteVectorStore struct {
	db  *sql.DB
	mu  sync.RWMutex
	mem map[int64][]float32 // in-memory cache for fast search
}

func NewSQLiteVectorStore(db *sql.DB) (*SQLiteVectorStore, error) {
	vs := &SQLiteVectorStore{
		db:  db,
		mem: make(map[int64][]float32),
	}
	err := vs.loadCache()
	return vs, err
}

func (s *SQLiteVectorStore) loadCache() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`SELECT id, vector FROM vectors`)
	if err != nil {
		return fmt.Errorf("failed to load vectors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var blob []byte
		if err := rows.Scan(&id, &blob); err != nil {
			continue
		}
		var vec []float32
		buf := bytes.NewReader(blob)
		if err := gob.NewDecoder(buf).Decode(&vec); err != nil {
			continue
		}
		s.mem[id] = vec
	}
	return nil
}

func (s *SQLiteVectorStore) AddVector(id int64, vec []float32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(vec); err != nil {
		return err
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO vectors (id, doc_name, vector)
		VALUES (?, (SELECT doc_name FROM chunks WHERE id = ?), ?)
	`, id, id, b.Bytes())

	if err != nil {
		return fmt.Errorf("failed to insert vector: %w", err)
	}

	s.mem[id] = vec
	return nil
}

func (s *SQLiteVectorStore) SearchSimilar(query []float32, topK int) ([]int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		id    int64
		score float64
	}
	var results []scored

	for id, vec := range s.mem {
		sim, err := utils.CosineSimilarity(query, vec)
		if err != nil {
			continue
		}
		results = append(results, scored{id: id, score: sim})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	limit := topK
	if len(results) < topK {
		limit = len(results)
	}

	ids := make([]int64, limit)
	for i := 0; i < limit; i++ {
		ids[i] = results[i].id
	}
	return ids, nil
}

func (s *SQLiteVectorStore) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mem = make(map[int64][]float32)
	_, err := s.db.Exec(`DELETE FROM vectors`)
	return err
}

func (s *SQLiteVectorStore) DeleteVectorsByDoc(doc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// remove from SQLite
	_, err := s.db.Exec(`DELETE FROM vectors WHERE doc_name = ?`, doc)
	if err != nil {
		return err
	}

	// remove from memory
	for id, _ := range s.mem {
		var docName string
		err := s.db.QueryRow(`SELECT doc_name FROM chunks WHERE id = ?`, id).Scan(&docName)
		if err == nil && docName == doc {
			delete(s.mem, id)
		}
	}
	return nil
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}
