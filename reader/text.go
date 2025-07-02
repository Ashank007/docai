package reader

import (
	"fmt"
	"os" // Use ioutil for simplicity, though os.ReadFile is preferred in newer Go
)

// TextReader implements the Reader interface for plain text files (.txt).
type TextReader struct{}

// NewTextReader creates a new TextReader.
func NewTextReader() *TextReader {
	return &TextReader{}
}

// Extract reads the content of a plain text file.
func (r *TextReader) Extract(filePath string) (string, error) {
	content, err := os.ReadFile(filePath) // Use os.ReadFile for Go 1.16+
	if err != nil {
		return "", fmt.Errorf("failed to read text file %s: %w", filePath, err)
	}
	return string(content), nil
}
