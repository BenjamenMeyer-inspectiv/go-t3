# Plan Restructure: Connectivity Check + CLI Parity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fold the already-merged ping/pong connectivity checks and CLI shell into `plans/` as intentional work (new Phase 1.5), and add CLI-parity PR entries to every remaining GUI-facing phase, per the approved design doc.

**Architecture:** This is a documentation-only restructuring of `plans/00-go-t3-tic-tac-toe.md` and 6 phase files. No source code changes. Each task edits one file to exact, final content — no code is written or run, since the deliverable is plan text, not software. "Verification" per task means grepping the edited file for the expected strings and confirming no leftover references to the old content.

**Tech Stack:** Markdown. Git for commits.

**Design doc:** `docs/superpowers/specs/2026-07-09-plan-restructure-cli-connectivity-design.md`

## Global Constraints

- No source code changes in this plan — `plans/` and `docs/` files only.
- No placeholders (TBD/TODO) in any edited file — every PR block must have concrete file paths, code, and verification commands, per this project's plan-writing convention (see `plans/phases/00-phase-0-bootstrap.md` for the house style).
- Phases 0 and 1 keep their real, already-resolved GitHub PR numbers — do not touch historical numbering.
- Phase 1.5 onward: no predicted exact GitHub PR numbers in phase files (per design decision 5). Use "N planned" + descriptions instead.
- Each edit must preserve the existing phase-file section structure (`**Complexity:**`, `**PRs:**`, `**Release Tag:**`, `**Branch prefix:**`, `## Goal`, `## Context Check`, PR blocks, `### Post-Phase`).
- Branch `phase/1-github-actions` is left untouched (no git action) — only referenced in the new Phase 1.5 doc as superseded.
- Commit each task separately (docs commits, not code — still needs its own commit per house rule of frequent, atomic commits).

---

### Task 1: Create Phase 1.5 file

**Files:**
- Create: `plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`

**Interfaces:**
- Produces: the Phase 1.5 doc that Task 2 (master plan) links to from its Phases table and Phase Files list.

- [ ] **Step 1: Write the file**

```markdown
# Phase 1.5: Connectivity Check + CLI Foundation

**Complexity:** Simple
**Actual PRs:** #4–#17 (retroactively documented; not originally planned)
**Release Tag:** v0.2.0 (pending — see Outstanding Items)
**Branch prefix:** none (branches used ad hoc: `enhancement-*`, `cleanup_*`, `dependabot/*`, `*-patch-1`)

## Goal

Formalize work that merged ahead of the original plan (PRs #4-#17) as an
intentional bridging phase between Phase 1 (structure) and Phase 2 (server
API): a server-side connectivity-check surface, and the foundation of a
headless CLI client that gains feature parity with the GUI phase-by-phase
from here on.

## Context Check
- [x] Phase 1 merged: PR #2 (server skeleton), PR #3 (client skeleton)
- [x] PRs #4-#17 already merged on `main`, off the original plan
- [x] See design doc: `docs/superpowers/specs/2026-07-09-plan-restructure-cli-connectivity-design.md`

---

## Connectivity Check (retroactive)

**Actual PRs:** #6–#13 (`enhancement-gitignore`, `enhancement-more`,
`enhancement-less`, `cleanup_1`, `cleanup_2`, `enhancement_ping2`,
`enhancement_ping-reset`, `enhancement-pong2`)

### What shipped
`internal/api/api.go` gained three endpoints alongside `/health`:
```go
mux.HandleFunc("/ping", pingHandler)
mux.HandleFunc("/ping/reset", pingResetHandler)
mux.HandleFunc("/pong", pongHandler)
```

### Formalized purpose
These are the connectivity-check surface both clients (GUI and CLI) probe
before attempting gameplay calls — confirms the server is reachable and
responding before a game session starts. `/health` remains the
machine-readable liveness check; `/ping`/`/pong` are the human/CLI-facing
round-trip check.

### Verification
```bash
go run ./cmd/server &
curl -s http://localhost:8080/health
# → {"status":"ok"}
curl -s http://localhost:8080/ping
# → {"status":"ok"}
curl -s http://localhost:8080/pong
# → {"status":"ok"}
kill %1
```

---

## CLI Foundation (retroactive)

**Actual PRs:** #14–#15 (`enhancement_basic-client`, `enhancment_cli_add-banner`)

### What shipped
`cmd/client/main.go` gained a flag-based CLI dispatch:
```go
flag.StringVar(&host, "host", "localhost", "server host")
flag.IntVar(&port, "port", 8080, "server port")
// ...
switch args[0] {
case "ping":
    // GET /ping
case "pong":
    // GET /pong
case "shell":
    ui.Run()
}
```

### Formalized purpose
This is the entrypoint the headless CLI client builds on for the rest of
the roadmap. Every future GUI-facing phase (3, 5, 7, 9, 11) adds one
CLI-parity PR here, reusing the shared `APIClient` from
`internal/ui/client.go` rather than duplicating HTTP logic. See the master
plan's PR Breakdown for the parity schedule.

### Verification
```bash
go run ./cmd/client ping
# → status: 200 OK / body: {"status":"ok"}
go run ./cmd/client shell
# → fyne window opens (delegates to ui.Run())
```

---

## Other Merged Work (retroactive, no formal follow-up needed)

- PRs #4–#5 (`BenjamenMeyer-inspectiv-patch-1`) — dependabot.yml create/update, config only.
- PR #16 (`enhancement-update-readme-1`) — README line addition, documentation only.
- PR #17 (`dependabot/go_modules/fyne.io/fyne/v2-2.7.4`) — routine dependency bump, no plan impact.

## Outstanding Items

- **CI still not implemented.** The original Phase 1 PR #4 concept
  (`.github/workflows/test.yml`, `lint.yml`, `build.yml`, `.golangci.yml`,
  documented in full in `plans/phases/01-phase-1-structure.md`) never
  merged, even though 13 more PRs landed after it without CI gating them.
  This is deferred, not abandoned — pick it up as the next PR in this phase
  before tagging.
- **Branch `phase/1-github-actions` is superseded.** It contains a prior CI
  attempt that also reverts ping/pong and the CLI shell back to bare
  skeletons — the opposite of this phase's direction. Do not merge it as-is;
  if reused, cherry-pick only the CI workflow files, not the reverts.

## Post-Phase
- Implement CI (see Outstanding Items above; content already specified in
  `plans/phases/01-phase-1-structure.md`'s PR #4 block)
- Tag `v0.2.0`
```

- [ ] **Step 2: Verify file content**

Run: `grep -c "^## " plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`
Expected: `6` (Goal, Context Check, Connectivity Check, CLI Foundation, Other Merged Work, Outstanding Items, Post-Phase — confirms all sections present; if the count differs, re-check the file against the content above)

Run: `grep -n "phase/1-github-actions" plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`
Expected: one match, in the Outstanding Items section

- [ ] **Step 3: Commit**

```bash
git add plans/phases/01a-phase-1.5-connectivity-cli-foundation.md
git commit -m "docs: add Phase 1.5 retroactively documenting connectivity + CLI work

Formalizes PRs #4-#17 (previously off-plan) as the connectivity-check
surface and CLI client foundation, per the approved design.

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 2: Update master plan — Phases table, Phase Files list, PR Breakdown, numbering policy

**Files:**
- Modify: `plans/00-go-t3-tic-tac-toe.md`

**Interfaces:**
- Consumes: Phase 1.5 file path from Task 1.
- Produces: the "PR Numbering Policy" section that Tasks 3-8 rely on when rewriting phase file headers.

- [ ] **Step 1: Update the complexity line**

Old (line 3):
```markdown
**Complexity:** Complex (multi-phase, 37 PRs across 12 phases, client/server architecture)
```

New:
```markdown
**Complexity:** Complex (multi-phase, 13 phases, client/server architecture; exact PR counts vary — see PR Numbering Policy below)
```

- [ ] **Step 2: Replace the Phases table**

Old (lines 57-70):
```markdown
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
```

New:
```markdown
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
```

- [ ] **Step 3: Replace the PR Breakdown table with a historical table + numbering policy**

Old (lines 72-112, the `## PR Breakdown` heading through the last row `| #37 | ... |`):
```markdown
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
```

New:
```markdown
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
```

- [ ] **Step 4: Update the Phase Files list**

Old (lines 114-127):
```markdown
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
```

New:
```markdown
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
```

- [ ] **Step 2: Verify**

Run: `grep -n "Phase 1.5" plans/00-go-t3-tic-tac-toe.md`
Expected: 2 matches (Phases table row, Phase Files list entry)

Run: `grep -n "^| #" plans/00-go-t3-tic-tac-toe.md | wc -l`
Expected: `10` (historical table now has 10 rows: #1, #2, #3, #4–#5, #6, #7–#13, #14, #15, #16, #17 — down from 37, confirming the predicted rows for phase 2+ were removed)

- [ ] **Step 3: Commit**

```bash
git add plans/00-go-t3-tic-tac-toe.md
git commit -m "docs: restructure master plan for Phase 1.5 and CLI parity

Inserts Phase 1.5, drops exact PR-number predictions for Phase 2+ in
favor of a numbering policy, and truncates the PR Breakdown table to
historical (already-merged) entries only.

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 3: Update Phase 1 file — flag CI as deferred to Phase 1.5

**Files:**
- Modify: `plans/phases/01-phase-1-structure.md`

**Interfaces:**
- Consumes: Phase 1.5 file path from Task 1.

- [ ] **Step 1: Update the header block**

Old (lines 4-6):
```markdown
**PRs:** #2–#4
**Release Tag:** v0.2.0 (on PR #4)
**Branch prefix:** phase/1-
```

New:
```markdown
**PRs:** #2–#4 (PR #4 not yet merged — deferred, see Phase 1.5)
**Release Tag:** v0.2.0 (moved to Phase 1.5 — see `phases/01a-phase-1.5-connectivity-cli-foundation.md`)
**Branch prefix:** phase/1-
```

- [ ] **Step 2: Update the PR #4 Post-Phase note**

Old (lines 314-317):
```markdown
### Post-Phase
- Merge PR #4
- Tag `v0.2.0`
- All future PRs will have CI status checks
```

New:
```markdown
### Post-Phase
- **Deferred.** This PR did not merge before PRs #5-#17 landed off-plan.
  See `phases/01a-phase-1.5-connectivity-cli-foundation.md` for current
  status and the merge/tag plan. The file contents specified above (test.yml,
  lint.yml, build.yml, .golangci.yml) are still the correct target — only
  the timing changed.
```

- [ ] **Step 3: Verify**

Run: `grep -n "Phase 1.5" plans/phases/01-phase-1-structure.md`
Expected: 2 matches

- [ ] **Step 4: Commit**

```bash
git add plans/phases/01-phase-1-structure.md
git commit -m "docs: flag Phase 1 CI PR as deferred to Phase 1.5

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 4: Add CLI parity PR to Phase 3 (GUI Frontend)

**Files:**
- Modify: `plans/phases/03-phase-3-gui-frontend.md`

**Interfaces:**
- Consumes: `internal/ui.APIClient` methods `CreateGame() (string, error)`, `GetState(gameID string) (*t3.GameStateResponse, error)`, `MakeMove(gameID, player string, row, col int) (*t3.MoveResponse, error)` — defined in this same phase's PR #9 block.
- Produces: `playCLI(client *ui.APIClient)` in `cmd/client/main.go`, consumed by Task 5's CLI parity PR (Phase 5 extends this function with a `--mode` flag).

- [ ] **Step 1: Update the header**

Old (lines 3-6):
```markdown
**Complexity:** Medium
**PRs:** #9–#11
**Release Tag:** v0.4.0 (on PR #11)
**Branch prefix:** phase/3-
```

New:
```markdown
**Complexity:** Medium
**PRs:** 4 planned (3 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v0.4.0 (on final PR)
**Branch prefix:** phase/3-
```

- [ ] **Step 2: Insert the CLI parity PR block before `### Post-Phase`**

Old (end of file, lines 119-134):
```markdown
### Verification
```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/client
# GUI opens; click cells to play X then O
# Win condition shows "Player X wins!" or "Draw!"
# "New Game" resets the board
```

### Post-Phase
- Merge PR #10
- Tag `v0.4.0`
```

New:
```markdown
### Verification
```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/client
# GUI opens; click cells to play X then O
# Win condition shows "Player X wins!" or "Draw!"
# "New Game" resets the board
```

---

## PR — CLI Parity: Play via Text Client
**Branch:** `phase/3-cli-play`

### Files

#### `cmd/client/main.go` (updated)
Add a `play` subcommand that drives a full 2-player game over stdin/stdout
using the same `internal/ui.APIClient` the GUI uses:
```go
case "play":
    playCLI(client)
```

```go
func playCLI(client *ui.APIClient) {
    gameID, _ := client.CreateGame()
    for {
        state, _ := client.GetState(gameID)
        renderBoard(state.Board)
        if state.Done {
            fmt.Printf("Result: %s\n", resultText(state.Winner))
            return
        }
        fmt.Printf("Turn: %s. Enter move as row,col (0-2): ", state.Current)
        var row, col int
        fmt.Scanf("%d,%d", &row, &col)
        client.MakeMove(gameID, state.Current, row, col)
    }
}

func renderBoard(board [3][3]string) {
    for _, r := range board {
        fmt.Println(strings.Join(r[:], " | "))
    }
}
```

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client play
# Text board renders after each move; alternates X/O prompts
# Win/draw prints "Result: ..." and exits
kill %1
```

### Post-Phase
- Merge final PR (CLI parity)
- Tag `v0.4.0`
```

- [ ] **Step 3: Verify**

Run: `grep -n "CLI Parity" plans/phases/03-phase-3-gui-frontend.md`
Expected: 1 match

Run: `grep -n "Merge PR #10" plans/phases/03-phase-3-gui-frontend.md`
Expected: no matches (stale reference removed)

- [ ] **Step 4: Commit**

```bash
git add plans/phases/03-phase-3-gui-frontend.md
git commit -m "docs: add CLI parity PR to Phase 3 (GUI Frontend)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 5: Add CLI parity PR to Phase 5 (GUI AI Integration)

**Files:**
- Modify: `plans/phases/05-phase-5-gui-ai-integration.md`

**Interfaces:**
- Consumes: `playCLI(client *ui.APIClient)` from Task 4; `CreateGame(mode, computerPlayer string) (string, error)` (updated signature from this phase's PR #15 block).
- Produces: `play --mode` flag on `cmd/client/main.go`, consumed by no later task (auth in Task 6 wraps the whole CLI, not this flag specifically).

- [ ] **Step 1: Update the header**

Old (lines 3-6):
```markdown
**Complexity:** Simple
**PRs:** #15–#16
**Release Tag:** v1.0.0 (on PR #16)
**Branch prefix:** phase/5-
```

New:
```markdown
**Complexity:** Simple
**PRs:** 3 planned (2 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.0.0 (on final PR)
**Branch prefix:** phase/5-
```

- [ ] **Step 2: Insert the CLI parity PR block before `### Post-Phase`**

Old (lines 76-80):
```markdown
# vs Computer, computer plays X:
# Computer immediately places at start (X goes first)
# Human plays O
```

### Post-Phase
- Merge PR #15
- Tag `v1.0.0`
```

New:
```markdown
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
```

- [ ] **Step 3: Verify**

Run: `grep -n "CLI Parity" plans/phases/05-phase-5-gui-ai-integration.md`
Expected: 1 match

- [ ] **Step 4: Commit**

```bash
git add plans/phases/05-phase-5-gui-ai-integration.md
git commit -m "docs: add CLI parity PR to Phase 5 (GUI AI Integration)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 6: Add CLI parity PR to Phase 7 (Auth Frontend)

**Files:**
- Modify: `plans/phases/07-phase-7-auth-frontend.md`

**Interfaces:**
- Consumes: `APIClient.Login(username, password string) error`, `APIClient.Register(username, password string) error`, `APIClient.CurrentUserID() string`, `APIClient.CurrentUsername() string` — defined in this phase's PR #22 block.
- Produces: `login`/`register` subcommands and a `saveToken`/`loadToken` pair in `cmd/client/main.go`, consumed by Task 7's lobby/invite CLI parity PR (needs an authenticated client) and Task 8's matchmaking CLI parity PR.

- [ ] **Step 1: Update the header**

Old (lines 3-6):
```markdown
**Complexity:** Simple
**PRs:** #22–#24
**Release Tag:** v1.2.0 (on PR #24)
**Branch prefix:** phase/7-
```

New:
```markdown
**Complexity:** Simple
**PRs:** 4 planned (3 GUI/backend-removal + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.2.0 (on final PR)
**Branch prefix:** phase/7-
```

- [ ] **Step 2: Insert the CLI parity PR block between PR #23 and PR #24**

Old (lines 116-121):
```markdown
kill %1 %2 %3
```

---

## PR #24 — Remove --no-auth Flag
```

New:
```markdown
kill %1 %2 %3
```

---

## PR — CLI Parity: Auth Subcommands
**Branch:** `phase/7-cli-auth`

### Files

#### `cmd/client/main.go` (updated)
Add `login` and `register` subcommands using the same
`APIClient.Login`/`Register` methods the GUI auth screen calls:
```go
case "login":
    var username, password string
    fmt.Print("Username: ")
    fmt.Scanln(&username)
    fmt.Print("Password: ")
    fmt.Scanln(&password)
    if err := client.Login(username, password); err != nil {
        fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
        os.Exit(1)
    }
    saveToken(client.Token())
case "register":
    // same prompt shape, calls client.Register then saveToken
```

Token persisted to `~/.go-t3/token` via a small `saveToken`/`loadToken` pair
shared by both subcommands; `play`, `lobby`, `queue`, etc. call `loadToken`
on startup to restore the session.

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client register
# Prompts for username/password, registers + logs in, saves token

go run ./cmd/client login
# Prompts for credentials, on success prints "logged in as <username>"
kill %1
```

---

## PR #24 — Remove --no-auth Flag
```

- [ ] **Step 3: Update the final Post-Phase note**

Old (lines 159-161):
```markdown
### Post-Phase
- Merge PR #24
- Tag `v1.2.0`
```

New:
```markdown
### Post-Phase
- Merge final PR (remove --no-auth, after CLI parity lands)
- Tag `v1.2.0`
```

- [ ] **Step 4: Verify**

Run: `grep -n "CLI Parity" plans/phases/07-phase-7-auth-frontend.md`
Expected: 1 match

Run: `grep -n "## PR #24" plans/phases/07-phase-7-auth-frontend.md`
Expected: still present, appearing after the CLI Parity block

- [ ] **Step 5: Commit**

```bash
git add plans/phases/07-phase-7-auth-frontend.md
git commit -m "docs: add CLI parity PR to Phase 7 (Auth Frontend)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 7: Add CLI parity PR to Phase 9 (GUI Opponent Selection)

**Files:**
- Modify: `plans/phases/09-phase-9-gui-opponent-selection.md`

**Interfaces:**
- Consumes: `APIClient.OnlineUsers() ([]t3.OnlineUser, error)`, `InviteUser(opponentUserID string, mode t3.GameMode) (string, error)`, `PendingInvites() ([]t3.PendingInvite, error)`, `AcceptInvite(gameID string) error`, `DeclineInvite(gameID string) error` — defined in this phase's PR #29 block; `loadToken`/authenticated client from Task 6.
- Produces: `lobby`, `invite`, `invites`, `accept`, `decline` subcommands in `cmd/client/main.go`.

- [ ] **Step 1: Update the header**

Old (lines 3-6):
```markdown
**Complexity:** Medium
**PRs:** #29–#31
**Release Tag:** v1.4.0 (on PR #31)
**Branch prefix:** phase/9-
```

New:
```markdown
**Complexity:** Medium
**PRs:** 4 planned (3 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.4.0 (on final PR)
**Branch prefix:** phase/9-
```

- [ ] **Step 2: Insert the CLI parity PR block before `### Post-Phase`**

Old (lines 113-121):
```markdown
# Take turns; each side only clicks their own cells
```

### Post-Phase
- Merge PR #31
- Tag `v1.4.0`
```

New:
```markdown
# Take turns; each side only clicks their own cells
```

---

## PR — CLI Parity: Lobby + Invite Subcommands
**Branch:** `phase/9-cli-lobby`

### Files

#### `cmd/client/main.go` (updated)
```go
case "lobby":
    users, _ := client.OnlineUsers()
    for _, u := range users {
        fmt.Println(u.Username)
    }
case "invite":
    // args[1] = username to invite
    gameID, _ := client.InviteUser(resolveUserID(args[1]), t3.ModeTwoPlayer)
    fmt.Printf("Invited. game_id=%s\n", gameID)
case "invites":
    invites, _ := client.PendingInvites()
    for _, inv := range invites {
        fmt.Printf("%s wants to play (game_id=%s)\n", inv.FromUsername, inv.GameID)
    }
case "accept":
    client.AcceptInvite(args[1]) // args[1] = game_id
case "decline":
    client.DeclineInvite(args[1])
```

### Verification
```bash
go run ./cmd/server &
go run ./cmd/client lobby
# → lists online usernames

go run ./cmd/client invite bob
# → "Invited. game_id=..."

go run ./cmd/client invites
# (as bob) → "alice wants to play (game_id=...)"

go run ./cmd/client accept <game_id>
kill %1
```

### Post-Phase
- Merge final PR (CLI parity)
- Tag `v1.4.0`
```

- [ ] **Step 3: Verify**

Run: `grep -n "CLI Parity" plans/phases/09-phase-9-gui-opponent-selection.md`
Expected: 1 match

- [ ] **Step 4: Commit**

```bash
git add plans/phases/09-phase-9-gui-opponent-selection.md
git commit -m "docs: add CLI parity PR to Phase 9 (GUI Opponent Selection)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 8: Add CLI parity PR to Phase 11 (GUI Random Matchmaking)

**Files:**
- Modify: `plans/phases/11-phase-11-gui-random-matchmaking.md`

**Interfaces:**
- Consumes: `APIClient.JoinMatchmaking() (*t3.MatchmakingResponse, error)`, `LeaveMatchmaking() error`, `MatchmakingStatus() (*t3.MatchmakingResponse, error)` — defined in this phase's PR #35 block; authenticated client from Task 6.
- Produces: `queue join|leave|status` subcommand in `cmd/client/main.go`. Final task — nothing downstream consumes this.

- [ ] **Step 1: Update the header**

Old (lines 3-6):
```markdown
**Complexity:** Simple
**PRs:** #35–#37
**Release Tag:** v1.6.0 (on PR #37)
**Branch prefix:** phase/11-
```

New:
```markdown
**Complexity:** Simple
**PRs:** 4 planned (3 GUI + 1 CLI parity); numbers assigned at merge — see PR Numbering Policy in master plan
**Release Tag:** v1.6.0 (on final PR)
**Branch prefix:** phase/11-
```

- [ ] **Step 2: Insert the CLI parity PR block before `### Post-Phase`**

Old (lines 108-116):
```markdown
# Login → Find Random Opponent → Start → Cancel
# → Returns to mode selector; carol removed from queue
```

### Post-Phase
- Merge PR #37
- Tag `v1.6.0`
```

New:
```markdown
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
```

- [ ] **Step 3: Verify**

Run: `grep -n "CLI Parity" plans/phases/11-phase-11-gui-random-matchmaking.md`
Expected: 1 match

- [ ] **Step 4: Commit**

```bash
git add plans/phases/11-phase-11-gui-random-matchmaking.md
git commit -m "docs: add CLI parity PR to Phase 11 (GUI Random Matchmaking)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>"
```

---

### Task 9: Final consistency pass

**Files:**
- Read-only verification across all files touched by Tasks 1-8.

**Interfaces:**
- Consumes: all files from Tasks 1-8.
- Produces: nothing — this task only verifies, no further edits expected unless a check fails.

- [ ] **Step 1: Confirm no dangling "Merge PR #NN" references remain in the 5 CLI-parity phase files**

Run:
```bash
grep -rn "Merge PR #" plans/phases/03-phase-3-gui-frontend.md plans/phases/05-phase-5-gui-ai-integration.md plans/phases/07-phase-7-auth-frontend.md plans/phases/09-phase-9-gui-opponent-selection.md plans/phases/11-phase-11-gui-random-matchmaking.md
```
Expected: only `plans/phases/07-phase-7-auth-frontend.md` may still reference `Merge PR #24` if Task 6 Step 3 was skipped — otherwise no output. If any unexpected match appears, re-open that file and confirm Post-Phase says "Merge final PR" per its task.

- [ ] **Step 2: Confirm the master plan links resolve to real files**

Run:
```bash
for f in $(grep -oE 'phases/[a-zA-Z0-9._-]+\.md' plans/00-go-t3-tic-tac-toe.md); do test -f "plans/$f" && echo "OK: $f" || echo "MISSING: $f"; done
```
Expected: every line prints `OK:` — no `MISSING:` lines.

- [ ] **Step 3: Confirm the design doc is still referenced from Phase 1.5**

Run: `grep -n "2026-07-09-plan-restructure-cli-connectivity-design.md" plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`
Expected: 1 match

- [ ] **Step 4: No commit needed** — this task is verification-only. If any check in Steps 1-3 fails, fix the offending file and amend that file's own task commit (per Tasks 1-8), not this task.
