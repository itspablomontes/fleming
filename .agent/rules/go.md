---
trigger: glob
description: Go backend development rules.
globs: "*.go"
---

# Go Coding Guidelines – Fleming Backend

> **See also**: `general.md` for project philosophy and engineering principles.  
> We prioritize **simple, readable, and maintainable** Go code that any experienced Go developer can understand at a glance. Favor clarity over cleverness, explicitness over implicitness, and correctness over premature optimization.  
> We use **Go 1.25+** to leverage modern language features, runtime improvements, and performance enhancements.

---

## 1. Tooling & Formatting (Non-Negotiable)

| Tool/Command              | Purpose                                    | Notes                                               |
| ------------------------- | ------------------------------------------ | --------------------------------------------------- |
| **Go Version**            | `1.25+` (pin to latest stable in `go.mod`) | Use `go.mod` `go` directive to enforce              |
| **Formatting**            | `gofumpt -w .` or `go fmt ./...`           | Stricter than `go fmt`; run pre-commit              |
| **Imports**               | `goimports -w .` or `gci write .`          | Group and sort imports (std, third-party, internal) |
| **Linting**               | `golangci-lint run --fix`                  | Run in pre-commit and CI                            |
| **Vet & Static Analysis** | `go vet ./...`                             | Always include in CI                                |

### Recommended `.golangci.yml` (minimal but effective set)

```yaml
linters:
  enable:
    - errcheck      # Checks for unchecked errors
    - gosimple      # Suggests simplifications (included in staticcheck)
    - staticcheck   # Comprehensive static analysis (SA, S, ST series)
    - unused        # Detects unused code
    - ineffassign   # Detects ineffective assignments
    - bodyclose     # Ensures HTTP response bodies are closed
    - nilerr        # Catches returning nil error with non-nil result
    - revive        # Drop-in replacement for golint with better rules
    - gofumpt       # Enforces stricter formatting
    - misspell      # Catches common spelling mistakes

issues:
  exclude-use-default: false
```

---

## 2. Naming Conventions

Follow standard Go naming (Effective Go).

| Scope                  | Convention   | Example                              |
| ---------------------- | ------------ | ------------------------------------ |
| Variables, functions   | `camelCase`  | `userID`, `fetchTimeline`            |
| Exported types/methods | `PascalCase` | `TimelineEvent`, `AuthService`       |
| Interfaces             | Descriptive  | `Repository`, not `UserRepositoryer` |
| Acronyms               | ALL CAPS     | `HTTPClient`, `JWTToken`             |
| Short-lived variables  | Short        | `r`, `w`, `ctx`, `db`, `mu`          |
| Constants              | `PascalCase` | `DefaultTimeout`, `maxRetries`       |

### Anti-Patterns
- ❌ Hungarian notation (`strName`, `iCount`)
- ❌ Unnecessary prefixes (`GetUser()` → `User()` is sufficient)
- ❌ Abbreviations unless widely known (`cfg` is fine, `usr` is not)

---

## 3. Project & Package Structure

Organize by feature / vertical slice.

```text
apps/backend/
├── cmd/
│   └── fleming/
│       └── main.go                  # Minimal: wire/bootstrap only
├── internal/
│   ├── auth/                        # Feature package
│   │   ├── entity.go                # Domain models (with GORM tags)
│   │   ├── repository.go            # Interface + impl
│   │   ├── service.go               # Business logic
│   │   ├── handler.go               # HTTP adapters
│   │   └── middleware.go            # Feature-specific middleware
│   ├── timeline/
│   │   └── (same structure)
│   ├── common/                      # Shared utilities (errors, types, helpers)
│   ├── config/                      # Configuration loading & structs
│   └── middleware/                  # Cross-cutting middleware (logging, auth, recovery)
├── pkg/                             # Optional: reusable public packages (rare)
├── router.go                        # Central Gin router setup
├── go.mod / go.sum
└── Dockerfile
```

### Key Rules
- `internal/` for all non-public code.
- No god packages like `handlers/`, `repositories/`, `models/`.
- `main.go` stays tiny — only wiring and server startup.
- Shared cross-feature code goes in `internal/common/` or dedicated feature packages.

---

## 4. Dependency Injection & Bootstrapping

- Use explicit constructor functions (`NewXxx(deps...)`).
- Avoid global variables or singletons.

**Example:**
```go
func NewTimelineHandler(service TimelineService, middleware ...gin.HandlerFunc) *TimelineHandler {
    return &TimelineHandler{service: service}
}
```

---

## 5. Database Access (GORM + Repository Pattern)

```go
// repository.go
type Repository interface {
    Create(ctx context.Context, e *Event) error
    FindByUserID(ctx context.Context, userID string) ([]Event, error)
}

type gormRepository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &gormRepository{db: db}
}
```

### Rules
- Services depend on the interface, not `*gorm.DB`.
- Entities live in `entity.go` with GORM tags.
- Use scoped sessions (`db.WithContext(ctx)`) to respect cancellation.

---

## 6. HTTP Layer (Gin)

Handlers are thin adapters — validate input, call service, translate output/errors.

```go
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
    userID := c.GetString("user_id")

    events, err := h.service.GetTimeline(c.Request.Context(), userID)
    if err != nil {
        handleError(c, err) // Centralized error response
        return
    }

    c.JSON(http.StatusOK, gin.H{"events": events})
}
```

### Rules
- No business logic in handlers.
- Use structured response types and centralized error handling.
- Validate request bodies with `c.ShouldBindJSON(&req)` + `validator.v10` tags.

---

## 7. Error Handling

```go
// Wrap with context
return fmt.Errorf("fetch timeline for user %s: %w", userID, err)

// Sentinel errors for public boundaries
var ErrNotFound = errors.New("not found")

// Checking
if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound) {
    return ErrNotFound
}
```

### Rules
- Never ignore errors except in rare, documented cases.
- Always wrap with `%w` for context.
- Use `errors.Is/As` and `errors.Join` (Go 1.20+).

---

## 8. Context & Cancellation

- `context.Context` should be the first parameter for any potentially blocking call.
- Propagate through entire call chain.
- Always respect `ctx.Done()` in long operations.
- Use `golang.org/x/sync/errgroup` for concurrent work.

---

## 9. Concurrency

| Pattern              | Preferred Use Case                        |
| -------------------- | ----------------------------------------- |
| `errgroup.Group`     | Parallel tasks with error propagation     |
| `sync.Mutex/RWMutex` | Protecting shared state                   |
| Channels + `select`  | When coordination adds clarity            |
| `sync.WaitGroup`     | Simple waiting (no error handling needed) |

- Never spawn naked `go func()` — always tie lifetime to context or a manager.

---

## 10. Testing & Benchmarking

- Table-driven tests as default.
- Prefer stdlib testing; use `testify` only for complex asserts.
- Use `t.Parallel()` for independent subtests.
- Use `t.Cleanup()` for resource teardown.
- Mock repositories with interfaces (no mocking frameworks needed).

### Coverage Targets
- Domain/services: 80%+
- Handlers: 60%+
- Boilerplate: pragmatic

---

## 11. Logging & Observability

- Use `log/slog` exclusively.

```go
slog.InfoContext(ctx, "timeline fetched",
    "user_id", userID,
    "event_count", len(events),
    "duration_ms", duration.Milliseconds(),
)

slog.ErrorContext(ctx, "failed to fetch timeline",
    "err", err,
    "user_id", userID,
)
```

### Rules
- Always use structured fields.
- Include correlation/request ID (via middleware).
- Production: JSON output; development: human-readable.

---

## 12. Performance & Resource Efficiency

- **Optimize only with evidence**: Always profile (`go tool pprof`) and benchmark before changing code.
- **Minimize allocations**: Reduce garbage collection pressure in hot paths.
- **Reuse objects**: Use `sync.Pool` for frequent, short-lived objects.
- **Choose efficient data structures**: Favor slices for cache locality; use maps judiciously.

---

## 13. Documentation & Comments

Code should be self-documenting.

- **Godoc comments required for exported types, functions, and constants.**
- **Avoid comments explaining what the code does** — use clear names instead.
- **Allow comments only for**:
    - Why a non-obvious decision was made.
    - Complex algorithms.
    - Workarounds for external bugs (with issue links).
- **Delete outdated comments immediately.**

---

## Guiding Mantra

> **"Clear is better than clever. Simple is better than fast — until profiling proves otherwise."**
> Write code that reads like prose, scales through understandability, and evolves gracefully.