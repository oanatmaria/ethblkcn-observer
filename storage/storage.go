package storage

type Transaction struct {
	Hash      string
	From      string
	To        string
	Value     string
	BlockHash string
	BlockNum  int
	Type      string
}

type Storage interface {
	AddObservedAddress(address string) bool
	IsObservedAddress(address string) bool
	GetTransactions(address string) []Transaction
	AddTransaction(tx Transaction)
}
