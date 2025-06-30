package chain

type Chain interface {
  Run(query string, docNameFilter string) (string, error)
}
