package reader

import (
	"bytes"
	"fmt"
	"io"
	//"os"

	"github.com/ledongthuc/pdf"
)

type PDFReader struct{}

func NewPDFReader() *PDFReader {
	return &PDFReader{}
}

func (p *PDFReader) Extract(path string) (string, error) {
	file, reader, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	text, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("failed to extract plain text: %w", err)
	}
	_, err = io.Copy(&buf, text)
	if err != nil {
		return "", fmt.Errorf("failed to copy text from buffer: %w", err)
	}

	return buf.String(), nil
}
