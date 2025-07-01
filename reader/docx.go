package reader

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DocxReader implements the Reader interface for .docx files.
type DocxReader struct{}

// NewDocxReader creates a new DocxReader.
func NewDocxReader() *DocxReader {
	return &DocxReader{}
}

// Extract reads the content of a .docx file and returns it as plain text.
func (r *DocxReader) Extract(filePath string) (string, error) {
	// DOCX files are essentially ZIP archives.
	// The main text content is typically found in 'word/document.xml'.
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open docx file %s as zip: %w", filePath, err)
	}
	defer zipReader.Close()

	var docText strings.Builder
	foundDocumentXML := false

	// Iterate through the files within the zip archive
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			foundDocumentXML = true
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open 'word/document.xml' in docx: %w", err)
			}
			defer rc.Close()

			// XML parser to extract text
			decoder := xml.NewDecoder(rc)
			for {
				token, err := decoder.Token()
				if err == io.EOF {
					break // End of XML file
				}
				if err != nil {
					return "", fmt.Errorf("error parsing 'word/document.xml' in docx: %w", err)
				}

				switch se := token.(type) {
				case xml.StartElement:
					// Look for <w:t> (text run) tags for actual text
					if se.Name.Local == "t" {
						var textContent string
						if err := decoder.DecodeElement(&textContent, &se); err != nil {
							return "", fmt.Errorf("error decoding text element: %w", err)
						}
						docText.WriteString(textContent)
					} else if se.Name.Local == "p" {
						// Add a newline for paragraph breaks to maintain readability
						docText.WriteString("\n")
					}
				}
			}
			break // We found and processed document.xml, no need to check other files
		}
	}

	if !foundDocumentXML {
		return "", fmt.Errorf("'word/document.xml' not found in docx file %s", filePath)
	}

	// Return the extracted text, trimming leading/trailing whitespace
	return strings.TrimSpace(docText.String()), nil
}

