package storage

import (
	"reflect"
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	if storage == nil {
		t.Errorf("Expected MemoryStorage instance to not be nil")
	}

	if storage.GetCurrentBlock() != 0 {
		t.Errorf("Expected initial block to be 0, got %d", storage.GetCurrentBlock())
	}
}

func TestAddObservedAddress(t *testing.T) {
	storage := NewMemoryStorage()

	if !storage.AddObservedAddress("address1") {
		t.Errorf("Expected adding new address to return true")
	}

	if storage.AddObservedAddress("address1") {
		t.Errorf("Expected adding duplicate address to return false")
	}
}

func TestGetTransactions(t *testing.T) {
	storage := NewMemoryStorage()

	txs := storage.GetTransactions("address1")
	if len(txs) != 0 {
		t.Errorf("Expected no transactions for unobserved address, got %d", len(txs))
	}

	storage.AddObservedAddress("address1")
	tx1 := Transaction{
		Hash:      "tx1",
		From:      "address1",
		To:        "address2",
		Value:     "100",
		BlockHash: "blockhash1",
		BlockNum:  1,
		Type:      "transfer",
	}
	tx2 := Transaction{
		Hash:      "tx2",
		From:      "address1",
		To:        "address3",
		Value:     "200",
		BlockHash: "blockhash2",
		BlockNum:  2,
		Type:      "transfer",
	}
	storage.AddTransactions(tx1, tx2)

	txs = storage.GetTransactions("address1")
	if len(txs) != 2 {
		t.Errorf("Expected 2 transactions for address1, got %d", len(txs))
	}
	if !reflect.DeepEqual(txs[0], tx1) || !reflect.DeepEqual(txs[1], tx2) {
		t.Errorf("Expected transactions to match tx1 and tx2, got %+v", txs)
	}
}

func TestAddTransactions(t *testing.T) {
	storage := NewMemoryStorage()

	storage.AddObservedAddress("address1")
	storage.AddObservedAddress("address2")

	tx1 := Transaction{
		Hash:      "tx1",
		From:      "address1",
		To:        "address2",
		Value:     "100",
		BlockHash: "blockhash1",
		BlockNum:  1,
		Type:      "transfer",
	}
	tx2 := Transaction{
		Hash:      "tx2",
		From:      "address3",
		To:        "address2",
		Value:     "150",
		BlockHash: "blockhash2",
		BlockNum:  2,
		Type:      "transfer",
	}
	tx3 := Transaction{
		Hash:      "tx3",
		From:      "address1",
		To:        "address4",
		Value:     "200",
		BlockHash: "blockhash3",
		BlockNum:  3,
		Type:      "transfer",
	}

	storage.AddTransactions(tx1, tx2, tx3)

	txsFromAddress1 := storage.GetTransactions("address1")
	if len(txsFromAddress1) != 2 || txsFromAddress1[0].From != "address1" {
		t.Errorf("Expected transaction.From to be address1, got %s", txsFromAddress1[0].From)
	}

	txsFromAddress2 := storage.GetTransactions("address2")
	if len(txsFromAddress2) != 2 || txsFromAddress2[0].To != "address2" {
		t.Errorf("Expected transaction.From to be address2, got %s", txsFromAddress2[0].To)
	}

	txsFromAddress3 := storage.GetTransactions("address3")
	if len(txsFromAddress3) != 0 {
		t.Errorf("Expected no transactions for unobserved address, got %d", len(txsFromAddress3))
	}
}

func TestGetAndUpdateCurrentBlock(t *testing.T) {
	storage := NewMemoryStorage()

	if storage.GetCurrentBlock() != 0 {
		t.Errorf("Expected initial block to be 0, got %d", storage.GetCurrentBlock())
	}

	storage.UpdateCurrentBlock(10)
	if storage.GetCurrentBlock() != 10 {
		t.Errorf("Expected current block to be updated to 10, got %d", storage.GetCurrentBlock())
	}
}
