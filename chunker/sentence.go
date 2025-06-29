package chunker

import (
	"strings"

	"github.com/Ashank007/docai/types"
)

type SentenceChunker struct {
	MaxWords int
}

func NewSentenceChunker(maxWords int) *SentenceChunker {
	if maxWords <= 0 {
		maxWords = 200
	}
	return &SentenceChunker{MaxWords: maxWords}
}

func (sc *SentenceChunker) Chunk(text string) ([]types.Chunk, error) {
	rawSentences := strings.Split(text, ".")
	var chunks []types.Chunk
	var currentChunk []string
	wordCount := 0

	for _, s := range rawSentences {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		sentence := s + "."
		sentenceWords := len(strings.Fields(sentence))

		if wordCount+sentenceWords > sc.MaxWords && len(currentChunk) > 0 {
			chunks = append(chunks, types.Chunk{
				Text: strings.Join(currentChunk, " "),
			})
			currentChunk = []string{}
			wordCount = 0
		}

		currentChunk = append(currentChunk, sentence)
		wordCount += sentenceWords
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, types.Chunk{
			Text: strings.Join(currentChunk, " "),
		})
	}

	return chunks, nil
}

func (sc *SentenceChunker) Name() string {
	return "sentence-word-chunker"
}
