package client

import "github.com/oanatmaria/ethblkcn-observer/storage"

type Client interface {
	GetLatestBlockNumber() (int, error)
	GetBlockByNumber(blockNum int) (*Block, error)
}

type Block struct {
	Number       int
	Transactions []storage.Transaction
}
