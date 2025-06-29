package embedder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type OllamaEmbedder struct {
	Model string
	URL   string // e.g., http://localhost:11434/api/embeddings
}

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// NewOllama returns an Embedder for Ollama's embedding API
func NewOllama(model, url string) *OllamaEmbedder {
	return &OllamaEmbedder{
		Model: model,
		URL:   url,
	}
}

func (o *OllamaEmbedder) Embed(text string) ([]float32, error) {
	reqBody := embedRequest{
		Model:  o.Model,
		Prompt: text,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embed request: %w", err)
	}

	resp, err := http.Post(o.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding request failed with status: %s", resp.Status)
	}

	var result embedResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.New("failed to decode embedding response")
	}

	return result.Embedding, nil
}

func (o *OllamaEmbedder) Name() string {
	return "ollama-embedder"
}
