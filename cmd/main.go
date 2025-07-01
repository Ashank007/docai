package main

import (
	"bufio" // For reading user input
	"fmt"
	"log"
	"os"    // For os.Stdin
	"path/filepath" // For getting file extension
	"strings"

	"github.com/Ashank007/docai/chain"
	"github.com/Ashank007/docai/chunker"
	"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/generator"
	"github.com/Ashank007/docai/reader" // Make sure this is imported
	"github.com/Ashank007/docai/retriever"
	"github.com/Ashank007/docai/store"
)

func main() {
	// STEP 1: Init all components
	ch := chunker.NewSentenceChunker(200)
	embed := embedder.NewOllama("nomic-embed-text", "http://localhost:11434/api/embeddings")
	gen := generator.NewOllama("llama3.1", "http://localhost:11434/api/generate")

	// Initialize all readers
	pdfReader := reader.NewPDFReader()
	textReader := reader.NewTextReader()
	docxReader := reader.NewDocxReader() // New DocxReader

	meta := store.NewSQLiteStore()
	if err := meta.Init("test.db"); err != nil {
		log.Fatal("❌ SQLite MetadataStore init failed:", err)
	}
	defer meta.Close()

	vector := store.NewMemoryVectorStore()

	retr := retriever.NewCosineRetriever(vector, meta, embed.Embed)

	// Build the chains explicitly for clarity
	// embedChain is now *chain.EmbedChain, not chain.Chain
	actualEmbedChain := &chain.EmbedChain{
		DocName:   "", // This will be set dynamically per document
		Chunker:   ch,
		EmbedFunc: embed.Embed,
		MetaStore: meta,
		VectorDB:  vector,
	}

	actualQueryChain := &chain.QueryChain{
		EmbedFunc: embed.Embed,
		Retriever: retr,
		Generator: gen.Generate,
	}

	// Use the ChainBuilder
	chainBuilder := chain.NewChainBuilder().
		WithEmbedChain(actualEmbedChain).
		WithQueryChain(actualQueryChain)

	// Retrieve the chains from the builder
	// Note: BuildEmbed returns *chain.EmbedChain (concrete type)
	// BuildQuery returns chain.Chain (interface type)
	embedChain := chainBuilder.BuildEmbed()
	queryChain := chainBuilder.BuildQuery()

	// STEP 2: Read and Process Files (PDF, DOCX, TXT)
	// Example: Process multiple documents of different types
	// Make sure these files exist in your testdata/ directory
	documentsToProcess := map[string]string{
		"sample_pdf":  "./testdata/sample.pdf",
		"another_doc": "./testdata/another_document.docx", // Replace with your actual docx file
		"plain_text":  "./testdata/notes.txt",             // Replace with your actual text file
	}

	for docName, filePath := range documentsToProcess {
		fmt.Printf("\nProcessing document: %s (%s)\n", docName, filePath)

		var currentReader reader.Reader
		fileExtension := strings.ToLower(filepath.Ext(filePath))

		switch fileExtension {
		case ".pdf":
			currentReader = pdfReader
		case ".txt":
			currentReader = textReader
		case ".docx":
			currentReader = docxReader
		default:
			log.Printf("⚠️ Unsupported file type for %s: %s. Skipping.\n", filePath, fileExtension)
			continue
		}

		text, err := currentReader.Extract(filePath)
		if err != nil {
			log.Printf("⚠️ File extract failed for %s: %v. Skipping this document.\n", docName, err)
			continue
		}

		embedChain.DocName = docName
		_, err = embedChain.Run(text) // Call Run method on *EmbedChain
		if err != nil {
			log.Fatalf("❌ EmbedChain failed for %s: %v", docName, err)
		}
		fmt.Printf("✅ Document '%s' processed and embedded successfully.\n", docName)
	}

	// STEP 3: Ask a question with optional document filtering
	readerInput := bufio.NewReader(os.Stdin) // Renamed to avoid conflict with `reader` package

	for {
		fmt.Print("\n❓ Enter your query (or type 'exit' to quit): ")
		query, _ := readerInput.ReadString('\n') // Use readerInput
		query = strings.TrimSpace(query)

		if query == "exit" {
			fmt.Println("Exiting. Goodbye!")
			break
		}

		fmt.Print("📄 Enter document name to filter (leave empty for all documents): ")
		docNameFilter, _ := readerInput.ReadString('\n') // Use readerInput
		docNameFilter = strings.TrimSpace(docNameFilter)

		fmt.Printf("\nSearching for: '%s' in document: '%s' (empty means all)\n", query, docNameFilter)

		// Call QueryChain.Run which now correctly expects two arguments
		answer, err := queryChain.Run(query, docNameFilter)
		if err != nil {
			log.Fatalf("❌ QueryChain failed: %v", err)
		}

		// STEP 4: Show the result
		fmt.Println("\n🧠 Final Answer:\n-----------------\n" + answer)
	}
}

