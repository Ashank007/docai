package reader

type Reader interface {
	Extract(path string) (string, error)
}
