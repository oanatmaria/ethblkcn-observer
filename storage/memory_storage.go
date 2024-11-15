package storage

import (
	"sync"
)

type MemoryStorage struct {
	observedAddresses map[string]struct{}
	transactions      map[string][]Transaction
	currentBlock      int
	mu                sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		observedAddresses: make(map[string]struct{}),
		transactions:      make(map[string][]Transaction),
		currentBlock:      0,
	}
}

func (s *MemoryStorage) AddObservedAddress(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.observedAddresses[address]; exists {
		return false
	}
	s.observedAddresses[address] = struct{}{}
	return true
}

func (s *MemoryStorage) GetTransactions(address string) []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transactions[address]
}

func (s *MemoryStorage) AddTransactions(txs ...Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, tx := range txs {
		if _, exists := s.observedAddresses[tx.From]; exists {
			s.transactions[tx.From] = append(s.transactions[tx.From], tx)
		}

		if tx.To != "" {
			if _, exists := s.observedAddresses[tx.To]; exists {
				s.transactions[tx.To] = append(s.transactions[tx.To], tx)
			}
		}
	}
}

func (s *MemoryStorage) GetCurrentBlock() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentBlock
}

func (s *MemoryStorage) UpdateCurrentBlock(block int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentBlock = block
}
