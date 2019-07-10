package main

import (
	"encoding/json"
	"log"
	"net/http"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"

	"erc721"
)

const (
	defaultContractAddress = "0x273f7F8E6489682Df756151F5525576E322d51A3"
)

// Response is a struct representing the response sent by the API
type Response struct {
	Owner string `json:"owner"`
}

// replyWithError writes an error to the response payload and appropriately sets
// the Content-Type and response code
func replyWithError(statusCode int, message string, w http.ResponseWriter) {
	responseBody, err := json.Marshal(message)
	if err != nil {
		// Marshalling a struct containing only serializable data types
		// cannot fail, but just in case, handle the error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(responseBody)
	if err != nil {
		log.Printf("Couldn't write the error response body: %s\n", err)
	}
}

// httpHandler is the main request handler of the HTTP API
type httpHandler struct {
	client *ethclient.Client
}

// ServerHTTP returns the owner of a token given an ERC721 contract
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// Only GET is supported
		replyWithError(http.StatusMethodNotAllowed, "Only the GET method is supported", w)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	parameters := strings.Split(path, "/")

	log.Println("Received a query with parameters:", parameters)

	if len(parameters) == 0 {
		replyWithError(http.StatusBadRequest, "You need to specify a token", w)
		return
	}

	if len(parameters) > 2 {
		replyWithError(http.StatusBadRequest, "Too many parameters", w)
		return
	}

	var contractAddress, tokenStr string

	if len(parameters) == 2 {
		contractAddress = parameters[0]
		tokenStr = parameters[1]
	} else {
		contractAddress = defaultContractAddress
		tokenStr = parameters[0]
	}

	token, err := strconv.ParseInt(tokenStr, 10, 64)
	if err != nil {
		replyWithError(http.StatusBadRequest, "Invalid token: " + err.Error(), w)
		return
	}

	// Load the contract
	contract, err := erc721.LoadContract(h.client, contractAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Checking for the owner of token %d in contract %s\n", token, contractAddress)

	owner, err := contract.OwnerOf(nil, big.NewInt(token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Owner: owner.String(),
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		// Marshalling a struct containing only serializable data types
		// cannot fail, but just in case, handle the error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		log.Printf("Couldn't write the response body: %s\n", err)
	}
}
