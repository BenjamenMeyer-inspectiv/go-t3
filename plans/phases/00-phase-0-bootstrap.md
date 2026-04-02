# Phase 0: Repository Bootstrap

**Complexity:** Simple
**PR:** #1 (single PR — Phase 0 exception)
**Release Tag:** v0.1.0
**Branch:** `phase/0-bootstrap`

## Goal

Initialize the Go module, add a minimal `main.go`, a `Makefile`, and all standard GitHub community files (templates, CODEOWNERS, CONTRIBUTORS, SECURITY, CONTRIBUTING, CHANGELOG). This gives the repository everything needed to start accepting well-structured PRs.

## Context Check
- [x] Repo has: .gitignore, LICENSE, README.md only
- [x] No go.mod yet
- [x] Module name: `github.com/benjamenmeyer/go-t3`

---

## Files to Create

### Go Module

#### `go.mod` (new)
```
module github.com/BenjamenMeyer-inspectiv/go-t3

go 1.25
```

#### `main.go` (new)
```go
package main

import "fmt"

func main() {
    fmt.Println("go-t3: Tic Tac Toe")
}
```

#### `Makefile` (new)
```makefile
.PHONY: build run test fmt lint help

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o bin/go-t3 .

run: ## Run the application
	go run .

test: ## Run tests
	go test ./...

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run
```

#### `README.md` (updated)
Expand from the current one-liner to include: project description, build instructions, usage, and links to CONTRIBUTING and SECURITY.

---

### GitHub Community Files

#### `.github/PULL_REQUEST_TEMPLATE.md` (new)
Standard PR checklist covering: description, type of change, testing done, checklist (tests pass, linter clean, docs updated).

```markdown
## Description
<!-- What does this PR do? Why? -->

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Refactor
- [ ] Documentation
- [ ] Chore / dependency update

## Testing
<!-- How was this tested? -->

## Checklist
- [ ] `go test ./...` passes
- [ ] `golangci-lint run` passes
- [ ] `go build ./...` succeeds
- [ ] Relevant documentation updated
- [ ] PR title follows conventional commit style (`feat:`, `fix:`, `chore:`, etc.)
```

#### `.github/ISSUE_TEMPLATE/config.yml` (new)
Disables blank issues and points users to the templates:
```yaml
blank_issues_enabled: false
contact_links:
  - name: Security Vulnerability
    url: https://github.com/benjamenmeyer/go-t3/security/advisories/new
    about: Please report security vulnerabilities via a private advisory.
```

#### `.github/ISSUE_TEMPLATE/bug_report.md` (new)
```markdown
---
name: Bug Report
about: Report a reproducible bug
title: 'bug: '
labels: bug
assignees: ''
---

## Describe the Bug
<!-- A clear and concise description of what the bug is. -->

## Steps to Reproduce
1.
2.
3.

## Expected Behavior
<!-- What did you expect to happen? -->

## Actual Behavior
<!-- What actually happened? -->

## Environment
- OS:
- Go version (`go version`):
- go-t3 version / commit:

## Additional Context
<!-- Logs, screenshots, etc. -->
```

#### `.github/ISSUE_TEMPLATE/feature_request.md` (new)
```markdown
---
name: Feature Request
about: Suggest a new feature or enhancement
title: 'feat: '
labels: enhancement
assignees: ''
---

## Problem Statement
<!-- What problem does this feature solve? -->

## Proposed Solution
<!-- Describe what you'd like to happen. -->

## Alternatives Considered
<!-- What other approaches did you consider? -->

## Additional Context
<!-- Mockups, related issues, etc. -->
```

#### `.github/CODEOWNERS` (new)
```
# Global owner — all files
* @BenjamenMeyer-inspectiv

# Game logic
/internal/game/ @BenjamenMeyer-inspectiv
/internal/ai/   @BenjamenMeyer-inspectiv

# Server API
/internal/api/  @BenjamenMeyer-inspectiv
/cmd/server/    @BenjamenMeyer-inspectiv

# Client / UI
/internal/ui/   @BenjamenMeyer-inspectiv
/cmd/client/    @BenjamenMeyer-inspectiv

# Auth
/internal/auth/ @BenjamenMeyer-inspectiv

# Shared types
/pkg/           @BenjamenMeyer-inspectiv
```

#### `CONTRIBUTING.md` (new)
Covers: development setup, branch naming (`phase/<n>-<slug>`), PR process, commit message style (conventional commits), running tests and linter, and linking to the issue templates.

#### `CONTRIBUTORS.md` (new)
```markdown
# Contributors

Thanks to everyone who has contributed to go-t3.

| Name | GitHub |
|------|--------|
| Benjamen Meyer | [@benjamenmeyer](https://github.com/benjamenmeyer) |
```

#### `SECURITY.md` (new)
Covers: supported versions table (current `main` branch supported), how to report a vulnerability (private GitHub advisory), response timeline (acknowledge within 3 business days, patch within 30 days), and out-of-scope items.

#### `CHANGELOG.md` (new)
```markdown
# Changelog

All notable changes to go-t3 are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.0] - 2026-04-02
### Added
- Initial repository bootstrap
- Go module, main.go, Makefile
- GitHub community files (PR template, issue templates, CODEOWNERS, SECURITY, CONTRIBUTING)
```

---

## Verification Steps

```bash
go build .
go run .
# Expected: "go-t3: Tic Tac Toe"

# Confirm community files visible on GitHub after push:
# - Repo shows "Contributing" link in sidebar
# - New issue shows bug/feature templates
# - New PR pre-fills with PR template
# - SECURITY tab shows security policy
```

## Post-Phase

- Create PR #1 from `phase/0-bootstrap` → `main`
- Merge PR
- Tag `v0.1.0`
