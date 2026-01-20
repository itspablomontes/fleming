---
trigger: glob
globs: *.go
---

# Go Code Rules

We write **simple, readable, boring** Go that any experienced Gopher can understand quickly.

## Formatting & Tooling (non-negotiable)
- **Go Version**: 1.25+ (Matches `go.mod` `1.25.5`)
- `go fmt ./...` on every file — always.
- Use `goimports` or `gci` to automatically sort + organize imports.
- Enforce via pre-commit or CI: `golangci-lint run --fix`
- Recommended linter presets: default + `gosimple`, `unused`, `revive`, `errcheck`, `nilerr`, `goheader` (for license/copyright if needed)
- Editor integration: VS Code / GoLand with `gopls` + golangci-lint plugin

## Naming Conventions
- `camelCase` for variables, functions, method receivers, local vars
- `PascalCase` for exported types, constants, package-level vars
- Short is good when scope is small: `db`, `cfg`, `cli`, `srv`, `h`, `w`, `r`
- Acronyms stay **UPPERCASE**: `HTTPRequest`, `JWTToken`, **not** `HttpRequest`
- No Hungarian notation, no `strUserName`, no `pInterface`
- Avoid `GetXxx()` prefix on simple getters unless they are expensive

## Package & Module Structure
- **Monorepo Layout**:
    - Backend API lives in `apps/backend`
    - Shared libraries in `pkg/` (at root)
- One package = one clear responsibility
- `main` package stays tiny — move real logic to `cmd/...` or `internal/...`
- Use `internal/` for code that should never be imported outside the module
- **Preferred Folder Pattern** (Domain-first):
  - `internal/domain/user/`
  - `internal/application/`
  - `internal/infrastructure/`
- No cyclic dependencies — use `golangci-lint`’s `depguard` or `gocyclo` if needed

## HTTP & API (Gin Framework)
- Project uses `github.com/gin-gonic/gin` (v1.11.0+)
- Handlers should be thin adapters between HTTP and your domain logic
- **Always** return appropriate status codes
- Do not leak internal error details to the client

## Error Handling
- **Always** check errors — no `_ =` or `//nolint:errcheck` unless truly justified
- Wrap errors with context: `fmt.Errorf("cannot fetch user %d: %w", id, err)`
- Use `errors.Is` / `errors.As` / `errors.Join` when inspecting wrapped errors
- Sentinel errors only for a very small set of public API cases

## Context & Cancellation
- `context.Context` is **first parameter** in almost every function that can block or spawn work (DB calls, HTTP requests)
- Pass `ctx` through the entire call chain
- Use `context.WithTimeout`, `context.WithCancel` for controlling request lifecycles
- **Never** start naked `go func()` without a way to wait / cancel / collect errors (use `errgroup`)

## Concurrency Idioms
- Prefer channels + `select` for coordination when it improves clarity
- Use mutexes when protecting shared mutable state is simplest
- `sync.WaitGroup` + `errgroup` > manual channel boilerplate in most cases
- Keep goroutine lifetimes short and obvious

## Performance & Allocation Awareness
- Prefer `[]T` over `[N]T` unless fixed size is semantically important
- Use `strings.Builder` / `bytes.Buffer` for concatenation in loops
- Be conscious of escape analysis — but **don’t** pre-optimize without profiling
- Use `defer` for cleanup (`Unlock`, `Close`) ensuring resource safety

## Testing Philosophy
- Table-driven tests are the default for pure logic
- Use `testing` stdlib first — `testify` only when it adds real clarity
- Test files live next to production code: `foo.go` → `foo_test.go`
- Cover happy path + important errors + edge cases
- Goal: High coverage on business logic (`internal/domain`), relaxed on boilerplate

## Logging & Observability
- Use `log/slog` (structured logging) — **never** `fmt.Printf` or `log.Printf` in prod code
- Default handler: JSON in production, text in development
- Log levels: `Debug` (disabled in prod by default), `Info`, `Warn`, `Error`
- Always include `err`, correlation IDs, user IDs, request path when relevant

## Quick Mantra
**"Clear is better than clever. Simple is better than fast (until proven otherwise)."**
