package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type EthClient struct {
}

func NewHttpClient() Client {
	return &EthClient{}
}

func (c *EthClient) GetLatestBlockNumber() (int, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}
	response, err := c.sendRequest(payload)
	if err != nil {
		return 0, err
	}
	blockHex := response["result"].(string)
	blockNum, _ := strconv.ParseInt(blockHex[2:], 16, 64)
	return int(blockNum), nil
}

func (c *EthClient) GetBlockByNumber(blockNum int) (*Block, error) {
	blockData, err := c.fetchBlockData(blockNum)
	if err != nil {
		return nil, err
	}

	transactions, err := c.parseTransactions(blockData["transactions"].([]interface{}), blockNum)
	if err != nil {
		return nil, err
	}

	return &Block{
		Number:       blockNum,
		Transactions: transactions,
	}, nil
}

func (c *EthClient) fetchBlockData(blockNum int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNum), true},
		"id":      1,
	}

	response, err := c.sendRequest(payload)
	if err != nil {
		return nil, err
	}

	return response["result"].(map[string]interface{}), nil
}

func (c *EthClient) parseTransactions(transactionsData []interface{}, blockNum int) ([]storage.Transaction, error) {
	transactions := []storage.Transaction{}
	for _, tx := range transactionsData {
		parsedTx, err := c.parseTransaction(tx.(map[string]interface{}), blockNum)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, parsedTx)
	}
	return transactions, nil
}

func (c *EthClient) parseTransaction(txMap map[string]interface{}, blockNum int) (storage.Transaction, error) {
	var txType, toAddress string

	if txMap["to"] == nil {
		txType = smartContractDeploymentType
	} else {
		toAddress = txMap["to"].(string)
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
		Hash:      txMap["hash"].(string),
		From:      txMap["from"].(string),
		To:        toAddress,
		Value:     txMap["value"].(string),
		BlockHash: txMap["blockHash"].(string),
		BlockNum:  blockNum,
		Type:      txType,
	}, nil
}

func (c *EthClient) isSmartContract(address string) (bool, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getCode",
		"params":  []interface{}{address, "latest"},
		"id":      1,
	}

	response, err := c.sendRequest(payload)
	if err != nil {
		return false, err
	}

	code := response["result"].(string)
	return code != "0x", nil
}

func (c *EthClient) sendRequest(payload map[string]interface{}) (map[string]interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ethRrpUrl, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
