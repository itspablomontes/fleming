---
trigger: model_decision
description: Anytime you need to create a new file or folder, make architectural decisions, or design system structure.
---

# Fleming Architecture Rules
Modular Monolith

> **See also**: general.md for project philosophy and engineering principles.

We follow a **strict monorepo structure** designed for modularity, data sovereignty, and maintainability.

---

## 1. High-Level Monorepo Structure

```
fleming/
├── apps/
│   ├── backend/          # Go API (Modular Monolith)
│   └── web/              # React Frontend (SPA)
├── pkg/                  # Shared Go libraries (external use)
├── contracts/            # Blockchain smart contracts (Solidity)
├── infra/                # Infrastructure-as-Code (Terraform, Docker)
├── docs/                 # ADRs and documentation
└── .agent/rules/         # AI agent rules
```

### Directory Responsibilities
| Directory | Purpose | Notes |
|-----------|---------|-------|
| `apps/backend` | Core Go API | All business logic here |
| `apps/web` | React SPA | Feature-based architecture |
| `pkg/` | Shared Go packages | Only if reusable across apps |
| `contracts/` | Smart contracts | Solidity, Foundry |
| `infra/` | IaC | Terraform, Docker Compose |

---

## 2. Backend Architecture (`apps/backend`)

**Pattern**: Modular Monolith with Domain-Driven Design (DDD) principles.

### Directory Layout
```
apps/backend/
├── cmd/                  # Entry points (main.go only)
├── internal/             # Private application code
│   ├── auth/             # Domain Module
│   ├── timeline/         # Domain Module
│   └── handlers/         # HTTP handlers
├── router.go             # Central routing
└── Dockerfile
```

### Domain Module Structure
Each module (`internal/auth/`, `internal/timeline/`, etc.) should contain:
- `models.go` — Domain entities and types
- `service.go` — Business logic (transport-agnostic)
- `repository.go` — Data access layer
- `handler.go` — HTTP handlers (optional, can be in `handlers/`)

### Architectural Rules
| Rule | Description |
|------|-------------|
| **Dependency Rule** | `cmd` → `internal`. Modules should be loosely coupled. |
| **No Cyclic Dependencies** | Use interfaces to break cycles. |
| **Transport Agnostic** | Services must NOT know about HTTP (`gin.Context`). |
| **Handler = Adapter** | Handlers translate HTTP ↔ Domain. |

---

## 3. Frontend Architecture (`apps/web`)

**Pattern**: Feature-based React SPA with TanStack Router.

### Directory Layout
```
apps/web/src/
├── routes/               # File-based routing (TanStack)
├── features/             # Vertical feature slices
│   ├── auth/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── pages/
│   │   └── types/
│   └── timeline/
├── components/           # Shared UI components
│   ├── ui/               # Primitives (Button, Input)
│   └── common/           # Generic (Logo, ErrorBoundary)
├── lib/                  # Utilities, API clients
├── hooks/                # Global hooks
└── types/                # Global types (named files, not index.ts)
```

### Architectural Rules
| Rule | Description |
|------|-------------|
| **Co-location** | Keep related code close. Feature-specific → `features/{feature}/` |
| **Thin Routes** | Route files define routing only, import Page from `features/` |
| **Server State First** | Prefer React Query over global state. |
| **Context for DI** | Use React Context for dependency injection, not global state. |

---

## 4. Cross-Cutting Concerns

| Concern | Approach |
|---------|----------|
| **API** | REST + JSON (protobuf-compatible if needed) |
| **Auth** | Wallet-based (SIWE) + JWT |
| **Database** | PostgreSQL via `pgx` |
| **Migrations** | `migrations/` directory, versioned SQL |
| **Storage** | MinIO (S3-compatible) |

---

## 5. Architectural Quality Attributes

### Self-Hostability
- Everything runs with `docker compose up`.
- No mandatory external services.

### Privacy First
- User data encrypted at rest.
- No telemetry without consent.

### Simplicity
- No microservices until proven necessary.
- Modular monolith scales further than you think.

---

## 6. Decision-Making Framework

Before adding new components, ask:

1. **Does this belong in an existing module?** → Add it there.
2. **Is this truly shared across modules?** → Put in `/lib` or `/pkg`.
3. **Is this a new domain concept?** → Create a new module in `internal/`.
4. **Is this generic infrastructure?** → Put in `infra/`.

---

## Quick Mantra

> **"Modular, not Micro. Simple, not Easy. Private by Default."**