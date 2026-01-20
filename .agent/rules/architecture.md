---
trigger: model_decision
description: Anytime you need to create a new file or folder, in plan or build.
---

# Fleming Architecture Rules – 2026 Edition
(Monorepo, Modular Monolith, TypeScript Frontend)

We follow a **strict strict monorepo structure** designed for modularity, data sovereignty, and maintainability.

## 1. High-Level Monorepo Structure
- **apps/**: Deployable applications.
  - `apps/backend`: The core Go API (Modular Monolith).
  - `apps/web`: The React Frontend (SPA).
- **pkg/**: Shared Go libraries *only* (intended for potential external use or sharing between multiple apps).
  - internal code must stay in `apps/backend/internal`.
- **contracts/**: Blockchain smart contracts (Solidity).
- **infra/**: Infrastructure-as-Code (Terraform, Docker configs).
- **docs/**: Architectural decisions (ADRs) and documentation.

## 2. Backend Architecture (`apps/backend`)
**Pattern:** Modular Monolith with Domain-Driven Design (DDD) principles.

### Directory Layout
- **cmd/**: Entry points only (`main.go`). No business logic.
- **internal/**: Private application code.
  - `internal/auth/`, `internal/timeline/`: **Domain Modules**.
  - Each module should generally contain its own:
    - `models.go` (Domain entities)
    - `service.go` (Business logic)
    - `repository.go` (Data access)
    - `handler.go` (HTTP transport, if needed)
- **router.go**: Central HTTP routing definition (Gin).

### Rules
- **Dependency Rule**: `cmd` depends on `internal/modules`. Modules should be loosely coupled.
- **No Cyclic Dependencies**: Use interfaces to break cycles between modules.
- **Transport Agnostic**: Business logic (`service`) must not know about HTTP (`gin.Context`). Handlers translate HTTP <-> Domain.

## 3. Frontend Architecture (`apps/web`)
**Pattern:** Feature-based React SPA.

### Directory Layout
- **src/routes/**: File-based routing (TanStack Router).
  - Map 1:1 to URL structure.
- **src/features/**: Vertical slices of functionality.
  - `src/features/auth/`: Components, hooks, and queries specific to Auth.
  - `src/features/timeline/`: Timeline visualization logic.
- **src/components/ui/**: Shared, dumb UI components (Buttons, Inputs).
- **src/lib/**: formatting, API clients, shared utilities.

### Rules
- **Co-location**: Keep related things close. A hook used only by the "Timeline" feature belongs in `src/features/timeline/hooks`.
- **Global vs Local**: Only truly generic components go to `src/components`.
- **State**: Prefer server state (React Query/TanStack Query) over global client state (Zustand/Context). Use Context only for dependency injection or strictly UI state (theme).

## 4. Cross-Cutting Concerns
- **API Communication**: REST JSON (proto-compatible if needed).
- **Authentication**: Wallet-based (SIWE) + JWT.
- **Database**: Postgres (via `pgx` or equivalent). Migrations in `migrations/`.

## 5. Architectural Quality Attributes
- **Self-Hostability**: Everything must run via `docker compose up`.
- **Privacy First**: User data is encrypted or sovereign.
- **Simplicity**: Do not introduce microservices until the monolith is proven too large (unlikely).
