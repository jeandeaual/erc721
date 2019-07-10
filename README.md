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
go build ./cmd/api
```

## Usage

* Deployment

  ```bash
  $ ./deploy PRIVATE_KEY

  Flags:
    -gasLimit int
          Gas limit for the deployment (default 300000)
    -gasPrice int
          Gas price for the deployment (defaults to the suggested price)
    -pattern string
          Connection endpoint (default "https://ropsten.infura.io/v3/bdffc1723bc8468b8bf8879d54a10cbc")
  ```

* API

  ```bash
  $ ./api -port 8080

  2019/07/10 18:13:23 Dialing https://ropsten.infura.io/v3/bdffc1723bc8468b8bf8879d54a10cbc
  2019/07/10 18:13:23 Listening on port 8080
  ```

  ```bash
  $ curl -w "\n" http://localhost:8080/0x273f7F8E6489682Df756151F5525576E322d51A3/50010001

  {"owner":"0xd868711BD9a2C6F1548F5f4737f71DA67d821090"}

  $ curl -w "\n" http://localhost:8080/50010001

  {"owner":"0xd868711BD9a2C6F1548F5f4737f71DA67d821090"}
  ```
