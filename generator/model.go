package generator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type OllamaGenerator struct {
	Model string
	URL   string // e.g., http://localhost:11434/api/generate
}

type genRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type genResponse struct {
	Response string `json:"response"`
}

// NewOllama returns a Generator using the Ollama APIa
func NewOllama(model, url string) *OllamaGenerator {
	return &OllamaGenerator{
		Model: model,
		URL:   url,
	}
}

// Generate constructs the prompt and gets completion from Ollama
func (g *OllamaGenerator) Generate(query string, contexts []string) (string, error) {
	prompt := g.constructPrompt(query, contexts)

	reqBody := genRequest{
		Model:  g.Model,
		Prompt: prompt,
		Stream: false,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(g.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("generation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result genResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("failed to decode generation response")
	}

	return strings.TrimSpace(result.Response), nil
}

func (g *OllamaGenerator) constructPrompt(query string, contexts []string) string {
	return fmt.Sprintf(`You are a helpful assistant AI. Use the following context to answer the question.

Context:
%s

Question: %s

Answer briefly.`, strings.Join(contexts, "\n"), query)
}

func (g *OllamaGenerator) Name() string {
	return "ollama"
}
