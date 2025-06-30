package retriever

import (
	"fmt"

	//"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/store"
	"github.com/Ashank007/docai/types"
)

type CosineRetriever struct {
	VectorDB   store.VectorStore
	MetaStore  store.MetadataStore
	EmbedFunc  func(string) ([]float32, error) // inject embedding logic
}

func NewCosineRetriever(vdb store.VectorStore, mdb store.MetadataStore, embed func(string) ([]float32, error)) *CosineRetriever {
	return &CosineRetriever{
		VectorDB:  vdb,
		MetaStore: mdb,
		EmbedFunc: embed,
	}
}

func (r *CosineRetriever) Retrieve(query string, topK int,docNameFilter string) ([]types.RetrievedChunk, error) {
	queryVec, err := r.EmbedFunc(query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	ids, err := r.VectorDB.SearchSimilar(queryVec, topK,docNameFilter)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	var results []types.RetrievedChunk
	for _, id := range ids {
		chunk, err := r.MetaStore.GetChunkByID(id)
		if err != nil {
			continue
		}
		results = append(results, types.RetrievedChunk{
			Chunk:     chunk,
			Embedding: nil, // optional
			Score:     0,   // optional
		})
	}

	return results, nil
}
