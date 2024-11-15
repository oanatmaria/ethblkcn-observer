package parser

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/oanatmaria/ethblkcn-observer/client"
	"github.com/oanatmaria/ethblkcn-observer/storage"
)

type EthParser struct {
	storage storage.Storage
	client  client.Client
}

func NewEthParser(storage storage.Storage, client client.Client) (Parser, error) {
	latestBlock, err := client.GetLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("error fetching latest block: %v", err)
	}

	storage.UpdateCurrentBlock(latestBlock)

	return &EthParser{
		storage: storage,
		client:  client,
	}, nil
}

func (p *EthParser) GetCurrentBlock() int {
	return p.storage.GetCurrentBlock()
}

func (p *EthParser) Subscribe(address string) bool {
	return p.storage.AddObservedAddress(address)
}

func (p *EthParser) GetTransactions(address string) []storage.Transaction {
	return p.storage.GetTransactions(address)
}
func (p *EthParser) ProcessNewBlocks(ctx context.Context) {
	latestBlock, err := p.client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Error fetching latest block: %v\n", err)
		return
	}

	currentBlock := p.storage.GetCurrentBlock()
	if currentBlock >= latestBlock {
		return
	}

	log.Println(currentBlock)

	blockChan := make(chan int)
	var wg sync.WaitGroup

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case blockNum, ok := <-blockChan:
					if !ok {
						return
					}
					block, err := p.client.GetBlockByNumber(blockNum)
					if err != nil {
						log.Printf("Error fetching block %d: %v\n", blockNum, err)
						continue
					}
					p.storage.AddTransactions(block.Transactions...)
				}
			}

		}()
	}

	for blockNum := currentBlock + 1; blockNum <= latestBlock; blockNum++ {
		select {
		case <-ctx.Done():
			close(blockChan)
			return
		case blockChan <- blockNum:
		}
	}
	close(blockChan)

	wg.Wait()

	p.storage.UpdateCurrentBlock(latestBlock)
}
