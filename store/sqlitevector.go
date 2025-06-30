package store

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"sort"
	"sync"
  "log"
	"github.com/Ashank007/docai/utils"
)


type SQLiteVectorStore struct {
	db  *sql.DB
	mu  sync.RWMutex
	mem map[int64]struct {
		Vector  []float32
		DocName string
	}
}


func NewSQLiteVectorStore(db *sql.DB) (*SQLiteVectorStore, error) {
	vs := &SQLiteVectorStore{
		db:  db,
		mem: make(map[int64]struct{ Vector []float32; DocName string }),
	}
	_, err := vs.db.Exec(`
		CREATE TABLE IF NOT EXISTS vectors (
			id INTEGER PRIMARY KEY,
			doc_name TEXT NOT NULL,
			vector BLOB NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create vectors table: %w", err)
	}
	err = vs.loadCache()
	return vs, err
}


func (s *SQLiteVectorStore) loadCache() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(`SELECT id, doc_name, vector FROM vectors`) // Select doc_name as well
	if err != nil {
		return fmt.Errorf("failed to load vectors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var docName string // To read doc_name
		var blob []byte
		if err := rows.Scan(&id, &docName, &blob); err != nil {
			log.Printf("Error scanning row in loadCache: %v", err) // Log error instead of continuing silently
			continue
		}
		var vec []float32
		buf := bytes.NewReader(blob)
		if err := gob.NewDecoder(buf).Decode(&vec); err != nil {
			log.Printf("Error decoding vector for id %d: %v", id, err) // Log error
			continue
		}
		s.mem[id] = struct{ Vector []float32; DocName string }{Vector: vec, DocName: docName}
	}
	return nil
}

func (s *SQLiteVectorStore) AddVector(id int64, vec []float32, docName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(vec); err != nil {
		return err
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO vectors (id, doc_name, vector)
		VALUES (?, ?, ?)
	`, id, docName, b.Bytes())

	if err != nil {
		return fmt.Errorf("failed to insert vector: %w", err)
	}

	s.mem[id] = struct{ Vector []float32; DocName string }{Vector: vec, DocName: docName}
	return nil
}


func (s *SQLiteVectorStore) SearchSimilar(query []float32, topK int, docNameFilter string) ([]int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		id    int64
		score float64
	}
	var results []scored

	for id, data := range s.mem {
		if docNameFilter != "" && data.DocName != docNameFilter {
			continue // Skip this vector if it doesn't match the filter
		}

		sim, err := utils.CosineSimilarity(query, data.Vector) // Use data.Vector
		if err != nil {
			log.Printf("Error calculating cosine similarity for id %d: %v", id, err)
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

	s.mem = make(map[int64]struct{ Vector []float32; DocName string }) // Reset in-memory map
	_, err := s.db.Exec(`DELETE FROM vectors`)
	return err
}

func (s *SQLiteVectorStore) DeleteVectorsByDoc(docName string) error { // Renamed parameter for clarity
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM vectors WHERE doc_name = ?`, docName)
	if err != nil {
		return fmt.Errorf("failed to delete vectors from DB for doc %s: %w", docName, err)
	}

	// remove from memory
	for id, data := range s.mem {
		if data.DocName == docName {
			delete(s.mem, id)
		}
	}
	return nil
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}
