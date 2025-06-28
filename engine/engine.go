package engine

func ProcessTextWrapper(docName, text string) {
	ProcessText(docName, text)
}

func AnswerQueryWrapper(query string) string {
	return AnswerQuery(query)
}

