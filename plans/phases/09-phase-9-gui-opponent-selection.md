# Phase 9: GUI Opponent Selection and Game Invitation

**Complexity:** Medium
**PRs:** #29–#31
**Release Tag:** v1.4.0 (on PR #31)
**Branch prefix:** phase/9-

## Goal

Extend the GUI to browse online users, send invitations, and respond to incoming invitations. Split into client methods, the lobby widget, and the full UI wiring.

## Context Check
- [x] Phase 8 merged: server supports invite/accept/decline, online user list, GameStatus lifecycle
- [x] GUI has auth screen and main game window

---

## PR #29 — APIClient Lobby + Invite Methods
**Branch:** `phase/9-client-lobby-methods`

### Files

#### `internal/ui/client.go` (updated)
```go
func (c *APIClient) OnlineUsers() ([]t3.OnlineUser, error)
func (c *APIClient) InviteUser(opponentUserID string, mode t3.GameMode) (string, error) // returns game_id
func (c *APIClient) PendingInvites() ([]t3.PendingInvite, error)
func (c *APIClient) AcceptInvite(gameID string) error
func (c *APIClient) DeclineInvite(gameID string) error
func (c *APIClient) PollGameStatus(gameID string) (t3.GameStatus, error)
```

All methods set `Authorization: Bearer` header.

### Verification
```bash
go build ./...
```

---

## PR #30 — Lobby Widget + Board Player Side
**Branch:** `phase/9-lobby-widget`

### Files

#### `internal/ui/lobby.go` (new)
```go
// LobbyWidget shows online users and pending invitations.
type LobbyWidget struct {
    widget.BaseWidget
    onInvite  func(userID, username string)
    onAccept  func(gameID string)
    onDecline func(gameID string)
}

func NewLobbyWidget(onInvite func(string, string), onAccept func(string), onDecline func(string)) *LobbyWidget

// Refresh updates the displayed user list and invite list.
func (l *LobbyWidget) Refresh(users []t3.OnlineUser, invites []t3.PendingInvite)
```

Layout:
```
Online Users:
  [ alice  ] [Invite]
  [ carol  ] [Invite]
  [Refresh]

Incoming Invitations:
  bob wants to play  [Accept] [Decline]
```

#### `internal/ui/board.go` (updated)
`Update` gains `playerSide` and `currentTurn` so the board disables opponent cells and cells when it is not the local player's turn:
```go
func (b *BoardWidget) Update(board [3][3]string, playerSide string, currentTurn string, done bool)
```

### Verification
```bash
go build ./internal/ui/...
```

---

## PR #31 — Invite Flow Wiring in ui.go
**Branch:** `phase/9-invite-flow`

### Files

#### `internal/ui/ui.go` (updated)
- Mode selector gains third option "Invite Opponent"
- When selected: show `LobbyWidget`
- Lobby auto-refreshes every 5 seconds (goroutine + `widget.Refresh`)
- On invite sent: replace lobby with "Waiting for <name>..." + [Cancel] button
- Polling goroutine: every 2s call `client.PollGameStatus`; on `active` → open board
- [Cancel] on waiting state: call `client.DeclineInvite`, return to lobby
- On incoming invite accepted: call `client.AcceptInvite`, open board for that game ID
- Board opened with `playerSide` derived from game state (`PlayerXUserID` / `PlayerOUserID` vs `client.CurrentUserID()`)

### Verification
```bash
# Terminal 1 — server
go run ./cmd/server

# Terminal 2 — alice
go run ./cmd/client
# Login as alice → select "Invite Opponent" → see bob → click [Invite]
# → "Waiting for bob to accept..."

# Terminal 3 — bob
go run ./cmd/client
# Login as bob → see "alice wants to play" → click [Accept]
# → Both windows open game board
# Take turns; each side only clicks their own cells
```

### Post-Phase
- Merge PR #31
- Tag `v1.4.0`
