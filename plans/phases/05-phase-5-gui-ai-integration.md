# Phase 5: GUI Computer-Play Integration

**Complexity:** Simple
**PRs:** 3 planned (2 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.0.0 (on final PR)
**Branch prefix:** phase/5-

## Goal

Extend the GUI to let the player choose to play against the computer. Add a mode selector and per-cell enable/disable logic so the human can only click their own turn cells.

## Context Check
- [x] Phase 4 merged: server supports `vs_computer` mode + auto AI move
- [x] GUI currently always creates two_player games

---

## PR #15 — Board Cell Enable/Disable + CreateGame Update
**Branch:** `phase/5-board-cell-control`

### Files

#### `internal/ui/board.go` (updated)
Add per-cell enable control to `Update`:
```go
// Update refreshes labels and enables only cells where enabled[row][col] is true.
func (b *BoardWidget) Update(board [3][3]string, enabled [3][3]bool)
```
Empty cells where `enabled` is false are rendered as disabled buttons (opponent's turn or game over).

#### `internal/ui/client.go` (updated)
`CreateGame` signature updated to accept mode and computer player:
```go
func (c *APIClient) CreateGame(mode, computerPlayer string) (string, error)
```

### Verification
```bash
go build ./...
```

---

## PR #16 — Mode Selector UI
**Branch:** `phase/5-mode-selector`

### Files

#### `internal/ui/ui.go` (updated)
- Add `modeSelect` widget (`widget.RadioGroup`: "2 Player" / "vs Computer")
- Add `computerPlayerSelect` widget (`widget.RadioGroup`: "X" / "O") shown only when vs Computer selected
- Pass selected mode + computerPlayer to `client.CreateGame`
- After move response: compute `enabled` grid — only empty cells on the human player's turn; pass to `board.Update`
- Status label logic:
  - "Your turn (X)" / "Your turn (O)"
  - "You win!" / "Computer wins!" / "Draw!" based on `ComputerPlayer` field vs winner

### Verification
```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/client

# 2-player mode: existing behavior unchanged

# vs Computer, computer plays O:
# Click a cell as X → board shows X move + computer's O response automatically
# Play to completion → correct win/loss/draw message shown

# vs Computer, computer plays X:
# Computer immediately places at start (X goes first)
# Human plays O
```

---

## PR — CLI Parity: Mode Selection
**Branch:** `phase/5-cli-mode-selector`

### Files

#### `cmd/client/main.go` (updated)
`play` subcommand gains a `--mode` flag:
```go
flag.StringVar(&mode, "mode", "two_player", "two_player or vs_computer")
flag.StringVar(&computerPlayer, "computer-player", "O", "X or O, only used with vs_computer")
```
When `vs_computer`, `CreateGame` is called with `(mode, computerPlayer)` per
the Phase 5 GUI signature change, and `playCLI` skips prompting for moves on
the computer's turn — it just re-polls `GetState` until `Current` matches
the human side.

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client play --mode=vs_computer --computer-player=O
# Human plays X, computer auto-responds as O after each human move
kill %1
```

### Post-Phase
- Merge final PR (CLI parity)
- Tag `v1.0.0`

## v1.0.0 Release Notes Template

```
## go-t3 v1.0.0

Tic Tac Toe in Go with fyne.io GUI.

### Features
- 2-player local game via REST API
- Play against computer (minimax AI, never loses)
- fyne.io cross-platform GUI

### Usage
    go run ./cmd/server
    go run ./cmd/client
```
