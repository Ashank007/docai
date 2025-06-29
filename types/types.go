package types

// Chunk is a unit of raw or processed text with optional metadata
type Chunk struct {
	ID       string // optional: unique ID or hash
	Text     string
	Source   string // filename or origin
	Page     int    // optional page number if from PDF
	Position int    // optional position/index in document
}

// ChunkWithEmbedding binds a chunk to its vector representation
type ChunkWithEmbedding struct {
	Chunk     Chunk
	Embedding []float32
}

// RetrievedChunk is the output of Retriever
type RetrievedChunk struct {
	Chunk     Chunk
	Embedding []float32 // optional: only if needed in downstream
	Score     float64   // similarity score
}

// FileMeta stores info about added files in the system
type FileMeta struct {
	ID       string // unique file ID
	Name     string // file name
	Path     string // absolute or relative path
	AddedAt  string // timestamp in RFC3339 format
	FileType string // e.g., pdf, txt
}
