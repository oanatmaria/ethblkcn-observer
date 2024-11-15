package parser

import "github.com/oanatmaria/ethblkcn-observer/storage"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []storage.Transaction

	ProcessNewBlocks()
}
