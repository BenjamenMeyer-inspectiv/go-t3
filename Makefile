.PHONY: build run test fmt lint help

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
