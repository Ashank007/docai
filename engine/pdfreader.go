package engine

import (
	"bytes"
	"io"
	"log"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(path string) string {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Fatal("Failed to open PDF:", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		log.Fatal("Failed to extract text:", err)
	}
	_, err = io.Copy(&buf, b)
	if err != nil {
		log.Fatal("Failed to copy text:", err)
	}
	return buf.String()
}


