package store

import "github.com/Ashank007/docai/types"

// MetadataStore manages files and their associated chunk metadata
type MetadataStore interface {
	Init(path string) error
	SaveChunk(docName string, chunk types.Chunk) (int64, error)
	GetChunkByID(id int64) (types.Chunk, error)
	ListFiles() ([]types.FileMeta, error)
	DeleteFile(name string) error
	Close() error
}

// VectorStore manages vector representations and similarity search
type VectorStore interface {
	AddVector(id int64, vec []float32, docName string) error
	SearchSimilar(query []float32, topK int,docNameFilter string) ([]int64, error)
	Reset() error
	DeleteVectorsByDoc(docName string) error
}
