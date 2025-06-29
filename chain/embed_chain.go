package chain

import (
	"fmt"

	"github.com/Ashank007/docai/chunker"
	//"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/store"
	//"github.com/Ashank007/docai/types"
)

type EmbedChain struct {
	DocName    string
	Chunker    chunker.Chunker
	EmbedFunc  func(string) ([]float32, error)
	MetaStore  store.MetadataStore
	VectorDB   store.VectorStore
}

func (e *EmbedChain) Run(input string) (string, error) {
	chunks, _ := e.Chunker.Chunk(input)
	for _, chunk := range chunks {
		id, err := e.MetaStore.SaveChunk(e.DocName, chunk)
		if err != nil {
			return "", fmt.Errorf("failed to save chunk: %w", err)
		}
		vec, err := e.EmbedFunc(chunk.Text)
		if err != nil {
			return "", fmt.Errorf("embedding error: %w", err)
		}
		err = e.VectorDB.AddVector(id, vec)
		if err != nil {
			return "", fmt.Errorf("vector store error: %w", err)
		}
	}
	return fmt.Sprintf("Document '%s' embedded successfully.", e.DocName), nil
}
