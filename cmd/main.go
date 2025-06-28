package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Ashank007/docai/engine"
)

func main() {
	engine.InitDB()
	engine.InitVectorStore()

	fmt.Print("Enter PDF path: ")
	var path string
	fmt.Scanln(&path)

	text := engine.ExtractTextFromPDF(path)
	engine.ProcessText(path, text)

	fmt.Println("PDF processed. You can ask questions now (type 'exit' to quit).")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nQuestion: ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)
		if query == "exit" {
			break
		}

		answer := engine.AnswerQuery(query)
		fmt.Println("Answer:", answer)
	}
}


