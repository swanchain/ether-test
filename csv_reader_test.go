package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// ReadAddressesFromCSV reads Ethereum addresses from a CSV file.
func ReadAddressesFromCSV(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	isHeader := true

	for scanner.Scan() {
		if isHeader {
			// Skip the header
			isHeader = false
			continue
		}

		line := scanner.Text()
		fields := strings.Split(line, ";")
		if len(fields) > 0 {
			addresses = append(addresses, fields[0])
		}
	}

	return addresses, scanner.Err()
}
func TestSendEthToAddresses(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	senderPrivateKey := os.Getenv("SENDER_PRIVATE_KEY")
	if senderPrivateKey == "" {
		t.Fatal("No private key found in .env file")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		t.Fatalf("Failed to get nonce: %v", err)
	}

	gasLimit := uint64(21000) // Standard limit for a transfer
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("Failed to suggest gas price: %v", err)
	}

	amount := big.NewInt(1e12) // 0.000001 ETH in Wei

	addresses, err := ReadAddressesFromCSV("./ethereum-address/000000000003.csv")
	if err != nil {
		t.Fatalf("Failed to read addresses from CSV: %v", err)
	}

	for _, recipientAddress := range addresses {
		tx := types.NewTransaction(nonce, common.HexToAddress(recipientAddress), amount, gasLimit, gasPrice, nil)

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			t.Fatalf("Failed to get network ID: %v", err)
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			t.Fatalf("Failed to sign transaction: %v", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			t.Fatalf("Failed to send transaction to %s: %v", recipientAddress, err)
		}

		t.Logf("Transaction sent to %s: %s", recipientAddress, signedTx.Hash().Hex())
		nonce++ // Increment nonce for next transaction
		time.Sleep(2 * time.Second)
	}
}
