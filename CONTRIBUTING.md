# Contributing to go-t3

Thank you for your interest in contributing!

## Development Setup

1. Install Go 1.25 or later: https://go.dev/dl/
2. Install golangci-lint: https://golangci-lint.run/usage/install/
3. Clone the repository and install dependencies:

```bash
git clone https://github.com/BenjamenMeyer-inspectiv/go-t3.git
cd go-t3
go mod tidy
```

## Branch Naming

Branches follow the pattern `phase/<n>-<short-description>`, for example:

```
phase/2-game-logic
phase/6-auth-middleware
```

For bug fixes outside a phase: `fix/<short-description>`
For chores: `chore/<short-description>`

## Commit Messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add minimax AI opponent
fix: correct win detection on diagonal
chore: update golangci-lint to v2
docs: add API endpoint examples to README
```

## Pull Request Process

1. Branch from `main`
2. Make focused, atomic changes
3. Ensure all checks pass locally:
   ```bash
   make test
   make lint
   go build ./...
   ```
4. Open a PR — the template will guide you through the checklist
5. A reviewer (see CODEOWNERS) must approve before merging

## Reporting Bugs

Use the [Bug Report](.github/ISSUE_TEMPLATE/bug_report.md) issue template.

## Requesting Features

Use the [Feature Request](.github/ISSUE_TEMPLATE/feature_request.md) issue template.

## Security Issues

Do **not** open a public issue for security vulnerabilities.
See [SECURITY.md](SECURITY.md) for how to report them privately.
