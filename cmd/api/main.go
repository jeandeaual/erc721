package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"os"

	"erc721"
)

const (
	defaultEndpoint = "https://ropsten.infura.io/v3/bdffc1723bc8468b8bf8879d54a10cbc"
	defaultPort = 80
)

func main() {
	var (
		endpoint string
		port int
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&endpoint, "endpoint", defaultEndpoint, "Connection endpoint")
	flag.IntVar(&port, "port", defaultPort, "Port to listen on")

	flag.Parse()

	client, err := erc721.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Expose an HTTP API
	mux := http.NewServeMux()
	mux.Handle("/", &httpHandler{
		client: client,
	})

	log.Printf("Listening on port %d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Printf("An error occurred when running the server: %s\n", err)
		os.Exit(1)
	}
}
