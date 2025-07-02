
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
	"github.com/Ashank007/docai/chain"
	"github.com/Ashank007/docai/chunker"
	"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/generator"
	"github.com/Ashank007/docai/reader"
	"github.com/Ashank007/docai/retriever"
	"github.com/Ashank007/docai/store"
	"github.com/Ashank007/docai/summarizer" // New import for the summarizer package
)

func main() {
	// STEP 1: Init all components
	// Chunkers, Embedders, Generators, Readers (No Change)
	ch := chunker.NewSentenceChunker(200)
	embed := embedder.NewOllama("nomic-embed-text", "http://localhost:11434/api/embeddings")
	gen := generator.NewOllama("llama3.1", "http://localhost:11434/api/generate")

	pdfReader := reader.NewPDFReader()
	textReader := reader.NewTextReader()
	docxReader := reader.NewDocxReader() // Ensure this uses your robust manual parser

	meta := store.NewSQLiteStore()
	if err := meta.Init("test.db"); err != nil {
		log.Fatal("❌ SQLite MetadataStore init failed:", err)
	}
	defer meta.Close()

	vector := store.NewMemoryVectorStore()

	retr := retriever.NewCosineRetriever(vector, meta, embed.Embed)

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

	chainBuilder := chain.NewChainBuilder().
		WithEmbedChain(actualEmbedChain).
		WithQueryChain(actualQueryChain)

	embedChain := chainBuilder.BuildEmbed()
	queryChain := chainBuilder.BuildQuery()

	// Initialize the new Summarizer library component
	docSummarizer := summarizer.NewSummarizer(ch, gen, pdfReader, textReader, docxReader)

	// ---

	// STEP 2: Read and Process Files (Initial setup documents - same as before)
	// These documents are processed and embedded into the vector store for querying.
	// The summarizer works independently on a given file path.
	documentsToProcess := map[string]string{
		"sample_pdf":        "./testdata/sample.pdf",
		"another_doc":       "./testdata/another_document.docx",
		"notes_data":        "./testdata/notes.txt",
	}

	for docName, filePath := range documentsToProcess {
		fmt.Printf("\nProcessing document for query indexing: %s (%s)\n", docName, filePath)
		
		var currentReader reader.Reader // We need to determine reader type here too for initial embedding
		fileExtension := strings.ToLower(filepath.Ext(filePath))

		switch fileExtension {
		case ".pdf":
			currentReader = pdfReader
		case ".txt":
			currentReader = textReader
		case ".docx":
			currentReader = docxReader
		default:
			log.Printf("⚠️ Unsupported file type for %s: %s. Skipping for query indexing.\n", filePath, fileExtension)
			continue
		}

		text, err := currentReader.Extract(filePath)
		if err != nil {
			log.Printf("⚠️ File extract failed for %s: %v. Skipping for query indexing.\n", docName, err)
			continue
		}

		embedChain.DocName = docName
		_, err = embedChain.Run(text)
		if err != nil {
			log.Fatalf("❌ EmbedChain failed for %s: %v", docName, err)
		}
		fmt.Printf("✅ Document '%s' indexed for querying successfully.\n", docName)
	}

	// ---

	// STEP 3: User Interaction Loop (Modified for Summarizer)
	readerInput := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nChoose an action:\n")
		fmt.Print("1. Query documents\n")
		fmt.Print("2. Summarize a document\n")
		fmt.Print("   (Type 'exit' to quit)\n")
		fmt.Print("Enter choice (1 or 2): ")
		choice, _ := readerInput.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "exit" {
			fmt.Println("Exiting. Goodbye!")
			break
		}

		switch choice {
		case "1": // Query existing documents
			fmt.Print("\n❓ Enter your query: ")
			query, _ := readerInput.ReadString('\n')
			query = strings.TrimSpace(query)

			fmt.Print("📄 Enter document name to filter (leave empty for all documents): ")
			docNameFilter, _ := readerInput.ReadString('\n')
			docNameFilter = strings.TrimSpace(docNameFilter)

			fmt.Printf("\nSearching for: '%s' in document: '%s' (empty means all)\n", query, docNameFilter)

			answer, err := queryChain.Run(query, docNameFilter)
			if err != nil {
				log.Printf("❌ QueryChain failed: %v", err) // Use log.Printf instead of log.Fatalf here for graceful error handling
			} else {
				fmt.Println("\n🧠 Final Answer:\n-----------------\n" + answer)
			}

		case "2": // Summarize a specific document using the new summarizer library
			fmt.Print("\n📂 Enter the full path to the document you want to summarize (e.g., './testdata/notes.txt'): ")
			filePathToSummarize, _ := readerInput.ReadString('\n')
			filePathToSummarize = strings.TrimSpace(filePathToSummarize)

			if filePathToSummarize == "" {
				fmt.Println("🚫 File path cannot be empty for summarization. Please try again.")
				continue
			}

			fmt.Printf("\nSummarizing document: '%s'\n", filePathToSummarize)

			summary, err := docSummarizer.SummarizeDocument(filePathToSummarize)
			if err != nil {
				log.Printf("❌ Summarization failed for '%s': %v", filePathToSummarize, err)
			} else {
				fmt.Println("\n📝 Summary:\n------------\n" + summary)
			}

		default:
			fmt.Println("Invalid choice. Please enter '1' or '2'.")
		}
	}
}

