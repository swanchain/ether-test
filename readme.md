# Go Ethereum Test Suite

## Overview

This project contains a suite of tests written in Go for interacting with Ethereum smart contracts and performing various blockchain-related operations. The tests demonstrate how to connect to the Ethereum network, send transactions, and interact with a custom smart contract.

## Prerequisites

- Golang: Ensure you have Go installed on your system. [Download Go](https://golang.org/dl/)
- Ethereum Node: Access to an Ethereum node is required. This can be either a local node or a remote service 
- Solidity Contract: A deployed instance of the smart contract on the Ethereum network.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/filswan/ether-test
   ```
2. Navigate to the project directory:
   ```bash
   cd ether-test
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Set up your `.env` file with the necessary environment variables:

```
SENDER_PRIVATE_KEY=your_private_key
```

## Test Files

- `ethclient_test.go`: Demonstrates how to connect to an Ethereum client and fetches basic blockchain data like the latest block number and network ID.

- `ethclient_fund_test.go`: Contains tests for transferring Ether to an Ethereum address. It demonstrates how to construct, sign, and send transactions.

- `config_test.go`: (Provide a brief description of what this file does).

## Running Tests

To run the tests, use the following command:

```bash
go test -v ./...
```
