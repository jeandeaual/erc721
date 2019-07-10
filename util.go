//go:generate abigen --sol contract/ERC721.sol --pkg erc721 --out contract.go

package erc721

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Connect to an Ethereum endpoint
func Connect(endpoint string) (*ethclient.Client, error) {
	log.Println("Dialing", endpoint)

	return ethclient.Dial(endpoint)
}

// LoadContract loads an ERC721 contract from an address
func LoadContract(client *ethclient.Client, hex string) (*ERC721, error) {
	address := common.HexToAddress(hex)

	return NewERC721(address, client)
}
