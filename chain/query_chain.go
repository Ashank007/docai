package chain

import (
	//"strings"
  "fmt"
	//"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/retriever"
	//"github.com/Ashank007/docai/generator"
)

type QueryChain struct {
	EmbedFunc func(string) ([]float32, error)
	Retriever retriever.Retriever
	Generator func(string, []string) (string, error)
}

func (q *QueryChain) Run(query string, docNameFilter string) (string, error) { // <--- CORRECTED LINE HERE
	// Retriever.Retrieve needs to be updated to accept docNameFilter
	// This call is now correct as retriever.Retriever interface is fixed.
	chunks, err := q.Retriever.Retrieve(query, 4, docNameFilter)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err) // Use fmt.Errorf for better error wrapping
	}

	var contexts []string
	for _, c := range chunks {
		contexts = append(contexts, c.Chunk.Text)
	}

	return q.Generator(query, contexts)
}

