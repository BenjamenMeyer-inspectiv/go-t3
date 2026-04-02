# Phase 2: Server API for 2-Player Game

**Complexity:** Medium
**PRs:** #5–#8
**Release Tag:** v0.3.0 (on PR #8)
**Branch prefix:** phase/2-

## Goal

Implement a complete REST API on the server supporting a full 2-player tic-tac-toe game: game creation, state retrieval, move submission, and win/draw detection. Games are stored in memory.

## Context Check
- [x] Phase 1 merged: internal/game, internal/api, pkg/t3 stubs exist

---

## PR #5 — Core Game Logic
**Branch:** `phase/2-game-logic`

### Files

#### `internal/game/game.go` (expanded)
Add to the existing type stubs:
- `NewGame() *GameState` — initializes board, sets Current to PlayerX
- `(g *GameState) ApplyMove(player Cell, row, col int) error`
  - Errors: out of range, wrong turn, cell occupied, game already done
- `(g *GameState) CheckWinner() Cell` — checks all rows, cols, 2 diagonals
- `(g *GameState) IsDraw() bool` — no empty cells and no winner

#### `internal/game/game_test.go` (new)
Tests covering:
- New game initial state (all empty, Current == PlayerX)
- Valid move application and turn alternation
- Win detection: row, column, both diagonals
- Draw detection on full board with no winner
- ApplyMove error cases: wrong turn, occupied cell, out of range, game done

### Verification
```bash
go test ./internal/game/...
# All tests pass
```

---

## PR #6 — Shared API Types
**Branch:** `phase/2-api-types`

### Files

#### `pkg/t3/types.go` (expanded)
```go
type CreateGameRequest struct{}

type CreateGameResponse struct {
    GameID string `json:"game_id"`
}

type GameStateResponse struct {
    GameID  string       `json:"game_id"`
    Board   [3][3]string `json:"board"`   // "X", "O", or ""
    Current string       `json:"current"` // "X" or "O"
    Winner  string       `json:"winner"`  // "X", "O", "draw", or ""
    Done    bool         `json:"done"`
}

type MoveRequest struct {
    Player string `json:"player"` // "X" or "O"
    Row    int    `json:"row"`    // 0-2
    Col    int    `json:"col"`    // 0-2
}

type MoveResponse struct {
    GameStateResponse
    Error string `json:"error,omitempty"`
}
```

### Verification
```bash
go build ./...
```

---

## PR #7 — POST /games and GET /games/{id}
**Branch:** `phase/2-create-and-get`

### Files

#### `internal/api/api.go` (expanded)
Add:
- `GameRecord` struct wrapping `*game.GameState` with a `GameID` string
- In-memory store: `map[string]*GameRecord` with `sync.RWMutex`
- UUID generation helper using `crypto/rand`
- `createGameHandler` — POST /games: allocates new GameState, stores it, returns GameID
- `getGameHandler` — GET /games/{id}: fetches game, serializes board to `[3][3]string`
- Route registration for both endpoints

### Verification
```bash
go run ./cmd/server &

GAME=$(curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq -r '.game_id')
echo "Game: $GAME"

curl -s http://localhost:8080/games/$GAME | jq .
# Expected: empty board, current "X", done false

kill %1
```

---

## PR #8 — POST /games/{id}/move
**Branch:** `phase/2-move-endpoint`

### Files

#### `internal/api/api.go` (expanded)
Add:
- `makeMoveHandler` — POST /games/{id}/move
  - Decodes `MoveRequest`, calls `game.ApplyMove`, calls `CheckWinner` and `IsDraw`
  - Updates `GameRecord.Done` when game ends
  - Returns `MoveResponse` (full board state + optional error string)
- Route registration for move endpoint

### Verification
```bash
go run ./cmd/server &

GAME=$(curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq -r '.game_id')

curl -s -X POST http://localhost:8080/games/$GAME/move \
  -H 'Content-Type: application/json' \
  -d '{"player":"X","row":0,"col":0}' | jq .

curl -s -X POST http://localhost:8080/games/$GAME/move \
  -H 'Content-Type: application/json' \
  -d '{"player":"O","row":1,"col":1}' | jq .

# Invalid move — wrong turn
curl -s -X POST http://localhost:8080/games/$GAME/move \
  -H 'Content-Type: application/json' \
  -d '{"player":"O","row":2,"col":2}' | jq .
# Expected: error field set

kill %1
```

### Post-Phase
- Merge PR #7
- Tag `v0.3.0`
