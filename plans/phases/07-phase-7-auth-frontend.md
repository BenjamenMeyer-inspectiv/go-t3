# Phase 7: Authentication Frontend

**Complexity:** Simple
**PRs:** #22–#24
**Release Tag:** v1.2.0 (on PR #24)
**Branch prefix:** phase/7-

## Goal

Wire the GUI client into the auth backend, then permanently remove the `--no-auth` flag now that the client can authenticate. PR #22 adds client auth methods, PR #23 adds the login screen, PR #24 removes the escape hatch.

## Context Check
- [x] Phase 6 merged: server has `/auth/register`, `/auth/login`, `/auth/me`; `/games/*` gated by AuthMiddleware when run without `--no-auth`
- [x] `pkg/t3` has `RegisterRequest`, `LoginRequest`, `LoginResponse`

---

## PR #22 — APIClient Auth Methods + Token Storage
**Branch:** `phase/7-client-auth-methods`

### Files

#### `internal/ui/client.go` (updated)
Add fields to `APIClient`:
```go
type APIClient struct {
    BaseURL     string
    httpClient  *http.Client
    token       string
    currentUser t3.RegisterResponse // stores user_id + username after login
}
```

Add methods:
```go
// Register calls POST /auth/register then Login.
func (c *APIClient) Register(username, password string) error

// Login calls POST /auth/login and stores the returned token.
func (c *APIClient) Login(username, password string) error

func (c *APIClient) CurrentUserID() string
func (c *APIClient) CurrentUsername() string
```

Update all existing methods (`CreateGame`, `GetState`, `MakeMove`) to include:
```go
req.Header.Set("Authorization", "Bearer "+c.token)
```

When `token` is empty the header is still set but the server (if running without `--no-auth`) will return 401 — handled by the login screen in PR #22.

### Verification
```bash
go build ./...
# No compilation errors
```

---

## PR #23 — Login/Register GUI Screen
**Branch:** `phase/7-auth-screen`

### Files

#### `internal/ui/auth.go` (new)
```go
// ShowAuthScreen displays the login/register window.
// Calls onSuccess after a successful login or registration.
func ShowAuthScreen(client *APIClient, onSuccess func()) fyne.Window
```

Layout:
```
┌─────────────────────────┐
│  go-t3 — Login          │
│                         │
│  Username: [          ] │
│  Password: [          ] │
│                         │
│  [Login]  [Register]    │
│                         │
│  Status: ...            │
└─────────────────────────┘
```

Behavior:
- `widget.Entry{Password: true}` for password field
- [Login]: calls `client.Login`; on error shows message in status label
- [Register]: calls `client.Register`; on error shows message
- Both buttons disabled during in-flight request, re-enabled on error
- On success: calls `onSuccess()`

#### `internal/ui/ui.go` (updated)
`Run()` updated:
1. Create `APIClient`
2. Call `ShowAuthScreen(client, func() { showGameWindow(client) })`

Main game window creation extracted to `showGameWindow(client *APIClient)`.

### Verification
```bash
# Terminal 1 — server with no-auth so both modes can be tested
go run ./cmd/server --no-auth &

# Terminal 2 — client uses login screen but token is optional on this server
go run ./cmd/client
# Login screen appears; wrong password → "Invalid credentials"
# Correct credentials → game window opens

# Terminal 3 — server without flag (auth required)
go run ./cmd/server &
go run ./cmd/client
# Must log in; unauthenticated game requests return 401 handled gracefully
kill %1 %2 %3
```

---

## PR #24 — Remove --no-auth Flag
**Branch:** `phase/7-remove-no-auth`

### Goal

The `--no-auth` flag was a temporary shim to keep the client functional while auth was being built. Now that the client has a login screen, the flag serves no purpose. Remove it entirely so auth is always enforced.

### Files

#### `cmd/server/main.go` (updated)
- Remove `noAuth` flag declaration and `flag.BoolVar` call
- Pass a hard-coded `false` (or remove the parameter entirely) when constructing the router, then clean up the parameter

#### `internal/api/api.go` (updated)
- Remove `noAuth bool` parameter from `NewRouter`
- Remove `noAuthMiddleware` function
- Remove conditional — `AuthMiddleware` is always used for `/games/*`

### Verification
```bash
go build ./...

go run ./cmd/server &

# No token → 401
curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq .
# → 401

# --no-auth flag no longer recognised — server exits with unknown flag error
go run ./cmd/server --no-auth
# → flag provided but not defined: -no-auth

# Client flow works end-to-end
go run ./cmd/client
# Login → game window → create game → play
kill %1
```

### Post-Phase
- Merge PR #24
- Tag `v1.2.0`
