package engine

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const ollamaGenURL = "http://localhost:11434/api/generate"

type genRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type genResponse struct {
	Response string `json:"response"`
}

func GenerateAnswer(query string, contexts []string) string {
	fullPrompt := `You are a helpful assistant AI. Use the following context to answer the question.

Context:
` + strings.Join(contexts, "\n") + `

Question: ` + query + `

Answer briefly.`

	req := genRequest{
		Model:  "llama3.1",
		Prompt: fullPrompt,
		Stream: false,
	}

	payload, _ := json.Marshal(req)
	resp, err := http.Post(ollamaGenURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Generation request failed:", err)
	}
	defer resp.Body.Close()

	var data genResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode generation response:", err)
	}

	return strings.TrimSpace(data.Response)
}


