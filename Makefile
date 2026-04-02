.PHONY: build run test fmt lint help run-server build-server test-cover

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o bin/go-t3 .

run: ## Run the application
	go run .

test: ## Run tests
	go test ./...

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

run-server: ## Run the server (PORT=8080)
	go run ./cmd/server --port $(or $(PORT),8080)

build-server: ## Build server binary
	go build -o bin/server ./cmd/server

test-cover: ## Run tests with coverage report
	go test -cover ./...
