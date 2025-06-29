package chain

import (
	//"strings"

	//"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/retriever"
	//"github.com/Ashank007/docai/generator"
)

type QueryChain struct {
	EmbedFunc func(string) ([]float32, error)
	Retriever retriever.Retriever
	Generator func(string, []string) (string, error)
}

func (q *QueryChain) Run(query string) (string, error) {
	chunks, err := q.Retriever.Retrieve(query, 4)
	if err != nil {
		return "", err
	}

	var contexts []string
	for _, c := range chunks {
		contexts = append(contexts, c.Chunk.Text)
	}

	return q.Generator(query, contexts)
}
