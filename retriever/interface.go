package retriever

import "github.com/Ashank007/docai/types"

type Retriever interface {
  Retrieve(query string, topK int, docNameFilter string) ([]types.RetrievedChunk, error)
}
