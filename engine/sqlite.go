package engine

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./rag_data.db")
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS chunks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		doc_name TEXT,
		chunk_text TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func SaveChunk(docName, chunkText string) (int64, error) {
	res, err := db.Exec("INSERT INTO chunks(doc_name, chunk_text) VALUES (?, ?)", docName, chunkText)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetChunkByID(id int64) (string, error) {
	var text string
	err := db.QueryRow("SELECT chunk_text FROM chunks WHERE id = ?", id).Scan(&text)
	return text, err
}


