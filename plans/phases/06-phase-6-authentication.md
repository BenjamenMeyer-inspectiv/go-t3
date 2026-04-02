# Phase 6: Authentication Backend

**Complexity:** Medium
**PRs:** #17–#21
**Release Tag:** v1.1.0 (on PR #21)
**Branch prefix:** phase/6-

## Goal

Add user registration and login to the server with JWT-based authentication. An early PR introduces a `--no-auth` flag so the existing GUI client continues to work against the server throughout this phase. Each endpoint and the middleware then ship as separate PRs.

## Context Check
- [x] Phase 5 merged: full game with AI working
- [x] No auth library exists yet

## Dependencies (added in PR #17)
```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

## Auth Design

```
Register:  POST /auth/register  { username, password } → { user_id, username }
Login:     POST /auth/login     { username, password } → { token, user_id, username }
Me:        GET  /auth/me                               → { user_id, username }
```

JWT claims: `UserID`, `Username`, 24h expiry.
Secret from `JWT_SECRET` env var; dev default if unset.

In-memory UserStore:
```go
type User struct {
    ID           string
    Username     string
    PasswordHash string // bcrypt
}
```

---

## PR #17 — --no-auth Server Flag
**Branch:** `phase/6-no-auth-flag`

### Goal

Introduce a `--no-auth` flag on the server so the GUI client (which has no login screen yet) can still connect and play games while auth endpoints are being built. When the flag is set, all middleware auth checks are skipped and a synthetic anonymous identity is injected into the request context.

### Files

#### `cmd/server/main.go` (updated)
```go
var noAuth bool
flag.BoolVar(&noAuth, "no-auth", false, "disable authentication (development only)")
flag.Parse()
```
Pass `noAuth` to `api.NewRouter(...)`.

#### `internal/api/api.go` (updated)
`NewRouter` accepts a `noAuth bool` parameter. When true, a pass-through middleware is used in place of `AuthMiddleware` and injects a fixed anonymous `UserID` into the context so downstream handlers that read user identity do not panic.

```go
func noAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), contextKeyUserID, "anon")
        ctx = context.WithValue(ctx, contextKeyUsername, "anon")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Verification
```bash
# With flag — client works without token
go run ./cmd/server --no-auth &
curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq .
# → game created successfully

# Without flag — should fail once middleware is wired (PR #21)
go run ./cmd/server &
# (games endpoint open until PR #20 gates it)
kill %1 %2
```

---

## PR #18 — User Store + POST /auth/register
**Branch:** `phase/6-register`

### Files

#### `internal/auth/auth.go` (new)
```go
// Package auth handles user management and JWT for go-t3.
package auth

type User struct {
    ID           string
    Username     string
    PasswordHash string
}

type UserStore struct { /* sync.RWMutex + map */ }

func NewUserStore(jwtSecret string) *UserStore
func (s *UserStore) Register(username, password string) (*User, error)
    // errors: username already taken, bcrypt failure
```

#### `internal/auth/auth_test.go` (new)
Tests:
- Register new user success
- Register duplicate username → error

#### `pkg/t3/types.go` (updated)
```go
type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterResponse struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
}
```

#### `internal/api/api.go` (updated)
- Accept `*auth.UserStore` as dependency
- Register `POST /auth/register` handler calling `store.Register`

#### `cmd/server/main.go` (updated)
```go
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    secret = "dev-secret-change-in-production"
}
store := auth.NewUserStore(secret)
router := api.NewRouter(store, noAuth)
```

### Verification
```bash
go test ./internal/auth/...

go run ./cmd/server --no-auth &
curl -s -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123"}' | jq .
# → { "user_id": "...", "username": "alice" }

# Duplicate → error
curl -s -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"other"}' | jq .
kill %1
```

---

## PR #19 — POST /auth/login
**Branch:** `phase/6-login`

### Files

#### `internal/auth/auth.go` (updated)
```go
type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func (s *UserStore) Login(username, password string) (string, error)
    // bcrypt compare, sign JWT, return token string
    // errors: user not found, wrong password
```

#### `internal/auth/auth_test.go` (updated)
Tests:
- Login success → valid JWT returned
- Login wrong password → error
- Login unknown user → error

#### `pkg/t3/types.go` (updated)
```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token    string `json:"token"`
    UserID   string `json:"user_id"`
    Username string `json:"username"`
}
```

#### `internal/api/api.go` (updated)
- Register `POST /auth/login` handler calling `store.Login`

### Verification
```bash
go test ./internal/auth/...

go run ./cmd/server --no-auth &
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')
echo "Token: $TOKEN"
kill %1
```

---

## PR #20 — GET /auth/me
**Branch:** `phase/6-me-endpoint`

### Files

#### `internal/auth/auth.go` (updated)
```go
func (s *UserStore) ValidateToken(tokenStr string) (*Claims, error)
```

#### `internal/auth/auth_test.go` (updated)
Tests:
- ValidateToken success
- ValidateToken expired/invalid → error

#### `internal/api/api.go` (updated)
- Register `GET /auth/me` handler: extract `Authorization: Bearer <token>`, call `ValidateToken`, return user info

### Verification
```bash
go run ./cmd/server --no-auth &
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')

curl -s http://localhost:8080/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq .
# → { "user_id": "...", "username": "alice" }
kill %1
```

---

## PR #21 — Auth Middleware + Gate /games/* Routes
**Branch:** `phase/6-auth-middleware`

### Files

#### `internal/auth/middleware.go` (new)
```go
// AuthMiddleware requires a valid Bearer token.
// On success, injects UserID and Username into request context.
// On failure, returns 401 Unauthorized.
func AuthMiddleware(store *UserStore, next http.Handler) http.Handler
```

#### `internal/api/api.go` (updated)
When `noAuth == false`, wrap all `/games/*` routes with `AuthMiddleware`.
When `noAuth == true`, use `noAuthMiddleware` (from PR #16) instead.

### Verification
```bash
# Auth required (no flag)
go run ./cmd/server &
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')

curl -s -X POST http://localhost:8080/games \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' -d '{}' | jq .
# → game created

curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq .
# → 401 Unauthorized

# --no-auth still works (for client compat during Phase 7)
go run ./cmd/server --no-auth &
curl -s -X POST http://localhost:8080/games \
  -H 'Content-Type: application/json' -d '{}' | jq .
# → game created without token
kill %1 %2
```

### Post-Phase
- Merge PR #21
- Tag `v1.1.0`
