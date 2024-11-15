package parser

import (
	"log"
	"sync"

	"github.com/oanatmaria/ethblkcn-observer/client"
	"github.com/oanatmaria/ethblkcn-observer/storage"
)

type EthParser struct {
	storage      storage.Storage
	client       client.Client
	currentBlock int
	mu           sync.Mutex
}

func NewEthParser(storage storage.Storage, client client.Client) Parser {
	latestBlock, err := client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Error fetching latest block: %v\n", err)
		return nil
	}

	return &EthParser{
		storage:      storage,
		client:       client,
		currentBlock: latestBlock,
	}
}

func (p *EthParser) GetCurrentBlock() int {
	return p.currentBlock
}

func (p *EthParser) Subscribe(address string) bool {
	return p.storage.AddObservedAddress(address)
}

func (p *EthParser) GetTransactions(address string) []storage.Transaction {
	return p.storage.GetTransactions(address)
}

func (p *EthParser) ProcessNewBlocks() {
	latestBlock, err := p.client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Error fetching latest block: %v\n", err)
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for blockNum := p.currentBlock + 1; blockNum <= latestBlock; blockNum++ {
		block, err := p.client.GetBlockByNumber(blockNum)
		if err != nil {
			log.Printf("Error fetching block %d: %v\n", blockNum, err)
			continue
		}

		for _, tx := range block.Transactions {
			if p.storage.IsObservedAddress(tx.From) || p.storage.IsObservedAddress(tx.To) {
				p.storage.AddTransaction(tx)
			}
		}
	}

	p.currentBlock = latestBlock
}
