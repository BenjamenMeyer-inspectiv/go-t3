# Phase 8: Server User Lobby for Game Setup

**Complexity:** Medium
**PRs:** #25–#28
**Release Tag:** v1.3.0 (on PR #28)
**Branch prefix:** phase/8-

## Goal

Extend the server to track online users and support named-opponent game invitations with an accept/decline lifecycle. Each logical area ships as its own PR.

## Context Check
- [x] Phase 7 merged: JWT auth always required, login screen in client
- [x] UserStore exists, AuthMiddleware gates /games/*

## Game Lifecycle States
```go
type GameStatus string

const (
    GameStatusPending  GameStatus = "pending"   // awaiting accept
    GameStatusActive   GameStatus = "active"    // in progress
    GameStatusDeclined GameStatus = "declined"
    GameStatusDone     GameStatus = "done"
)
```

---

## PR #25 — Online User Registry + GET /users/online
**Branch:** `phase/8-online-users`

### Files

#### `internal/auth/auth.go` (updated)
```go
// MarkOnline updates the user's last-seen timestamp (called by AuthMiddleware).
func (s *UserStore) MarkOnline(userID string)

// OnlineUsers returns users seen within the last 5 minutes, excluding excludeUserID.
func (s *UserStore) OnlineUsers(excludeUserID string) []User
```

#### `internal/auth/middleware.go` (updated)
After successful token validation, call `store.MarkOnline(claims.UserID)`.

#### `internal/auth/auth_test.go` (updated)
Tests:
- `MarkOnline` + `OnlineUsers` returns user within TTL
- User not returned after TTL expires
- Self excluded from results

#### `pkg/t3/types.go` (updated)
```go
type OnlineUser struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
}

type OnlineUsersResponse struct {
    Users []OnlineUser `json:"users"`
}
```

#### `internal/api/api.go` (updated)
- Register `GET /users/online` handler (requires Bearer token)
- Returns `OnlineUsersResponse` of users currently online

### Verification
```bash
go test ./internal/auth/...

go run ./cmd/server &
TOKEN_A=... TOKEN_B=...  # register + login alice and bob

curl -s http://localhost:8080/users/online \
  -H "Authorization: Bearer $TOKEN_A" | jq .
# → bob appears in users list
kill %1
```

---

## PR #26 — GameStatus Lifecycle + Shared Types
**Branch:** `phase/8-game-status-types`

### Files

#### `pkg/t3/types.go` (updated)
```go
type GameStatus string  // pending, active, declined, done

type InviteRequest struct {
    OpponentUserID string   `json:"opponent_user_id"`
    Mode           GameMode `json:"mode"`
    ComputerPlayer string   `json:"computer_player,omitempty"`
}

type InviteResponse struct {
    GameID string     `json:"game_id"`
    Status GameStatus `json:"status"`
}

type PendingInvite struct {
    GameID          string `json:"game_id"`
    InviterUserID   string `json:"inviter_user_id"`
    InviterUsername string `json:"inviter_username"`
    Mode            string `json:"mode"`
}

type PendingInvitesResponse struct {
    Invites []PendingInvite `json:"invites"`
}
```

#### `internal/api/api.go` (updated)
- `GameRecord` gains: `Status GameStatus`, `PlayerXUserID`, `PlayerOUserID`, `InviteeUserID`
- Existing `POST /games` sets `Status = active` immediately (self-started game)
- `makeMoveHandler` rejects moves if `Status != active`

### Verification
```bash
go build ./...
```

---

## PR #27 — POST /games/invite + GET /games/invites
**Branch:** `phase/8-invite-endpoints`

### Files

#### `internal/api/invites.go` (new)
```go
func (h *Handler) inviteUserHandler(w http.ResponseWriter, r *http.Request)
    // Creates GameRecord{Status: pending, InviteeUserID: ...}
    // Errors: cannot invite self, invitee not online

func (h *Handler) listInvitesHandler(w http.ResponseWriter, r *http.Request)
    // Returns all GameRecords where InviteeUserID == current user and Status == pending
```

#### `internal/api/api.go` (updated)
Register routes (both require Bearer token):
```
POST /games/invite
GET  /games/invites
```

### Verification
```bash
go run ./cmd/server &
# alice invites bob
GAME=$(curl -s -X POST http://localhost:8080/games/invite \
  -H "Authorization: Bearer $TOKEN_A" -H 'Content-Type: application/json' \
  -d "{\"opponent_user_id\":\"$BOB_ID\",\"mode\":\"two_player\"}" | jq -r '.game_id')

# bob checks invites
curl -s http://localhost:8080/games/invites \
  -H "Authorization: Bearer $TOKEN_B" | jq .
# → alice's invite appears
kill %1
```

---

## PR #28 — POST /games/{id}/accept + POST /games/{id}/decline + Tests
**Branch:** `phase/8-accept-decline`

### Files

#### `internal/api/invites.go` (updated)
```go
func (h *Handler) acceptInviteHandler(w http.ResponseWriter, r *http.Request)
    // Sets Status = active; only InviteeUserID may accept

func (h *Handler) declineInviteHandler(w http.ResponseWriter, r *http.Request)
    // Sets Status = declined; only InviteeUserID may decline
```

#### `internal/api/api.go` (updated)
Register routes:
```
POST /games/{id}/accept
POST /games/{id}/decline
```

#### `internal/api/invites_test.go` (new)
Tests:
- Invite → status pending
- Accept → status active, moves now allowed
- Decline → status declined, moves rejected
- Cannot invite self → error
- Non-invitee cannot accept → 403
- Cannot move on pending game → error

### Verification
```bash
go test ./internal/api/...

go run ./cmd/server &
# bob accepts alice's invite
curl -s -X POST http://localhost:8080/games/$GAME/accept \
  -H "Authorization: Bearer $TOKEN_B" | jq .
# → { "game_id": "...", "status": "active" }

# alice makes a move
curl -s -X POST http://localhost:8080/games/$GAME/move \
  -H "Authorization: Bearer $TOKEN_A" -H 'Content-Type: application/json' \
  -d '{"player":"X","row":0,"col":0}' | jq .
kill %1
```

### Post-Phase
- Merge PR #28
- Tag `v1.3.0`
