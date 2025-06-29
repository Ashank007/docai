package generator

// Generator defines the interface for any text generation model
type Generator interface {
	Generate(query string, contexts []string) (string, error)
	Name() string
}
