package main // Use the appropriate package name

import (
	"context"
	"crypto/ecdsa"
	"ether-test/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"math/big"
	"os"
	"testing"
	"time"
)

// TestConnectToTestnet tests the ability to connect to the testnet

func TestTransferEth(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	senderPrivateKey := os.Getenv("SENDER_PRIVATE_KEY")
	if senderPrivateKey == "" {
		t.Fatal("No private key found in .env file")
	}

	recipientAddress := "0x96216849c49358B10257cb55b28eA603c874b05E" // Replace with recipient's address
	amount := big.NewInt(100000000000000)                            // 0.0000001 ETH in Wei

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
		t.Fatalf("Failed to send transaction: %v", err)
	}

	t.Logf("Transaction sent: %s", signedTx.Hash().Hex())
}

func TestWriteMessageToContract(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	// Fetching the network ID
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Assuming rpcURL is defined as a constant or variable that contains your Ethereum testnet RPC URL
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		t.Fatalf("Failed to connect to the testnet: %v", err)
	}
	defer client.Close()

	senderPrivateKey := os.Getenv("SENDER_PRIVATE_KEY")
	if senderPrivateKey == "" {
		t.Fatal("No private key found in .env file")
	}

	privateKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	networkID, err := client.NetworkID(ctx)
	if err != nil {
		t.Fatalf("Failed to get network ID: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, networkID) // Use the correct chain ID
	if err != nil {
		t.Fatalf("Failed to create authorized transactor: %v", err)
	}

	contractAddress := common.HexToAddress("0xBEE4684EA19D09ae89038B09381c548e3202AaaA")
	msgContract, err := contract.NewContract(contractAddress, client)
	if err != nil {
		t.Fatalf("Failed to instantiate the contract: %v", err)
	}

	tx, err := msgContract.WriteMessage(auth, "Hello, Swan Chain!")
	if err != nil {
		t.Fatalf("Failed to send transaction: %v", err)
	}

	t.Logf("Transaction sent! Tx Hash: %s", tx.Hash().Hex())
}

func TestReadMessageFromContract(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("0xBEE4684EA19D09ae89038B09381c548e3202AaaA")
	msgContract, err := contract.NewContract(contractAddress, client)
	if err != nil {
		t.Fatalf("Failed to instantiate the contract: %v", err)
	}

	// Call the readMessage function
	message, err := msgContract.ReadMessage(&bind.CallOpts{})
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	t.Logf("Read message: %s", message)
}
