package client

import "github.com/oanatmaria/ethblkcn-observer/storage"

//go:generate mockgen -destination=mock_client.go -package=client github.com/oanatmaria/ethblkcn-observer/client Client

type Client interface {
	GetLatestBlockNumber() (int, error)
	GetBlockByNumber(blockNum int) (Block, error)
}

type Block struct {
	Number       int
	Transactions []storage.Transaction
}
