# Design: Fold Connectivity Checks + CLI Client into the Plan

**Date:** 2026-07-09
**Status:** Approved

## Background

The original plan (`plans/00-go-t3-tic-tac-toe.md`) specifies 12 phases / 37 PRs
delivering a tic-tac-toe game with a fyne GUI client. Reality diverged after
Phase 1 (PRs #2-#3 merged): PRs #4-#17 merged off-plan, adding:

- `/ping`, `/pong`, `/ping/reset` endpoints in `internal/api/api.go`
- A flag-based CLI in `cmd/client/main.go` (`ping`, `pong`, `shell` subcommands)
- A banner, a README tweak, a dependabot version bump

These were built to generate test cases for another task but provide real
product value: ping/pong validates network connectivity, and the CLI shell
gives a non-GUI way to interact with the server. This design folds both into
the plan going forward instead of treating them as drift to revert.

There is also an unmerged branch, `phase/1-github-actions`, containing a prior
attempt at the Phase 1 CI PR. That attempt **strips out** ping/pong and reverts
the CLI back to a bare `ui.Run()` call — the opposite of this design's
direction. It has no PR and predates PR #16/#17. It is left alone (not
deleted) but marked superseded — **do not merge as-is**.

## Decisions

1. **Ping/pong → permanent connectivity-check feature.** Not dev scaffolding;
   part of the shipped product surface both clients use before gameplay calls.
2. **CLI → first-class second frontend, full parity, phase-by-phase.** Not a
   thin smoke-test tool. Every future GUI-facing phase gets a matching CLI PR.
3. **Retro-fit via bridging phase (1.5), not a Phase 1 rewrite.** Phase 1's
   original scope (server/client skeleton) stays as originally written;
   PRs #4-#17 get documented under a new Phase 1.5 with real purpose.
4. **CI (original Phase 1 PR #4 concept) is deferred, not abandoned.** It's
   called out explicitly in Phase 1.5 as an outstanding item so it doesn't
   silently get lost a second time.
5. **Stop predicting exact GitHub PR numbers in phase docs.** Numbers already
   drifted once (predicted #4 = CI workflows; actual #4 = dependabot patch).
   Phase files describe PR purpose/order only; real numbers get filled in
   retroactively after merge, in an "Actual PRs" line.
6. **CLI parity lands via Approach A: one extra PR per existing GUI phase**,
   not dedicated CLI-only phases (Approach B) or one catch-up phase at the end
   (Approach C). Keeps CLI testable throughout, avoids roadmap bloat, matches
   the original ping/pong smoke-test rationale.

## Phase 1.5: Connectivity + CLI Foundation

New file: `plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`

Documents already-merged work under real purpose:

- **Connectivity check** — `/health`, `/ping`, `/pong`, `/ping/reset` in
  `internal/api/api.go`, formalized as the surface both clients probe before
  attempting gameplay calls.
- **CLI foundation** — `cmd/client/main.go`'s flag-based dispatch (ping/pong/
  shell) declared as the base for the headless client; all later CLI-parity
  PRs build on this entrypoint.
- Notes on the superseded `phase/1-github-actions` branch.
- CI explicitly flagged as outstanding (deferred per decision, tracked here
  so it surfaces again when this phase is picked up).
- Tag `v0.2.0` lands once this phase's outstanding items (at minimum CI) are
  resolved.

## PR numbering policy (going forward from Phase 1.5)

- Master plan's PR Breakdown table: drop the `#NN` column, use ordered
  sequencing per phase (PR A, PR B, ...) instead of predicted numbers.
- Phase file headers: `**PRs:** #22-24` becomes `**PRs:** 3 planned
  (description, description, description)`.
- After a phase's PRs actually merge, the phase file gains a short
  `**Actual PRs:** #22, #23, #24` line. This is the only place real numbers
  appear, and only in hindsight.
- Phases 0-1 are unaffected (already resolved with real numbers).

## CLI parity by phase (Approach A)

Backend phases (2, 4, 6, 8, 10) are unaffected — the API they build serves
both clients already. Each GUI-facing phase gains one CLI-parity PR:

| Phase | Existing GUI PR(s) | + CLI parity PR |
|---|---|---|
| 3 — GUI Frontend | api-client, board-widget, game-window | CLI: connect + play via stdin (text board render, move input) |
| 5 — GUI AI Integration | board cell control, mode selector | CLI: `--mode=ai` flag, difficulty select via stdin prompt |
| 7 — Auth Frontend | auth screen | CLI: `login`/`register` subcommands, token stored same location as GUI |
| 9 — GUI Opponent Selection | lobby widget, invite flow | CLI: `lobby` subcommand (list online users), `invite`/`accept`/`decline` subcommands |
| 11 — GUI Random Matchmaking | matchmaking widget | CLI: `queue join`/`queue leave`/`queue status` subcommands |

CLI PRs reuse the shared `APIClient` (from `internal/ui/client.go`) rather
than duplicating HTTP logic. CLI-specific rendering/parsing lives in
`cmd/client/main.go`, or a new `internal/cli/` package once subcommand count
grows past a handful.

## Files impacted by this design

**New:**
- `plans/phases/01a-phase-1.5-connectivity-cli-foundation.md`

**Edited:**
- `plans/00-go-t3-tic-tac-toe.md` — add Phase 1.5 to the Phases table and
  Phase Files list; switch PR Breakdown table to non-numeric sequencing.
- `plans/phases/01-phase-1-structure.md` — note CI as outstanding, not
  abandoned.
- `plans/phases/03-phase-3-gui-frontend.md`
- `plans/phases/05-phase-5-gui-ai-integration.md`
- `plans/phases/07-phase-7-auth-frontend.md`
- `plans/phases/09-phase-9-gui-opponent-selection.md`
- `plans/phases/11-phase-11-gui-random-matchmaking.md`
  — each gains one CLI-parity PR entry per the table above.

No source code changes in this pass — this is a plan/doc restructuring only.
Actual CLI parity code is written when each phase is executed.

## Verification

```bash
grep -rn "PR #4" plans/                       # CI reference still points to the right place
grep -rn "phase/1-github-actions" plans/      # superseded-branch note present
```

Manual re-read of the master plan, Phase 1.5, and the five edited phase files
to confirm no contradictions (numbering policy applied consistently, CLI
parity rows present, CI gap flagged exactly once).
