package chain

type Chain interface {
	Run(input string) (string, error)
}
