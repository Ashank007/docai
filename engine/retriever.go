package engine

import (
	"log"
)

func ProcessText(docName, text string) {
	chunks := Chunk(text)

	for _, chunk := range chunks {
		id, err := SaveChunk(docName, chunk)
		if err != nil {
			log.Fatalf("Failed to save chunk: %v", err)
		}

		embedding := GenerateEmbedding(chunk)
		AddVector(id, embedding)
	}
}

func AnswerQuery(query string) string {
	embedding := GenerateEmbedding(query)
	topIDs := SearchVector(embedding, 3)

	var contexts []string
	for _, id := range topIDs {
		text, err := GetChunkByID(id)
		if err != nil {
			log.Printf("Failed to get chunk for id %d: %v", id, err)
			continue
		}
		contexts = append(contexts, text)
	}

	return GenerateAnswer(query, contexts)
}


