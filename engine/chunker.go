package engine

import "strings"

func Chunk(text string) []string {
	sentences := strings.Split(text, ".")
	var chunks []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if len(s) > 30 {
			chunks = append(chunks, s+".")
		}
	}
	return chunks
}


