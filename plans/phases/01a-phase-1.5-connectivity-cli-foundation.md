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
