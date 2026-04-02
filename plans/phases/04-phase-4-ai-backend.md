# Phase 4: Computer AI Backend (Minimax)

**Complexity:** Medium
**PRs:** #12–#14
**Release Tag:** v0.5.0 (on PR #14)
**Branch prefix:** phase/4-

## Goal

Add a computer opponent to the server using a minimax depth-search algorithm. Extend the API to support a "vs computer" game mode where the server automatically plays for the computer player after each human move.

## Context Check
- [x] Phase 3 merged: full 2-player game working via GUI
- [x] internal/game has Board, GameState, ApplyMove, CheckWinner

---

## PR #12 — Minimax AI Package
**Branch:** `phase/4-minimax`

### Algorithm

Classic minimax without alpha-beta pruning (3x3 board is small enough).

```
minimax(state, depth, isMaximizing, computer, human):
  if winner == computer: return +10 - depth
  if winner == human:    return -10 + depth
  if draw:               return 0

  if isMaximizing:
    best = -∞
    for each empty cell:
      apply move (computer); score = minimax(..., false); undo
      best = max(best, score)
    return best
  else:
    best = +∞
    for each empty cell:
      apply move (human); score = minimax(..., true); undo
      best = min(best, score)
    return best
```

### Files

#### `internal/ai/ai.go` (new)
```go
// Package ai provides the minimax computer opponent for go-t3.
package ai

import "github.com/benjamenmeyer/go-t3/internal/game"

// BestMove returns the optimal (row, col) for computerPlayer using minimax.
func BestMove(state game.GameState, computerPlayer game.Cell) (row, col int)

func minimax(state game.GameState, depth int, isMaximizing bool,
    computer, human game.Cell) int
```

#### `internal/ai/ai_test.go` (new)
Tests covering:
- Win-in-one: computer takes the winning move
- Block-in-one: computer blocks human from winning
- Draw scenario: returns a valid move from a near-draw board
- Full game series: computer never loses against random human play

### Verification
```bash
go test ./internal/ai/...
# All tests pass
```

---

## PR #13 — Shared Types for Game Modes
**Branch:** `phase/4-game-mode-types`

### Files

#### `pkg/t3/types.go` (updated)
```go
type GameMode string

const (
    ModeTwoPlayer  GameMode = "two_player"
    ModeVsComputer GameMode = "vs_computer"
)

// CreateGameRequest updated
type CreateGameRequest struct {
    Mode           GameMode `json:"mode"`            // default: two_player
    ComputerPlayer string   `json:"computer_player"` // "X" or "O", default "O"
}

// GameStateResponse updated (add fields)
//   Mode           GameMode `json:"mode"`
//   ComputerPlayer string   `json:"computer_player,omitempty"`
```

### Verification
```bash
go build ./...
```

---

## PR #14 — Server AI Integration
**Branch:** `phase/4-server-ai`

### Files

#### `internal/api/api.go` (updated)
- `GameRecord` gains `Mode GameMode` and `ComputerPlayer game.Cell`
- `createGameHandler`: reads mode + computer_player from request, stores on GameRecord
- `makeMoveHandler`: after successful human `ApplyMove`, if `mode == vs_computer` and game not done and it is now the computer's turn → call `ai.BestMove`, apply result; response reflects full state after both moves

### Verification
```bash
go run ./cmd/server &

GAME=$(curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' \
  -d '{"mode":"vs_computer","computer_player":"O"}' | jq -r '.game_id')

# Human plays X at (0,0); server auto-responds with O's move
curl -s -X POST http://localhost:8080/games/$GAME/move \
  -H 'Content-Type: application/json' \
  -d '{"player":"X","row":0,"col":0}' | jq .
# Board shows both X at (0,0) and computer's O move

kill %1
```

### Post-Phase
- Merge PR #13
- Tag `v0.5.0`
