# Fleming Architecture

> Modular Monolith. Simple, not Easy. Private by Default.

## 1. Monorepo Structure
- `apps/backend/`: Go API (Modular Monolith)
- `apps/web/`: React SPA (Feature-based)
- `pkg/protocol/`: Canonical truth for types, schemas, and crypto.
- `contracts/`: Blockchain smart contracts.
- `infra/`: Docker and IaC.

## 2. Backend Design (`apps/backend`)
Follow **Modular Monolith** with Domain-Driven Design (DDD).
- **Module Structure**: `entity.go` (models), `service.go` (logic), `repository.go` (data), `handler.go` (HTTP).
- **Dependency Rule**: `cmd` → `internal`. Services MUST NOT know about HTTP (`gin.Context`).
- **No Cycles**: Use interfaces to break cyclic dependencies.

## 3. Frontend Design (`apps/web`)
Follow **Feature-based** architecture.
- **Location**: `src/features/{feature}/` for components, hooks, and logic.
- **Thin Routes**: `src/routes/` are config only; import Page components from `features/`.
- **State**: Server state (React Query) > Global state. Context for DI.

## 4. Quality Attributes
- **Self-Hostable**: Must run via `docker compose up`.
- **Privacy**: Encryption at rest; minified PII storage.
- **Scale**: Modular monolith first. No microservices without proof of necessity.

## 5. Domain Decision Tree
1. Belongs to existing domain? → Add to module.
2. Truly shared? → `pkg/protocol` (types) or `internal/common` (logic).
3. New domain? → Create new module in `internal/`.