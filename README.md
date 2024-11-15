# ethblkcn-observer

A HTTP server for observing Ethereum blockchain activities. You can subscribe to an Ethereum addresses, retrieving transactions, and querying the current block.

---

## Project Structure

```
ethblkcn-observer/
├── client/                # HTTP client that handles the calls to Blockchain
├── parser/                # Blockchain parser implementation
├── server/                # HTTP server implementation
├── storage/               # Storage module for blockchain data
├── main.go                # Entry point of the application
├── go.mod                 
├── go.sum                 
└── README.md              
```

## Features

- **Subscribe to Ethereum Addresses:** Allows clients to subscribe to Ethereum addresses to monitor transactions.
- **Retrieve Transactions:** Fetches transactions associated with a given Ethereum address (both from and to).
- **Current Block Information:** Provides the latest block number that was processed by the server corresponding to the block on the Ethereum blockchain.
- **Concurrent Block Processing:** Periodically (every 10 seconds) processes new blocks using a background worker.
- **Starts From Current Block:** The system processes transactions starting from the current block when the server starts. Historical transactions are not handled by default, but this can be easily extended.

---

## Usage

### Start the Server

The server starts on `localhost:8080` by default.

```bash
go run main.go
```

### API Endpoints and Examples

#### Subscribe to an Ethereum Address

Request:

```bash
curl -X POST "http://localhost:8080/subscribe?address=0x1234567890abcdef1234567890abcdef12345678"
```

Successful Response:
```
Subscribed to address: 0x1234567890abcdef1234567890abcdef12345678
```

#### Retrieve Transactions

Request:

```bash
curl -X GET "http://localhost:8080/transactions?address=0x1234567890abcdef1234567890abcdef12345678"
```

Successful Response (JSON):

```
[
    {
        "from": "0xabcdef1234567890abcdef1234567890abcdef12",
        "to": "",
        "value": 1000000000000000000,
        "hash": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef"
        "type": "Contract deployment"
        "blockNum": 21196366
    }
]
```

Here type is the transaction type. There are 3 posible types:
 - Regular transaction (from wallet to wallet)
 - Contract deployment (for smart contracts deployments, the to address will be empty)
 - Contract execution (for smart contracts executions)

 #### Get Current Block

 Request:

 ```bash
 curl -X GET "http://localhost:8080/current_block"
 ```

 Successful response:

```
{
    "block": 12345678
}
```

### Notes on Historical Data
This project does not process historical transactions by default. It starts observing from the current block at the time of server startup. If historical transaction processing is needed, the parser can be extended to include this functionality.