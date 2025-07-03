
![Go Build Status](https://img.shields.io/badge/Go-1.20%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Ollama Compatible](https://img.shields.io/badge/Ollama-Compatible-brightgreen)
![SQLite Backed](https://img.shields.io/badge/SQLite-Backed-orange)
![Supported Document Types](https://img.shields.io/badge/Docs-PDF%2C%20DOCX%2C%20TXT-red)
# DocAI: Intelligent Document Processing & Querying Toolkit
## 🚀 Overview

DocAI is a versatile Go-based toolkit designed for intelligent document processing, enabling users to effortlessly extract information, query document content, and generate concise summaries using local Large Language Models (LLMs) via Ollama. It combines robust file parsing, efficient text chunking, vector embedding, and retrieval-augmented generation (RAG) to provide a powerful command-line interface for interacting with your document collection.

The project is structured to be modular, allowing core functionalities like document summarization to be easily integrated into other Go applications as a standalone library.

## ✨ Features

* **Multi-Format Document Parsing**: Supports `.pdf`, `.docx`, and `.txt` file types for comprehensive data ingestion.
* **Intelligent Text Chunking**: Breaks down large documents into manageable, semantically relevant chunks for efficient LLM processing.
* **Local LLM Integration (Ollama)**: Leverages local Ollama installations for privacy-preserving and cost-effective text embeddings (`nomic-embed-text`) and response generation (`llama3.1`).
* **Vector Database & Metadata Storage**: Utilizes an in-memory vector store for semantic search and SQLite for document metadata management.
* **Retrieval-Augmented Generation (RAG)**: Enhances LLM responses by retrieving relevant document snippets based on user queries, providing accurate and contextual answers.
* **Document Querying**: Ask questions about your processed documents and get AI-generated answers based on the content.
* **Document Summarization (Library & CLI)**: Get concise summaries of entire documents. This feature is also exposed as a reusable Go library (`pkg/summarizer`).
* **Modular Design**: Cleanly separated concerns (readers, chunkers, embedders, generators, chains, stores) for maintainability and extensibility.

## 📦 Installation & Setup

### Prerequisites

1.  **Go Language**: Ensure you have Go 1.20 or higher installed.
    * [Download Go](https://golang.org/doc/install)
2.  **Ollama**: Install Ollama to run local LLMs.
    * [Download Ollama](https://ollama.com/download)
3.  **Ollama Models**: Pull the required models using Ollama CLI:
    ```bash
    ollama pull nomic-embed-text # Can use any
    ollama pull llama3.1 # Or your preferred LLM like llama3, mistral, etc.
    ```

### Project Setup

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/Ashank007/docai.git
    cd docai
    ```
2.  **Prepare Test Data**:
    * Create a testdata directory in the project root:
        ```bash
        mkdir testdata
        ```
    * Place your sample documents (`.pdf`, `.docx`, `.txt`) inside the `testdata` directory. Ensure they are named as expected in `cmd/main.go` (e.g., `sample.pdf`, `another_document.docx`, `notes.txt`). **Crucially, ensure your `.docx` files are valid and not corrupted.**

3.  **Install Go Dependencies**:
    ```bash
    go mod tidy
    ```

## 🚀 Usage

### Running the CLI Application

Navigate to the project root and run the main application:

```bash
go run cmd/main.go
```

The application will first process and embed the documents specified in cmd/main.go. After successful processing, it will present an interactive menu:

Processing document for query indexing: sample_pdf (./testdata/sample.pdf)
✅ Document 'sample_pdf' indexed for querying successfully.

similar output for other documents

### Choose an action:
- Query documents
- Summarize a document
   (Type 'exit' to quit)
Enter choice (1 or 2):

### Option 1: Query Documents

Enter 1 to query your indexed documents.

Enter choice (1 or 2): 1

❓ Enter your query: What are the key topics in Unit 3 of the notes?

📄 Enter document name to filter (leave empty for all documents): notes_data

Searching for: 'What are the key topics in Unit 3 of the notes?' in document: 'notes_data' (empty means all)

🧠 Final Answer:

Unit 3 focuses on transactions, concurrency control protocols (lock-based, timestamp-based, validation-based), and deadlock handling. It also covers transaction definition in SQL.

### Option 2: Summarize a Document

Enter 2 to summarize a specific document. You will be prompted to enter the full file path (relative to the project root).

Enter choice (1 or 2): 2

📂 Enter the full path to the document you want to summarize (e.g., './testdata/notes.txt'): ./testdata/another_document.docx

Summarizing document: './testdata/another_document.docx'

📝 Summary:

This document provides a summary of the key features of the DocAI toolkit, emphasizing its capabilities in document processing, querying, and summarization using local LLMs. It highlights support for various document formats (PDF, DOCX, TXT), intelligent text chunking, and integration with Ollama for embedding and generation. The toolkit leverages a vector database for semantic search and offers a modular design for reusability, particularly for its document summarization functionality.

## 📚 Library Usage: Document Summarization

The core summarization logic is exposed as a Go package pkg/summarizer, allowing you to integrate document summarization into your own Go applications.

Installation
```
go get https://github.com/Ashank007/docai/summarizer
```
### Example Usage
```
package main

import (
	"fmt"
	"log"

	"https://github.com/Ashank007/docai/chunker"
	"https://github.com/Ashank007/docai/generator"
	"https://github.com/Ashank007/docai/reader"
	"https://github.com/Ashank007/docai/summarizer"
)

func main() {
    // Initialize required components for the summarizer
    // (These dependencies should be initialized once in your application)
    ch := chunker.NewSentenceChunker(200) // Or your preferred chunker implementation
    gen := generator.NewOllama("llama3.1", "http://localhost:11434/api/generate")

    pdfReader := reader.NewPDFReader()
    textReader := reader.NewTextReader()
    docxReader := reader.NewDocxReader()

    // Create a new Summarizer instance
    docSummarizer := summarizer.NewSummarizer(ch, gen, pdfReader, textReader, docxReader)

    // Specify the path to the document you want to summarize
    filePath := "./path/to/your/document.txt" // Replace with actual path
    // filePath := "./path/to/your/report.pdf"
    // filePath := "./path/to/your/notes.docx"

    fmt.Printf("Attempting to summarize: %s\n", filePath)

    // Call the SummarizeDocument method
    summary, err := docSummarizer.SummarizeDocument(filePath)
    if err != nil {
        log.Fatalf("Error summarizing document: %v", err)
    }

    fmt.Println("\n--- Document Summary ---")
    fmt.Println(summary)
}
```

## 🏗️ Project Structure
```
.
├── cmd/
│   └── main.go           # Main application entry point
├── chain/
│   ├── embed.go          # Handles document embedding workflow
│   ├── query.go          # Manages query processing and RAG
│   └── builder.go        # Chain builder for structured setup
├── chunker/
│   ├── chunker.go        # Chunker interface
│   └── sentence.go       # Sentence-based chunking implementation
├── embedder/
│   └── ollama.go         # Ollama API integration for embeddings
├── generator/
│   └── ollama.go         # Ollama API integration for text generation (LLM)
├── reader/
│   ├── reader.go         # Document reader interface
│   ├── pdf.go            # PDF reading implementation
│   ├── text.go           # Plain text reading implementation
│   └── docx.go           # DOCX reading implementation (robust standard lib parsing)
├── retriever/
│   └── cosine.go         # Cosine similarity based document retrieval
├── store/
│   ├── store.go          # Interfaces for metadata and vector stores
│   ├── sqlite.go         # SQLite implementation for metadata
│   └── memory.go         # In-memory implementation for vector store
├── types/
│   └── types.go          # Core data structures (e.g., Chunk, Document)
└── summarizer/
        └── summarizer.go # Standalone library for document summarization
```
## 🤝 Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the MIT License
