package store

import (
	"database/sql"
	"fmt"

	"github.com/Ashank007/docai/types"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore() *SQLiteStore {
	return &SQLiteStore{}
}

func (s *SQLiteStore) Init(path string) error {
	var err error
	s.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open sqlite DB: %w", err)
	}
	stmt := `
	CREATE TABLE IF NOT EXISTS chunks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		doc_name TEXT,
		chunk_text TEXT,
		page INT,
		position INT
	);
	`
	_, err = s.db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create chunk table: %w", err)
	}
	stmt2 := `
		CREATE TABLE IF NOT EXISTS vectors (
			id INTEGER PRIMARY KEY,
			doc_name TEXT,
			vector BLOB
		);
		`
		_, err = s.db.Exec(stmt2)
	if err != nil {
		return fmt.Errorf("failed to create vectors table: %w", err)
	}
	return err
}

func (s *SQLiteStore) SaveChunk(docName string, chunk types.Chunk) (int64, error) {
	res, err := s.db.Exec(`
	INSERT INTO chunks (doc_name, chunk_text, page, position)
	VALUES (?, ?, ?, ?)`, docName, chunk.Text, chunk.Page, chunk.Position)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *SQLiteStore) GetChunkByID(id int64) (types.Chunk, error) {
	var chunk types.Chunk
	err := s.db.QueryRow(`
	SELECT chunk_text, doc_name, page, position FROM chunks WHERE id = ?`, id).
		Scan(&chunk.Text, &chunk.Source, &chunk.Page, &chunk.Position)
	chunk.ID = fmt.Sprintf("%d", id)
	return chunk, err
}

func (s *SQLiteStore) ListFiles() ([]types.FileMeta, error) {
	rows, err := s.db.Query("SELECT DISTINCT doc_name FROM chunks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []types.FileMeta
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		files = append(files, types.FileMeta{Name: name})
	}
	return files, nil
}


func (s *SQLiteStore) DeleteFile(name string) error {
	_, err := s.db.Exec(`DELETE FROM vectors WHERE doc_name = ?`, name)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM chunks WHERE doc_name = ?`, name)
	return err
}


func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
