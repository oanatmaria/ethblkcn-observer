package storage

import (
	"sync"
)

type MemoryStorage struct {
	observedAddresses map[string]struct{}
	transactions      map[string][]Transaction
	mu                sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		observedAddresses: make(map[string]struct{}),
		transactions:      make(map[string][]Transaction),
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

func (s *MemoryStorage) IsObservedAddress(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.observedAddresses[address]
	return exists
}

func (s *MemoryStorage) GetTransactions(address string) []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transactions[address]
}

func (s *MemoryStorage) AddTransaction(tx Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[tx.From] = append(s.transactions[tx.From], tx)
	if tx.To != "" {
		s.transactions[tx.To] = append(s.transactions[tx.To], tx)
	}
}
