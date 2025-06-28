package engine

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const ollamaEmbedURL = "http://localhost:11434/api/embeddings"

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GenerateEmbedding(text string) []float32 {
	req := embedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	payload, _ := json.Marshal(req)
	resp, err := http.Post(ollamaEmbedURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Embedding request failed:", err)
	}
	defer resp.Body.Close()

	var data embedResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode embedding response:", err)
	}

	return data.Embedding
}


