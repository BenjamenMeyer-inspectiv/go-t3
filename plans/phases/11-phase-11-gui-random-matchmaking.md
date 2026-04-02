# Phase 11: GUI Random Matchmaking Integration

**Complexity:** Simple
**PRs:** #35–#37
**Release Tag:** v1.6.0 (on PR #37)
**Branch prefix:** phase/11-

## Goal

Add "Find Random Opponent" to the GUI. Split into client methods, the waiting widget, and the full mode integration.

## Context Check
- [x] Phase 10 merged: server matchmaking queue working
- [x] Phase 9 merged: GUI has 3-option mode selector, lobby/invite flow

---

## PR #35 — APIClient Matchmaking Methods
**Branch:** `phase/11-client-matchmaking-methods`

### Files

#### `internal/ui/client.go` (updated)
```go
func (c *APIClient) JoinMatchmaking() (*t3.MatchmakingResponse, error)
func (c *APIClient) LeaveMatchmaking() error
func (c *APIClient) MatchmakingStatus() (*t3.MatchmakingResponse, error)
```

All set `Authorization: Bearer` header.

### Verification
```bash
go build ./...
```

---

## PR #36 — MatchmakingWidget
**Branch:** `phase/11-matchmaking-widget`

### Files

#### `internal/ui/matchmaking.go` (new)
```go
// MatchmakingWidget shows "searching for opponent..." with a cancel button.
type MatchmakingWidget struct {
    widget.BaseWidget
    statusLabel *widget.Label
    onCancel    func()
}

func NewMatchmakingWidget(onCancel func()) *MatchmakingWidget

// SetStatus updates the displayed status text.
func (m *MatchmakingWidget) SetStatus(text string)
```

Layout:
```
┌──────────────────────────────┐
│  Looking for an opponent...  │
│                              │
│  [Cancel]                    │
└──────────────────────────────┘
```

### Verification
```bash
go build ./internal/ui/...
```

---

## PR #37 — Mode Selector 4th Option + Polling Flow
**Branch:** `phase/11-random-match-ui`

### Files

#### `internal/ui/ui.go` (updated)
- Mode selector gains 4th option: "Find Random Opponent"
- On [Start] with random mode:
  1. Show `MatchmakingWidget`
  2. Call `client.JoinMatchmaking()`
  3. If response is `matched` → open board immediately
  4. If `waiting` → start polling goroutine (every 2s calls `client.MatchmakingStatus`)
- Polling goroutine: on `matched` → fyne main thread: show board, stop goroutine
- [Cancel] button: call `client.LeaveMatchmaking()`, stop goroutine, return to mode selector
- Board opened with `playerSide` derived from game state `PlayerXUserID` / `PlayerOUserID` vs `client.CurrentUserID()`

### Verification
```bash
# Terminal 1 — server
go run ./cmd/server

# Terminal 2 — alice
go run ./cmd/client
# Login → select "Find Random Opponent" → click [Start]
# → "Looking for an opponent..."

# Terminal 3 — bob
go run ./cmd/client
# Login → select "Find Random Opponent" → click [Start]
# → Both windows immediately show active game board
# Each player can only click their own cells

# Cancel flow:
# Terminal 4 — carol
go run ./cmd/client
# Login → Find Random Opponent → Start → Cancel
# → Returns to mode selector; carol removed from queue
```

### Post-Phase
- Merge PR #37
- Tag `v1.6.0`
