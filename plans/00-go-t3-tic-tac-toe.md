# Plan: Go Tic-Tac-Toe (go-t3)

**Complexity:** Complex (multi-phase, 13 phases, client/server architecture; exact PR counts vary — see PR Numbering Policy below)

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
| 1 | #2–#3 | Project structure (server + client skeleton) | — (moved to 1.5) |
| 1.5 | #4–#17 (actual, retroactive) | Connectivity check + CLI foundation + CI (outstanding) | v0.2.0 (pending) |
| 2 | 4 planned | Server API for 2-player game | v0.3.0 |
| 3 | 4 planned (incl. 1 CLI parity) | GUI + CLI frontend using API | v0.4.0 |
| 4 | 3 planned | Computer AI backend (minimax) | v0.5.0 |
| 5 | 3 planned (incl. 1 CLI parity) | GUI + CLI computer-play integration | v1.0.0 |
| 6 | 5 planned | Authentication backend (JWT, register/login) | v1.1.0 |
| 7 | 4 planned (incl. 1 CLI parity) | Authentication GUI + CLI + remove --no-auth | v1.2.0 |
| 8 | 4 planned | Server user lobby + game invitations | v1.3.0 |
| 9 | 4 planned (incl. 1 CLI parity) | GUI + CLI opponent selection + invite flow | v1.4.0 |
| 10 | 3 planned | Server random matchmaking queue | v1.5.0 |
| 11 | 4 planned (incl. 1 CLI parity) | GUI + CLI random matchmaking integration | v1.6.0 |

## PR Breakdown (Historical — Phases 0, 1, 1.5)

| PR | Branch | Description |
|----|--------|-------------|
| #1 | phase/0-bootstrap | go.mod, main.go, Makefile, GitHub community files |
| #2 | phase/1-server-skeleton | cmd/server, internal/game stub, internal/api stub |
| #3 | phase/1-client-skeleton | cmd/client, internal/ui stub, pkg/t3, fyne dep |
| #4–#5 | BenjamenMeyer-inspectiv-patch-1 | dependabot.yml create/update |
| #6 | enhancement-gitignore | .gitignore config |
| #7–#13 | enhancement-more, enhancement-less, cleanup_1, cleanup_2, enhancement_ping2, enhancement_ping-reset, enhancement-pong2 | Iterative /ping, /pong, /ping/reset build-out |
| #14 | enhancement_basic-client | CLI dispatch (ping/pong/shell) in cmd/client/main.go |
| #15 | enhancment_cli_add-banner | CLI banner output |
| #16 | enhancement-update-readme-1 | README line addition |
| #17 | dependabot/go_modules/fyne.io/fyne/v2-2.7.4 | fyne 2.7.3 → 2.7.4 bump |

Full detail for the connectivity-check and CLI-foundation work is in
[Phase 1.5](phases/01a-phase-1.5-connectivity-cli-foundation.md).

## PR Numbering Policy (Phase 2 onward)

Phase files describe planned PRs by purpose and order, not by predicted
GitHub number — numbers get assigned automatically at PR creation and drift
whenever off-plan work merges in between (as already happened once: PR #4
was predicted to be CI workflows, actually landed as a dependabot patch).
Once a phase's PRs merge, its phase file gains a short
`**Actual PRs:** #N, #N+1, ...` line recording what really happened.

## Phase Files

- [Phase 0](phases/00-phase-0-bootstrap.md) — Bootstrap + GitHub Community Files
- [Phase 1](phases/01-phase-1-structure.md) — Structure (CI deferred to 1.5)
- [Phase 1.5](phases/01a-phase-1.5-connectivity-cli-foundation.md) — Connectivity Check + CLI Foundation
- [Phase 2](phases/02-phase-2-server-api.md) — Server API
- [Phase 3](phases/03-phase-3-gui-frontend.md) — GUI + CLI Frontend
- [Phase 4](phases/04-phase-4-ai-backend.md) — AI Backend
- [Phase 5](phases/05-phase-5-gui-ai-integration.md) — GUI + CLI AI Integration
- [Phase 6](phases/06-phase-6-authentication.md) — Authentication Backend
- [Phase 7](phases/07-phase-7-auth-frontend.md) — Authentication GUI + CLI Frontend
- [Phase 8](phases/08-phase-8-user-lobby.md) — User Lobby + Invitations
- [Phase 9](phases/09-phase-9-gui-opponent-selection.md) — GUI + CLI Opponent Selection
- [Phase 10](phases/10-phase-10-random-matchmaking.md) — Server Random Matchmaking
- [Phase 11](phases/11-phase-11-gui-random-matchmaking.md) — GUI + CLI Random Matchmaking
