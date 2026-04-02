# Phase 10: Server Random User Matchmaking

**Complexity:** Simple
**PRs:** #32–#34
**Release Tag:** v1.5.0 (on PR #34)
**Branch prefix:** phase/10-

## Goal

Add a server-side matchmaking queue. Users join and the server pairs them randomly, creating an active game. Split into the queue package, join/leave endpoints, and the status endpoint.

## Context Check
- [x] Phase 9 merged: invite/accept flow working, GameStatus lifecycle exists
- [x] UserStore has online tracking (MarkOnline / OnlineUsers)

## Matchmaking Flow
```
POST   /matchmaking/join    → { status: "waiting" | "matched", game_id? }
DELETE /matchmaking/join    → leave queue
GET    /matchmaking/status  → { status: "waiting" | "matched" | "not_queued", game_id? }
```

Queue logic:
- FIFO, in-memory, `sync.Mutex` guarded
- On join: if another user is waiting → pop + pair → create active game (X/O random) → return `matched`
- Otherwise: push self → return `waiting`
- Idempotent: joining twice is a no-op
- Purge offline users (TTL check) on each join attempt

---

## PR #32 — Matchmaking Queue Package
**Branch:** `phase/10-queue-package`

### Files

#### `internal/matchmaking/matchmaking.go` (new)
```go
// Package matchmaking manages the random pairing queue.
package matchmaking

type Queue struct {
    mu    sync.Mutex
    items []string // user IDs in join order
}

func NewQueue() *Queue
func (q *Queue) Join(userID string) (paired bool, partnerID string)
func (q *Queue) Leave(userID string)
func (q *Queue) Status(userID string) (inQueue bool)
func (q *Queue) RemoveOffline(isOnline func(userID string) bool)
```

#### `internal/matchmaking/matchmaking_test.go` (new)
Tests:
- Single user joins → not paired (waiting)
- Two users join → both paired, correct partner IDs returned
- User leaves → removed from queue
- Duplicate join is idempotent (no duplicate entry)
- `RemoveOffline` purges users whose `isOnline` returns false

### Verification
```bash
go test ./internal/matchmaking/...
```

---

## PR #33 — POST + DELETE /matchmaking/join
**Branch:** `phase/10-join-leave`

### Files

#### `pkg/t3/types.go` (updated)
```go
type MatchmakingStatus string

const (
    MatchmakingWaiting   MatchmakingStatus = "waiting"
    MatchmakingMatched   MatchmakingStatus = "matched"
    MatchmakingNotQueued MatchmakingStatus = "not_queued"
)

type MatchmakingResponse struct {
    Status MatchmakingStatus `json:"status"`
    GameID string            `json:"game_id,omitempty"`
}
```

#### `internal/api/api.go` (updated)
- Add `Queue *matchmaking.Queue` to handler struct
- `joinMatchmakingHandler` (POST /matchmaking/join):
  1. Call `queue.RemoveOffline` using `store.OnlineUsers`
  2. Call `queue.Join(userID)`
  3. If paired: create active `GameRecord` with random X/O assignment, return `matched` + game_id
  4. If not: return `waiting`
- `leaveMatchmakingHandler` (DELETE /matchmaking/join): call `queue.Leave(userID)`
- Register routes (both require Bearer token)

#### `cmd/server/main.go` (updated)
Initialise `matchmaking.NewQueue()` and pass to router.

### Verification
```bash
go run ./cmd/server &

# Alice joins
curl -s -X POST http://localhost:8080/matchmaking/join \
  -H "Authorization: Bearer $TOKEN_A" | jq .
# → { "status": "waiting" }

# Bob joins → match
curl -s -X POST http://localhost:8080/matchmaking/join \
  -H "Authorization: Bearer $TOKEN_B" | jq .
# → { "status": "matched", "game_id": "..." }

kill %1
```

---

## PR #34 — GET /matchmaking/status
**Branch:** `phase/10-status-endpoint`

### Files

#### `internal/api/api.go` (updated)
- `matchmakingStatusHandler` (GET /matchmaking/status):
  - Check if user is in queue → `waiting`
  - Check if user has a recently-matched game (GameRecord created by matchmaking with this user) → `matched` + game_id
  - Otherwise → `not_queued`
- Register route (requires Bearer token)

### Verification
```bash
go run ./cmd/server &

curl -s -X POST http://localhost:8080/matchmaking/join \
  -H "Authorization: Bearer $TOKEN_A" | jq .

curl -s http://localhost:8080/matchmaking/status \
  -H "Authorization: Bearer $TOKEN_A" | jq .
# → { "status": "waiting" }

curl -s -X POST http://localhost:8080/matchmaking/join \
  -H "Authorization: Bearer $TOKEN_B" | jq .

curl -s http://localhost:8080/matchmaking/status \
  -H "Authorization: Bearer $TOKEN_A" | jq .
# → { "status": "matched", "game_id": "..." }

kill %1
```

### Post-Phase
- Merge PR #34
- Tag `v1.5.0`
