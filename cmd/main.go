package main

import (
	"fmt"
	"log"

	"github.com/Ashank007/docai/chunker"
	"github.com/Ashank007/docai/embedder"
	"github.com/Ashank007/docai/generator"
	"github.com/Ashank007/docai/reader"
	"github.com/Ashank007/docai/retriever"
	"github.com/Ashank007/docai/store"
	//"github.com/Ashank007/docai/types"
)

func main() {
	// STEP 1: Init all components
	ch := chunker.NewSentenceChunker(200)
	embed := embedder.NewOllama("nomic-embed-text", "http://localhost:11434/api/embeddings")
	gen := generator.NewOllama("llama3", "http://localhost:11434/api/generate")
	pdf := reader.NewPDFReader()

	meta := store.NewSQLiteStore()
	if err := meta.Init("test.db"); err != nil {
		log.Fatal("❌ SQLite init failed:", err)
	}
	defer meta.Close()

	vector := store.NewMemoryVectorStore()
	retr := retriever.NewCosineRetriever(vector, meta, embed.Embed)

	// STEP 2: Read PDF
	text, err := pdf.Extract("./testdata/sample.pdf")
	if err != nil {
		log.Fatal("❌ PDF extract failed:", err)
	}

	// STEP 3: Chunk the PDF text
	chunks, err := ch.Chunk(text)
	if err != nil {
		log.Fatal("❌ Chunking failed:", err)
	}

	// STEP 4: Embed and Store each chunk
	docName := "sample"
	for i, chunk := range chunks {
		chunk.Source = docName
		chunk.Position = i

		id, err := meta.SaveChunk(docName, chunk)
		if err != nil {
			log.Fatalf("❌ Failed to save chunk %d: %v", i, err)
		}

		vec, err := embed.Embed(chunk.Text)
		if err != nil {
			log.Fatalf("❌ Failed to embed chunk %d: %v", i, err)
		}

		if err := vector.AddVector(id, vec); err != nil {
			log.Fatalf("❌ Failed to store vector for chunk %d: %v", i, err)
		}
	}

	// STEP 5: Ask a question
	query := "What is the main topic of this document?"
	retrievedChunks, err := retr.Retrieve(query, 3)
	if err != nil {
		log.Fatal("❌ Retrieval failed:", err)
	}

	var contextTexts []string
	for _, c := range retrievedChunks {
		contextTexts = append(contextTexts, c.Chunk.Text)
	}

	answer, err := gen.Generate(query, contextTexts)
	if err != nil {
		log.Fatal("❌ Generation failed:", err)
	}

	// STEP 6: Show the result
	fmt.Println("\n🧠 Final Answer:\n-----------------\n" + answer)
}
