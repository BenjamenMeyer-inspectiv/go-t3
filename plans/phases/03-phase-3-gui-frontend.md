# Phase 3: GUI Frontend Using API

**Complexity:** Medium
**PRs:** 4 planned (3 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v0.4.0 (on final PR)
**Branch prefix:** phase/3-

## Goal

Build a functional fyne.io GUI that connects to the running server to play a 2-player tic-tac-toe game. Players take turns clicking cells on a 3x3 board driven by the server API.

## Context Check
- [x] Phase 2 merged: server API complete and tested
- [x] fyne.io dependency in go.mod from Phase 1

## UI Layout

```
┌──────────────────────────────┐
│  go-t3: Tic Tac Toe          │
│  Turn: X                     │
├────────┬────────┬────────────┤
│   X    │        │   O        │  3x3 Button grid
├────────┼────────┼────────────┤
│        │   X    │            │
├────────┼────────┼────────────┤
│        │        │            │
└────────┴────────┴────────────┘
  [New Game]    Status: ...
```

---

## PR #9 — API Client
**Branch:** `phase/3-api-client`

### Files

#### `internal/ui/client.go` (new)
```go
package ui

// APIClient calls the go-t3 server.
type APIClient struct {
    BaseURL    string
    httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient

// CreateGame calls POST /games, returns game ID.
func (c *APIClient) CreateGame() (string, error)

// GetState calls GET /games/{id}.
func (c *APIClient) GetState(gameID string) (*t3.GameStateResponse, error)

// MakeMove calls POST /games/{id}/move.
func (c *APIClient) MakeMove(gameID, player string, row, col int) (*t3.MoveResponse, error)
```

### Verification
```bash
go build ./...
# No compilation errors
```

---

## PR #10 — Board Widget
**Branch:** `phase/3-board-widget`

### Files

#### `internal/ui/board.go` (new)
```go
package ui

// BoardWidget is a 3x3 grid of fyne buttons.
type BoardWidget struct {
    widget.BaseWidget
    buttons [3][3]*widget.Button
    onPress func(row, col int)
}

func NewBoardWidget(onPress func(row, col int)) *BoardWidget

// Update refreshes button labels from board state [3][3]string and
// disables all buttons when done is true.
func (b *BoardWidget) Update(board [3][3]string, done bool)
```

### Verification
```bash
go build ./internal/ui/...
```

---

## PR #11 — Game Window Wiring
**Branch:** `phase/3-game-window`

### Files

#### `internal/ui/ui.go` (rewritten)
Game flow:
1. `Run()`: create APIClient, create new game via `client.CreateGame()`
2. Lay out board + status label + "New Game" button
3. On cell click: call `client.MakeMove`, update board from response
4. On win/draw: show result in status label, board buttons disabled via `board.Update`
5. "New Game" button: call `client.CreateGame`, reset board

#### `cmd/client/main.go` (updated)
Add `--server` flag:
```go
flag.StringVar(&serverURL, "server", "http://localhost:8080", "server URL")
flag.Parse()
```

### Verification
```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/client
# GUI opens; click cells to play X then O
# Win condition shows "Player X wins!" or "Draw!"
# "New Game" resets the board
```

---

## PR — CLI Parity: Play via Text Client
**Branch:** `phase/3-cli-play`

### Files

#### `cmd/client/main.go` (updated)
Add a `play` subcommand that drives a full 2-player game over stdin/stdout
using the same `internal/ui.APIClient` the GUI uses:
```go
case "play":
    playCLI(client)
```

```go
func playCLI(client *ui.APIClient) {
    gameID, _ := client.CreateGame()
    for {
        state, _ := client.GetState(gameID)
        renderBoard(state.Board)
        if state.Done {
            fmt.Printf("Result: %s\n", resultText(state.Winner))
            return
        }
        fmt.Printf("Turn: %s. Enter move as row,col (0-2): ", state.Current)
        var row, col int
        fmt.Scanf("%d,%d", &row, &col)
        client.MakeMove(gameID, state.Current, row, col)
    }
}

func renderBoard(board [3][3]string) {
    for _, r := range board {
        fmt.Println(strings.Join(r[:], " | "))
    }
}
```

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client play
# Text board renders after each move; alternates X/O prompts
# Win/draw prints "Result: ..." and exits
kill %1
```

### Post-Phase
- Merge final PR (CLI parity)
- Tag `v0.4.0`
