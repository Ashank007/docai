
package summarizer

import (
	"fmt"
	"strings"
  "path/filepath"
	"github.com/Ashank007/docai/chunker"
	"github.com/Ashank007/docai/generator"
	"github.com/Ashank007/docai/reader"
)

// Summarizer holds the necessary components for document summarization.
type Summarizer struct {
	Chunker    chunker.Chunker
	Generator  generator.Generator
	PDFReader  reader.Reader
	TextReader reader.Reader
	DocxReader reader.Reader
}

// NewSummarizer creates and returns a new Summarizer instance.
func NewSummarizer(
	ch chunker.Chunker,
	gen generator.Generator,
	pdfR reader.Reader,
	textR reader.Reader,
	docxR reader.Reader,
) *Summarizer {
	return &Summarizer{
		Chunker:    ch,
		Generator:  gen,
		PDFReader:  pdfR,
		TextReader: textR,
		DocxReader: docxR,
	}
}

// SummarizeDocument reads a document from the given filePath,
// chunks its content, and uses the LLM to generate a summary.
func (s *Summarizer) SummarizeDocument(filePath string) (string, error) {
	var currentReader reader.Reader
	fileExtension := strings.ToLower(filepath.Ext(filePath))

	switch fileExtension {
	case ".pdf":
		currentReader = s.PDFReader
	case ".txt":
		currentReader = s.TextReader
	case ".docx":
		currentReader = s.DocxReader
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileExtension)
	}

	// 1. Read the document content
	fullText, err := currentReader.Extract(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from %s: %w", filePath, err)
	}

	// If the document is empty after extraction, return an appropriate message.
	if strings.TrimSpace(fullText) == "" {
		return "The document is empty or contains no extractable text.", nil
	}

	// 2. Chunk the document
	// FIX 1: Handle the error returned by s.Chunker.Chunk
	chunks, err := s.Chunker.Chunk(fullText) // Now correctly handles both return values
	if err != nil {
		return "", fmt.Errorf("failed to chunk document: %w", err)
	}

	if len(chunks) == 0 {
		return "No meaningful chunks could be created from the document for summarization.", nil
	}

	// 3. Prepare content for LLM
	// Convert types.Chunk slice to a []string slice for the Generator
	var chunkStrings []string
	for _, chunk := range chunks {
		chunkStrings = append(chunkStrings, chunk.Text)
	}

	// 4. Create the summarization prompt
	// The prompt structure is crucial for good summaries.
	// We instruct the LLM to summarize the provided text.
	// Note: We are now passing the context as a separate slice to the Generator.
	// The Generator implementation (e.g., generator/ollama.go) will combine these.
	prompt := "Please provide a concise and comprehensive summary of the following document. Focus on the main ideas and key information."

	// 5. Send to LLM for summarization
	// FIX 2: Pass both the prompt (query) and the chunkStrings (context) to s.Generator.Generate
	summary, err := s.Generator.Generate(prompt, chunkStrings) // Now passing two arguments
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return summary, nil
}

