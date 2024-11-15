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

//go:generate mockgen -destination=mock_storage.go -package=storage github.com/oanatmaria/ethblkcn-observer/storage Storage
type Storage interface {
	AddObservedAddress(address string) bool
	GetTransactions(address string) []Transaction
	AddTransactions(txs ...Transaction)
	GetCurrentBlock() int
	UpdateCurrentBlock(block int)
}
