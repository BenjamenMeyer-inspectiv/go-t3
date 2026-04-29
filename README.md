# go-t3

Go-based Tic Tac Toe — a REST API server and fyne.io GUI client.

## Features

- 2-player local game via REST API
- Play against a computer opponent (minimax AI)
- User authentication (JWT)
- Invite a specific opponent or find a random match
- Cross-platform GUI built with [fyne.io](https://fyne.io)
- Testing utility

## Requirements

- Go 1.25+
- [golangci-lint](https://golangci-lint.run) (for linting)

## Build

```bash
# Build server
make build-server

# Build client
make build-client

# Or build both
go build ./...
```

## Usage

```bash
# Start the server (default :8080)
make run-server

# Start the GUI client
make run-client

# Connect client to a non-default server
go run ./cmd/client --server http://localhost:9090
```

## Development

```bash
make test       # run tests
make lint       # run golangci-lint
make fmt        # format code
make help       # list all targets
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for branch naming, PR process, and commit style.
See [SECURITY.md](SECURITY.md) to report a vulnerability.

