package chunker

import "github.com/Ashank007/docai/types"

type Chunker interface {
	Chunk(text string) ([]types.Chunk, error)
	Name() string
}
