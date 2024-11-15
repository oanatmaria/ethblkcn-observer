package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/oanatmaria/ethblkcn-observer/storage"
)

const (
	ethRrpUrl                   = "https://ethereum-rpc.publicnode.com"
	regularTransactionType      = "Regular transaction"
	smartContractDeploymentType = "Contract deployment"
	smartContractExecutionType  = "Contract execution"
)

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RpcError   `json:"error,omitempty"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockResponse struct {
	Number       string              `json:"number"`
	Transactions []TransactionDetail `json:"transactions"`
}

type TransactionDetail struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to,omitempty"`
	Value     string `json:"value"`
	BlockHash string `json:"blockHash"`
}

type EthClient struct {
	baseUrl string
}

func NewEthClient() Client {
	return &EthClient{}
}

func (c *EthClient) GetLatestBlockNumber() (int, error) {
	payload := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	response, err := c.sendRequest(payload)
	if err != nil {
		return 0, err
	}

	blockHex, ok := response.Result.(string)
	if !ok {
		return 0, errors.New("unexpected response format for block number")
	}

	blockNum, err := strconv.ParseInt(blockHex[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %v", err)
	}

	return int(blockNum), nil
}

func (c *EthClient) GetBlockByNumber(blockNum int) (Block, error) {
	blockData, err := c.fetchBlockData(blockNum)
	if err != nil {
		return Block{}, fmt.Errorf("failed to fetch block data: %v", err)
	}

	transactions, err := c.parseTransactions(blockData.Transactions, blockNum)
	if err != nil {
		return Block{}, fmt.Errorf("failed to parse transactions: %v", err)
	}

	// Construct the Block struct to return
	return Block{
		Number:       blockNum,
		Transactions: transactions,
	}, nil
}

func (c *EthClient) fetchBlockData(blockNum int) (BlockResponse, error) {
	payload := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%x", blockNum), true},
		ID:      1,
	}

	response, err := c.sendRequest(payload)
	if err != nil {
		return BlockResponse{}, err
	}

	var block BlockResponse
	err = mapToStruct(response.Result, &block)
	if err != nil {
		return BlockResponse{}, err
	}

	return block, nil
}

func (c *EthClient) parseTransactions(transactionsData []TransactionDetail, blockNum int) ([]storage.Transaction, error) {
	transactions := []storage.Transaction{}
	for _, tx := range transactionsData {
		parsedTx, err := c.parseTransaction(tx, blockNum)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, parsedTx)
	}
	return transactions, nil
}

func (c *EthClient) parseTransaction(txDetail TransactionDetail, blockNum int) (storage.Transaction, error) {
	var txType, toAddress string

	if txDetail.To == "" {
		txType = smartContractDeploymentType
	} else {
		toAddress = txDetail.To
		isSmartContract, err := c.isSmartContract(toAddress)
		if err != nil {
			return storage.Transaction{}, err
		}

		if isSmartContract {
			txType = smartContractExecutionType
		} else {
			txType = regularTransactionType
		}
	}

	return storage.Transaction{
		Hash:      txDetail.Hash,
		From:      txDetail.From,
		To:        toAddress,
		Value:     txDetail.Value,
		BlockHash: txDetail.BlockHash,
		BlockNum:  blockNum,
		Type:      txType,
	}, nil
}

func (c *EthClient) isSmartContract(address string) (bool, error) {
	payload := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getCode",
		Params:  []interface{}{address, "latest"},
		ID:      1,
	}

	response, err := c.sendRequest(payload)
	if err != nil {
		return false, err
	}

	code, ok := response.Result.(string)
	if !ok {
		return false, errors.New("unexpected response format for smart contract code")
	}

	return code != "0x", nil
}

func (c *EthClient) sendRequest(payload RpcRequest) (*RpcResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	var url string
	// only for tests
	if c.baseUrl != "" {
		url = c.baseUrl
	} else {
		url = ethRrpUrl
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected HTTP response: %s - %s", resp.Status, string(body))
	}

	var rpcResponse RpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %d - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return &rpcResponse, nil
}

func mapToStruct(data interface{}, target interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
