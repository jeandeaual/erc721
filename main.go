//go:generate abigen --sol contract/ERC721.sol --pkg main --out contract.go

package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	defaultEndpoint = "https://ropsten.infura.io/v3/bdffc1723bc8468b8bf8879d54a10cbc"
	defaultGasPrice = 0
	defaultGasLimit = 300000
)

func main() {
	var (
		endpoint string
		gasPrice int
		gasLimit int
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s PRIVATE_KEY\n\nFlags:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&endpoint, "pattern", defaultEndpoint, "Connection endpoint")
	flag.IntVar(&gasPrice, "gasPrice", defaultGasPrice, "Gas price for the deployment (defaults to the suggested price)")
	flag.IntVar(&gasLimit, "gasLimit", defaultGasLimit, "Gas limit for the deployment")

	flag.Parse()

	key := flag.Arg(0)
	if len(key) == 0 {
		fmt.Fprintf(os.Stderr, "No private key provided\n")
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Dialing", endpoint)

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Suggested gas price:", suggestedGasPrice)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(gasLimit)
	if gasPrice == defaultGasPrice {
		auth.GasPrice = suggestedGasPrice
	} else {
		auth.GasPrice = big.NewInt(int64(gasPrice))
	}

	addr, tx, _, err := DeployERC721(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy new contract: %v", err)
	}
	fmt.Println("Contract deployed at address", addr.String())
	fmt.Println("Transaction:", tx.Hash().String())
}
