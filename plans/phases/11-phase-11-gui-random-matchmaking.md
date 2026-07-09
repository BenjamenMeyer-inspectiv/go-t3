# Phase 11: GUI Random Matchmaking Integration

**Complexity:** Simple
**PRs:** 4 planned (3 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.6.0 (on final PR)
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

---

## PR — CLI Parity: Matchmaking Subcommands
**Branch:** `phase/11-cli-matchmaking`

### Files

#### `cmd/client/main.go` (updated)
```go
case "queue":
    if len(args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: cli queue <join|leave|status>")
        os.Exit(1)
    }
    switch args[1] {
    case "join":
        resp, _ := client.JoinMatchmaking()
        fmt.Printf("status: %s\n", resp.Status)
    case "leave":
        client.LeaveMatchmaking()
        fmt.Println("left queue")
    case "status":
        resp, _ := client.MatchmakingStatus()
        fmt.Printf("status: %s\n", resp.Status)
    }
```

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client queue join
# → "status: waiting" or "status: matched"

go run ./cmd/client queue status
go run ./cmd/client queue leave
kill %1
```

### Post-Phase
- Merge final PR (CLI parity)
- Tag `v1.6.0`
