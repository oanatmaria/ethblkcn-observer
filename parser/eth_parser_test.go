package parser_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/oanatmaria/ethblkcn-observer/client"
	"github.com/oanatmaria/ethblkcn-observer/parser"
	"github.com/oanatmaria/ethblkcn-observer/storage"
)

func TestNewEthParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := client.NewMockClient(ctrl)
	mockStorage := storage.NewMockStorage(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)

	ethParser, err := parser.NewEthParser(mockStorage, mockClient)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ethParser == nil {
		t.Errorf("expected ethParser to be non-nil")
	}
}

func TestNewEthParser_ErrorFetchingLatestBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := client.NewMockClient(ctrl)
	mockStorage := storage.NewMockStorage(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(0, errors.New("network error"))

	ethParser, err := parser.NewEthParser(mockStorage, mockClient)
	if err == nil {
		t.Errorf("expected an error but got none")
	}
	if ethParser != nil {
		t.Errorf("expected ethParser to be nil")
	}
}

func TestEthParser_GetCurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockClient := client.NewMockClient(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)
	mockStorage.EXPECT().GetCurrentBlock().Return(100)

	ethParser, _ := parser.NewEthParser(mockStorage, mockClient)
	currentBlock := ethParser.GetCurrentBlock()
	if currentBlock != 100 {
		t.Errorf("expected currentBlock to be 100, got %d", currentBlock)
	}
}

func TestEthParser_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockClient := client.NewMockClient(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)
	mockStorage.EXPECT().AddObservedAddress("0xAddress").Return(true)

	ethParser, _ := parser.NewEthParser(mockStorage, mockClient)
	result := ethParser.Subscribe("0xAddress")
	if !result {
		t.Errorf("expected Subscribe to return true")
	}
}

func TestEthParser_GetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockClient := client.NewMockClient(ctrl)

	transactions := []storage.Transaction{
		{Hash: "tx1"},
		{Hash: "tx2"},
	}

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)
	mockStorage.EXPECT().GetTransactions("0xAddress").Return(transactions)

	ethParser, _ := parser.NewEthParser(mockStorage, mockClient)
	result := ethParser.GetTransactions("0xAddress")
	if len(result) != len(transactions) {
		t.Errorf("expected %d transactions, got %d", len(transactions), len(result))
	}
	for i, tx := range result {
		if tx != transactions[i] {
			t.Errorf("expected transaction %v, got %v", transactions[i], tx)
		}
	}
}

func TestEthParser_ProcessNewBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockClient := client.NewMockClient(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)
	mockStorage.EXPECT().GetCurrentBlock().Return(100)
	mockClient.EXPECT().GetLatestBlockNumber().Return(105, nil)

	for i := 101; i <= 105; i++ {
		mockClient.EXPECT().GetBlockByNumber(i).Return(client.Block{
			Transactions: []storage.Transaction{{Hash: fmt.Sprintf("tx%d", i)}},
		}, nil)
		mockStorage.EXPECT().AddTransactions(storage.Transaction{Hash: fmt.Sprintf("tx%d", i)})
	}

	mockStorage.EXPECT().UpdateCurrentBlock(105)

	ethParser, _ := parser.NewEthParser(mockStorage, mockClient)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ethParser.ProcessNewBlocks(ctx)
}

func TestEthParser_ProcessNewBlocks_ErrorFetchingBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	mockClient := client.NewMockClient(ctrl)

	mockClient.EXPECT().GetLatestBlockNumber().Return(100, nil)
	mockStorage.EXPECT().UpdateCurrentBlock(100)
	mockStorage.EXPECT().GetCurrentBlock().Return(100)
	mockClient.EXPECT().GetLatestBlockNumber().Return(105, nil)

	mockClient.EXPECT().GetBlockByNumber(101).Return(client.Block{}, errors.New("block fetch error")).Times(1)
	mockClient.EXPECT().GetBlockByNumber(gomock.Any()).Return(client.Block{
		Transactions: []storage.Transaction{{Hash: "tx"}},
	}, nil).AnyTimes()

	mockStorage.EXPECT().AddTransactions(gomock.Any()).AnyTimes()
	mockStorage.EXPECT().UpdateCurrentBlock(105)

	ethParser, _ := parser.NewEthParser(mockStorage, mockClient)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ethParser.ProcessNewBlocks(ctx)
}
