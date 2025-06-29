package retriever

import "github.com/Ashank007/docai/types"

type Retriever interface {
	Retrieve(query string, topK int) ([]types.RetrievedChunk, error)
}