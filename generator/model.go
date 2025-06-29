package generator

import (
	"bufio"
	"bytes"
	"encoding/json"
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

// NewOllama returns a Generator using the Ollama API
func NewOllama(model, url string) *OllamaGenerator {
	return &OllamaGenerator{
		Model: model,
		URL:   url,
	}
}

// Generate sends the prompt to Ollama and prints the response as it streams
func (g *OllamaGenerator) Generate(query string, contexts []string) (string, error) {
	prompt := g.constructPrompt(query, contexts)
	reqBody := genRequest{
		Model:  g.Model,
		Prompt: prompt,
		Stream: true, // ✅ Enable streaming
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

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		var chunk genResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue // Ignore malformed lines
		}
		fmt.Print(chunk.Response)                 // 🔥 Live terminal output
		fullResponse.WriteString(chunk.Response) // Collect full response
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("stream read error: %w", err)
	}

	return strings.TrimSpace(fullResponse.String()), nil
}

// constructPrompt builds the final prompt from the context and question
func (g *OllamaGenerator) constructPrompt(query string, contexts []string) string {
	return fmt.Sprintf(`You are a helpful assistant AI. Use the following context to answer the question.

Context:
%s

Question: %s

Answer briefly.`, strings.Join(contexts, "\n"), query)
}

// Name returns the name of this generator
func (g *OllamaGenerator) Name() string {
	return "ollama"
}


