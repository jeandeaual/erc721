# ERC721

ERC-721 implementation and deployment tool.

## Requirements

* Go 1.12+
* solc
* abigen

## Build

```bash
go generate
go build ./cmd/deploy
```

## Usage

```bash
./deploy PRIVATE_KEY

Flags:
  -gasLimit int
        Gas limit for the deployment (default 300000)
  -gasPrice int
        Gas price for the deployment (defaults to the suggested price)
  -pattern string
        Connection endpoint (default "https://ropsten.infura.io/v3/bdffc1723bc8468b8bf8879d54a10cbc")
```
