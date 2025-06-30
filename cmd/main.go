package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Ashank007/docai/chain"
	"github.com/Ashank007/docai/chunker"
	"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/generator"
	"github.com/Ashank007/docai/reader"
	"github.com/Ashank007/docai/retriever"
	"github.com/Ashank007/docai/store"
)

func main() {
	// STEP 1: Init all components
	ch := chunker.NewSentenceChunker(200)
	embed := embedder.NewOllama("nomic-embed-text", "http://localhost:11434/api/embeddings")
	gen := generator.NewOllama("llama3", "http://localhost:11434/api/generate")
	pdf := reader.NewPDFReader()

	meta := store.NewSQLiteStore()
	if err := meta.Init("test.db"); err != nil {
		log.Fatal("❌ SQLite MetadataStore init failed:", err)
	}
	defer meta.Close()

	vector := store.NewMemoryVectorStore()

	retr := retriever.NewCosineRetriever(vector, meta, embed.Embed)

	// Build the chains explicitly for clarity
	actualEmbedChain := &chain.EmbedChain{
		DocName:   "",
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
	// CORRECTED LINE: embedChain now receives *chain.EmbedChain
	embedChain := chainBuilder.BuildEmbed() // <--- CHANGE IS HERE
	queryChain := chainBuilder.BuildQuery()

	// STEP 2: Read and Process PDF(s)
	documentsToProcess := map[string]string{
		"sample_pdf":  "./testdata/sample.pdf",
		"another_doc": "./testdata/another_document.pdf",
	}

	for docName, filePath := range documentsToProcess {
		fmt.Printf("\nProcessing document: %s (%s)\n", docName, filePath)
		text, err := pdf.Extract(filePath)
		if err != nil {
			log.Printf("⚠️ PDF extract failed for %s: %v. Skipping this document.\n", docName, err)
			continue
		}

		embedChain.DocName = docName
		_, err = embedChain.Run(text)
		if err != nil {
			log.Fatalf("❌ EmbedChain failed for %s: %v", docName, err)
		}
		fmt.Printf("✅ Document '%s' processed and embedded successfully.\n", docName)
	}

	// STEP 3: Ask a question with optional document filtering
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n❓ Enter your query (or type 'exit' to quit): ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)

		if query == "exit" {
			fmt.Println("Exiting. Goodbye!")
			break
		}

		fmt.Print("📄 Enter document name to filter (leave empty for all documents): ")
		docNameFilter, _ := reader.ReadString('\n')
		docNameFilter = strings.TrimSpace(docNameFilter)

		fmt.Printf("\nSearching for: '%s' in document: '%s' (empty means all)\n", query, docNameFilter)

		answer, err := queryChain.Run(query, docNameFilter)
		if err != nil {
			log.Fatalf("❌ QueryChain failed: %v", err)
		}

		// STEP 4: Show the result
		fmt.Println("\n🧠 Final Answer:\n-----------------\n" + answer)
	}
}

