# Plan: Go Tic-Tac-Toe (go-t3)

**Complexity:** Complex (multi-phase, 37 PRs across 12 phases, client/server architecture)

## Overview

Build a full tic-tac-toe game in Go with:
- A REST API backend supporting 2-player games
- A fyne.io GUI frontend
- Computer opponent using minimax depth-search AI
- JWT authentication and user accounts
- Named-opponent invitations and random matchmaking
- Each phase delivered across multiple PRs + a release tag

## Repository Structure (Final State)

```
go-t3/
├── .github/
│   ├── workflows/
│   │   ├── test.yml        # go test matrix
│   │   ├── lint.yml        # golangci-lint
│   │   └── build.yml       # go build server + client
│   ├── ISSUE_TEMPLATE/
│   │   ├── config.yml
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   ├── CODEOWNERS
│   └── PULL_REQUEST_TEMPLATE.md
├── cmd/
│   ├── server/             # Backend entrypoint
│   │   └── main.go
│   └── client/             # GUI frontend entrypoint
│       └── main.go
├── internal/
│   ├── game/               # Core game logic (board, rules, win detection)
│   ├── api/                # HTTP handlers and routes
│   ├── ai/                 # Computer opponent (minimax)
│   ├── auth/               # JWT auth, user store, middleware
│   ├── matchmaking/        # Random pairing queue
│   └── ui/                 # fyne.io UI components
├── pkg/
│   └── t3/                 # Public types shared between server and client
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CONTRIBUTORS.md
├── SECURITY.md
└── plans/
```

## Phases and PRs

| Phase | PRs | Description | Tag |
|-------|-----|-------------|-----|
| 0 | #1 | Repository bootstrap + GitHub community files | v0.1.0 |
| 1 | #2–#4 | Project structure + GitHub Actions CI | v0.2.0 |
| 2 | #5–#8 | Server API for 2-player game | v0.3.0 |
| 3 | #9–#11 | GUI frontend using API | v0.4.0 |
| 4 | #12–#14 | Computer AI backend (minimax) | v0.5.0 |
| 5 | #15–#16 | GUI computer-play integration | v1.0.0 |
| 6 | #17–#21 | Authentication backend (JWT, register/login) | v1.1.0 |
| 7 | #22–#24 | Authentication frontend + remove --no-auth | v1.2.0 |
| 8 | #25–#28 | Server user lobby + game invitations | v1.3.0 |
| 9 | #29–#31 | GUI opponent selection + invite flow | v1.4.0 |
| 10 | #32–#34 | Server random matchmaking queue | v1.5.0 |
| 11 | #35–#37 | GUI random matchmaking integration | v1.6.0 |

## PR Breakdown

| PR | Branch | Description |
|----|--------|-------------|
| #1 | phase/0-bootstrap | go.mod, main.go, Makefile, GitHub community files |
| #2 | phase/1-server-skeleton | cmd/server, internal/game stub, internal/api stub |
| #3 | phase/1-client-skeleton | cmd/client, internal/ui stub, pkg/t3, fyne dep |
| #4 | phase/1-github-actions | test.yml, lint.yml, build.yml, .golangci.yml |
| #5 | phase/2-game-logic | internal/game expanded + tests |
| #6 | phase/2-api-types | pkg/t3 request/response types |
| #7 | phase/2-create-and-get | POST /games + GET /games/{id} |
| #8 | phase/2-move-endpoint | POST /games/{id}/move |
| #9 | phase/3-api-client | internal/ui/client.go |
| #10 | phase/3-board-widget | internal/ui/board.go |
| #11 | phase/3-game-window | internal/ui/ui.go wiring + --server flag |
| #12 | phase/4-minimax | internal/ai package + tests |
| #13 | phase/4-game-mode-types | GameMode types in pkg/t3 |
| #14 | phase/4-server-ai | API auto-play computer moves |
| #15 | phase/5-board-cell-control | Board per-cell enable/disable + CreateGame sig |
| #16 | phase/5-mode-selector | Mode selector UI + computer player selector |
| #17 | phase/6-no-auth-flag | --no-auth server flag for client compat |
| #18 | phase/6-register | UserStore + POST /auth/register |
| #19 | phase/6-login | POST /auth/login + JWT signing |
| #20 | phase/6-me-endpoint | GET /auth/me + ValidateToken |
| #21 | phase/6-auth-middleware | AuthMiddleware + gate /games/* |
| #22 | phase/7-client-auth-methods | APIClient token storage + auth methods |
| #23 | phase/7-auth-screen | Login/register GUI screen |
| #24 | phase/7-remove-no-auth | Remove --no-auth flag, auth always required |
| #25 | phase/8-online-users | UserStore MarkOnline + GET /users/online |
| #26 | phase/8-game-status-types | GameStatus lifecycle + pkg/t3 invite types |
| #27 | phase/8-invite-endpoints | POST /games/invite + GET /games/invites |
| #28 | phase/8-accept-decline | POST /games/{id}/accept + /decline + tests |
| #29 | phase/9-client-lobby-methods | APIClient lobby/invite methods |
| #30 | phase/9-lobby-widget | LobbyWidget + board playerSide |
| #31 | phase/9-invite-flow | Invite flow wiring in ui.go |
| #32 | phase/10-queue-package | internal/matchmaking package + tests |
| #33 | phase/10-join-leave | POST + DELETE /matchmaking/join |
| #34 | phase/10-status-endpoint | GET /matchmaking/status |
| #35 | phase/11-client-matchmaking-methods | APIClient matchmaking methods |
| #36 | phase/11-matchmaking-widget | MatchmakingWidget (status label + cancel) |
| #37 | phase/11-random-match-ui | Mode selector 4th option + polling flow |

## Phase Files

- [Phase 0](phases/00-phase-0-bootstrap.md) — Bootstrap + GitHub Community Files
- [Phase 1](phases/01-phase-1-structure.md) — Structure + GitHub Actions CI
- [Phase 2](phases/02-phase-2-server-api.md) — Server API
- [Phase 3](phases/03-phase-3-gui-frontend.md) — GUI Frontend
- [Phase 4](phases/04-phase-4-ai-backend.md) — AI Backend
- [Phase 5](phases/05-phase-5-gui-ai-integration.md) — GUI AI Integration
- [Phase 6](phases/06-phase-6-authentication.md) — Authentication Backend
- [Phase 7](phases/07-phase-7-auth-frontend.md) — Authentication Frontend
- [Phase 8](phases/08-phase-8-user-lobby.md) — User Lobby + Invitations
- [Phase 9](phases/09-phase-9-gui-opponent-selection.md) — GUI Opponent Selection
- [Phase 10](phases/10-phase-10-random-matchmaking.md) — Server Random Matchmaking
- [Phase 11](phases/11-phase-11-gui-random-matchmaking.md) — GUI Random Matchmaking
