//go:generate abigen --sol contract/ERC721.sol --pkg main --out contract.go

package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(10000000000)}}, 10000000)

	// Deploy a token contract on the simulated blockchain
	addr, _, contract, err := DeployERC721(auth, sim)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	fmt.Println("Contract deployed at address", addr.String())

	// Commit all pending transactions in the simulator and print the names again
	sim.Commit()

	// Check the balance
	balance, err := contract.BalanceOf(nil, auth.From)
	if err != nil {
		log.Fatalf("Failed to retrieve balance of %v: %v", auth.From, err)
	}

	fmt.Println("Balance:", *balance)

	// Mint a new token
	tx, err := contract.Mint(nil, auth.From, big.NewInt(123))
	if err != nil {
		log.Fatalf("Failed to mint new token to %v: %v", auth.From, err)
	}

	fmt.Println("Token minted, transaction", tx.Hash().String())

	// Retrieve the owner of the token
	owner, err := contract.OwnerOf(nil, big.NewInt(1))
	if err != nil {
		log.Fatalf("Failed to retrieve token owner: %v", err)
	}

	fmt.Println("Token owner:", owner.String())
}
