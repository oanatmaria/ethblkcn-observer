package parser

import (
	"context"

	"github.com/oanatmaria/ethblkcn-observer/storage"
)

//go:generate mockgen -destination=mock_perser.go -package=parser github.com/oanatmaria/ethblkcn-observer/parser Parser

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []storage.Transaction

	ProcessNewBlocks(ctx context.Context)
}
