# Phase 1: Project Structure — Server + Client Skeleton

**Complexity:** Medium
**PRs:** #2–#4
**Release Tag:** v0.2.0 (on PR #4)
**Branch prefix:** phase/1-

## Goal

Establish the directory layout and skeleton code for both the server backend and GUI frontend, then add GitHub Actions CI so every subsequent PR is automatically tested and linted.

## Context Check
- [x] Phase 0 merged: go.mod, main.go, Makefile, all GitHub community files exist
- [x] fyne.io/fyne/v2 will be added as dependency in PR #3

---

## PR #2 — Server Skeleton
**Branch:** `phase/1-server-skeleton`

### Files

#### `cmd/server/main.go` (new)
```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/benjamenmeyer/go-t3/internal/api"
)

func main() {
    fmt.Println("go-t3 server starting on :8080")
    router := api.NewRouter()
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

#### `internal/game/game.go` (new)
```go
// Package game contains core tic-tac-toe board logic.
package game

type Cell int

const (
    Empty Cell = iota
    PlayerX
    PlayerO
)

type Board [3][3]Cell

type GameState struct {
    Board   Board
    Current Cell
    Winner  Cell
    Done    bool
}
```

#### `internal/api/api.go` (new)
```go
// Package api provides the HTTP router and handlers for go-t3.
package api

import "net/http"

func NewRouter() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}
```

#### `main.go` (deleted)
Root-level main.go removed; replaced by `cmd/server`.

#### `Makefile` (updated)
```makefile
run-server: ## Run the server
	go run ./cmd/server

build-server: ## Build server binary
	go build -o bin/server ./cmd/server

test-cover: ## Run tests with coverage report
	go test -cover ./...
```

### Verification
```bash
go build ./...
go run ./cmd/server &
curl http://localhost:8080/health
# Expected: {"status":"ok"}
kill %1
```

---

## PR #3 — Client Skeleton
**Branch:** `phase/1-client-skeleton`

### Dependencies
```bash
go get fyne.io/fyne/v2
go get fyne.io/fyne/v2/app
go get fyne.io/fyne/v2/widget
```

### Files

#### `cmd/client/main.go` (new)
```go
package main

import "github.com/benjamenmeyer/go-t3/internal/ui"

func main() {
    ui.Run()
}
```

#### `internal/ui/ui.go` (new)
```go
// Package ui provides the fyne.io GUI for go-t3.
package ui

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
)

func Run() {
    a := app.New()
    w := a.NewWindow("go-t3: Tic Tac Toe")
    w.SetContent(widget.NewLabel("Tic Tac Toe — coming soon"))
    w.ShowAndRun()
}
```

#### `pkg/t3/types.go` (new)
```go
// Package t3 contains shared public types for go-t3.
package t3

type Player string

const (
    PlayerX  Player = "X"
    PlayerO  Player = "O"
    NoPlayer Player = ""
)
```

#### `Makefile` (updated)
```makefile
run-client: ## Run the GUI client
	go run ./cmd/client

build-client: ## Build client binary
	go build -o bin/client ./cmd/client
```

### Verification
```bash
go build ./...
go run ./cmd/client
# Expected: fyne window opens with "coming soon" label
```

---

## PR #4 — GitHub Actions CI
**Branch:** `phase/1-github-actions`

### Goal

Add three workflow files so every PR and push to `main` is automatically validated. Each workflow is focused on a single concern so failures are immediately obvious.

### Files

#### `.github/workflows/test.yml` (new)
Runs `go test ./...` on every PR and push to `main`. Uses a matrix over the two latest stable Go versions.

```yaml
name: Tests

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    name: go test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go: ["1.25", "1.26"]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test ./...
```

#### `.github/workflows/lint.yml` (new)
Runs `golangci-lint` using the official action. Fails the PR if any linter reports an issue.

```yaml
name: Lint

on:
  push:
    branches: [main]
  pull_request:

jobs:
  lint:
    name: golangci-lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest
```

#### `.github/workflows/build.yml` (new)
Verifies both `cmd/server` and `cmd/client` compile cleanly. Runs on every PR and push to `main`.

```yaml
name: Build

on:
  push:
    branches: [main]
  pull_request:

jobs:
  build:
    name: go build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - name: Build server
        run: go build ./cmd/server
      - name: Build client
        # fyne requires display libs; use virtual framebuffer
        run: |
          sudo apt-get install -y libgl1-mesa-dev xorg-dev
          go build ./cmd/client
```

#### `.golangci.yml` (new)
Linter configuration at the repo root:

```yaml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign

linters-settings:
  gofmt:
    simplify: true

issues:
  exclude-rules:
    # Test files can use unkeyed struct literals
    - path: _test\.go
      linters: [govet]
```

#### `Makefile` (updated)
```makefile
ci: test lint build-server build-client ## Run all CI checks locally
```

### Verification
```bash
# Local lint check
golangci-lint run
# No issues

# Push branch → GitHub shows 3 green checks on PR:
# ✓ Tests (1.22), Tests (1.23)
# ✓ golangci-lint
# ✓ Build server, Build client
```

### Post-Phase
- Merge PR #4
- Tag `v0.2.0`
- All future PRs will have CI status checks
