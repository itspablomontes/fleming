---
trigger: glob
description: Go backend development rules.
globs: "*.go"
---

# Go Code Rules – Fleming Backend

> **See also**: general.md for project philosophy and engineering principles.

We write **simple, readable, boring** Go that any experienced Gopher can understand quickly.

---

## 1. Formatting & Tooling (Non-Negotiable)

| Tool | Command | Purpose |
|------|---------|---------|
| **Go Version** | 1.25+ | Match `go.mod` |
| **Format** | `go fmt ./...` | Always run before commit |
| **Imports** | `goimports` or `gci` | Auto-sort imports |
| **Lint** | `golangci-lint run --fix` | Pre-commit or CI |

### Recommended Linters
```yaml
# .golangci.yml
linters:
  enable:
    - gosimple
    - unused
    - revive
    - errcheck
    - nilerr
```

---

## 2. Naming Conventions

| Scope | Convention | Example |
|-------|------------|---------|
| Variables, functions | `camelCase` | `userID`, `fetchData` |
| Exported types | `PascalCase` | `TimelineEvent`, `AuthService` |
| Acronyms | UPPERCASE | `HTTPRequest`, `JWTToken` |
| Short vars (small scope) | Short | `db`, `cfg`, `h`, `w`, `r` |

### Anti-Patterns
- ❌ Hungarian notation: `strUserName`, `pInterface`
- ❌ Unnecessary prefixes: `GetUser()` when `User()` suffices

---

## 3. Package Structure

```
apps/backend/
├── cmd/                  # Entry points only
│   └── fleming/main.go
├── internal/             # Feature modules
│   ├── auth/             # Auth feature (handler, service, repo, entity)
│   ├── timeline/         # Timeline feature (handler, service, repo, entity)
│   └── middleware/       # Shared middleware
├── router.go             # Central routing
└── Dockerfile
```

### Rules
- **Package by Feature**: Group `handler.go`, `service.go`, `repository.go`, and `entity.go` within the feature folder (e.g., `internal/auth/`).
- **No global `handlers` or `repositories` packages**.
- **`main` package stays tiny** — real logic in `internal/`.
- **Use `internal/`** for code that shouldn't be imported externally.

---

## 4. GORM & Repository Pattern

### Repository Interface
Define interfaces for data access in `repository.go` within the feature package.

```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}
```

### GORM Implementation
Implement the interface using GORM.

```go
type GormRepository struct {
    db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
    return &GormRepository{db: db}
}
```

### Rules
- **Use GORM** for database interactions.
- **Decouple Service from GORM**: Services should depend on the `Repository` interface, not `*gorm.DB` directly.
- **Entities** should have GORM tags in `entity.go`.

---

## 4. HTTP & API (Gin Framework)

### Handler Pattern
```go
// Handlers are thin adapters: HTTP ↔ Domain
func (h *TimelineHandler) HandleGetTimeline(c *gin.Context) {
    // Extract from HTTP context
    userID := c.GetString("user_id")
    
    // Call domain service
    events, err := h.service.GetTimeline(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
        return
    }
    
    // Return HTTP response
    c.JSON(http.StatusOK, gin.H{"events": events})
}
```

### Rules
- **Handlers don't contain business logic**.
- **Always return appropriate status codes**.
- **Never leak internal error details to clients**.

---

## 5. Error Handling

```go
// ✅ Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("fetch user %d: %w", id, err)
}

// ✅ Good: Check errors with Is/As
if errors.Is(err, sql.ErrNoRows) {
    return ErrNotFound
}
```

### Rules
- **Always check errors** — no `_ =` unless justified.
- **Wrap with context** — use `fmt.Errorf("...: %w", err)`.
- **Sentinel errors** — only for public API boundaries.

---

## 6. Context & Cancellation

```go
// ✅ Good: Context as first parameter
func (s *Service) GetTimeline(ctx context.Context, userID string) ([]Event, error) {
    return s.repo.FindByUser(ctx, userID)
}

// ✅ Good: Respect cancellation
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result, nil
}
```

### Rules
- **`context.Context` is first parameter** for blocking operations.
- **Pass context through the chain**.
- **Use `errgroup`** for concurrent operations.
- **Never start naked `go func()`** without cancellation.

---

## 7. Concurrency

| Pattern | Use When |
|---------|----------|
| Channels + select | Coordination improves clarity |
| Mutexes | Protecting shared mutable state |
| `sync.WaitGroup` | Waiting for goroutines |
| `errgroup` | Concurrent operations with error handling |

---

## 8. Testing

### Philosophy
- **Table-driven tests** are the default.
- **Use `testing` stdlib** — `testify` only when it adds clarity.
- **Test files** live next to production code: `foo.go` → `foo_test.go`.

### Coverage Goals
| Layer | Coverage |
|-------|----------|
| Domain logic (`internal/domain`) | 80%+ |
| Handlers | 60%+ |
| Boilerplate | Relaxed |

```go
func TestService_GetTimeline(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        want    []Event
        wantErr bool
    }{
        {"happy path", "user1", []Event{{ID: "1"}}, false},
        {"user not found", "unknown", nil, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

---

## 9. Logging & Observability

```go
// ✅ Good: Structured logging with slog
slog.Info("user logged in",
    "user_id", userID,
    "ip", clientIP,
)

slog.Error("failed to fetch timeline",
    "err", err,
    "user_id", userID,
)
```

### Rules
- **Use `log/slog`** — not `fmt.Printf` or `log.Printf`.
- **JSON in production**, text in development.
- **Always include**: `err`, correlation IDs, user IDs, request path.

---

## Quick Mantra

> **"Clear is better than clever. Simple is better than fast (until proven otherwise)."**