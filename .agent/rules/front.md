---
trigger: glob
globs: *.ts , *.tsx
---

# React + TypeScript Rules – 2026 Edition (Fleming Codebase)
(React 19, Vite + Rolldown, TanStack Stack, Biome)

We write **clean, maintainable, type-safe, SOLID-compliant** React code.
We lean heavily on **React 19** primitives, **TanStack** libraries (Router, Form), and **Biome** for strict enforcement.

## Tooling & Language Discipline
- **Build**: Vite 7 (Rolldown) + `@vitejs/plugin-react-swc`.
- **Lint/Format**: **Biome** (`@biomejs/biome`). No ESLint/Prettier.
  - Run `pnpm lint` (or `biome check --write` to fix).
  - Enforces import sorting, formatting, and best practices automatically.
- **Styling**: Tailwind CSS v4 (`@tailwindcss/vite`) + `clsx` + `tailwind-merge` + `cva`.

## TypeScript Standards & Best Practices
- **Strict Typing**: `strict: true` is non-negotiable.
- **No `any`**: Strictly forbidden. Use `unknown` with narrowing if necessary.
- **Explicit Returns**: Exported functions and hooks should have explicit return types to prevent accidental API leaks.
- **Extract Union Types**: Do not inline union types in interfaces setup. Always extract them to a named type.
  - **Bad**: `interface User { role: "admin" | "user"; }`
  - **Good**: `interface User { role: UserRole; }`
- **Prefer "as const Enums"**: Instead of simple string union types for known states, use the Object-as-Const pattern. This provides better scalability, refactoring support, and runtime value access.
  ```typescript
  // Avoid:
  // type Status = "draft" | "published";

  // Prefer:
  export const Status = {
    Draft: "draft",
    Published: "published",
  } as const;
  export type Status = (typeof Status)[keyof typeof Status];
  ```
- **Type Definitions**:
  - **`type`** vs **`interface`**: Prefer `interface` for extendable object shapes (like Component Props) and `type` for unions/primitives.
  - **Zod Integration**: Use `z.infer<typeof Schema>` for types derived from validation schemas.


## Routing (TanStack Router)
- **File-Based Routing**: All routes live in `apps/web/src/routes`.
- **Typesafe**: Rely on `routeTree.gen.ts`. Do not manually edit it.
- **Loaders**: Use `loader` in routes for data fetching (parallelized by default).
- **Search Params**: Validate with Zod schemas in `validateSearch`.
- **Links**: Use type-safe `<Link>` component. Never `<a>` tags for internal nav.
- **Thin Routes**: Route files should have **NO logic** or UI implementation. They should only define the route (loader, params) and import a **Page Component** from `features/{feature}/pages/`.

## Forms & Validation (TanStack Form + Zod)
- Use `@tanstack/react-form` for complex form state.
- Define schemas with **Zod** (`z`).
- Validations should be shared between backend and frontend if possible.

## SOLID in React (Component/Hook Design)
- **S**ingle Responsibility: One component = one UI concern. Extract queries/logic to hooks.
- **O**pen/Closed: Extend via props/slots (Composition), not by adding 10 boolean flags.
- **L**iskov Substitution: A custom `Button` should accept standard button props.
- **I**nterface Segregation: Props should be minimal. Don't pass a whole `User` object if only `name` is needed.
- **D**ependency Inversion: Components depend on props/hooks, not global fetches.

## Component & Typing Best Practices
- **Function Components Only**.
- **Avoid** `React.FC`.
- **Composition**: Use compound components or slots for complex UIs.
- **Hook Rules**: Custom hooks start with `use`. Colocate state unless shared.

```tsx
// Example: Composition & Typing
interface CardProps {
  title: string;
  footer?: React.ReactNode;
  children: React.ReactNode;
}

export function Card({ title, footer, children }: CardProps) {
  return (
    <div className="border rounded-md p-4">
      <h3 className="font-bold">{title}</h3>
      <div>{children}</div>
      {footer && <div className="mt-4">{footer}</div>}
    </div>
  );
}
```

## React 19 Specifics
- **Actions**: Use `useActionState` (or aliases) for mutations.
- **Suspense**: Embrace `<Suspense>` boundaries.
- **"use"**: Use `use(Promise)` or `use(Context)` instead of `useContext` or effects for data unwrapping.
- **Ref**: Pass `ref` as a prop (no `forwardRef` needed in React 19).

## Directory Structure (`apps/web/src`)
- `features/` — Domain-specific features (Auth, Timeline, Consent).
  - `{feature}/components/`
  - `{feature}/pages/` — Route-level page components.
  - `{feature}/hooks/`
  - `{feature}/types/`
  - `{feature}/api/`
- `components/` — Shared UI components (dumb/presentational).
  - `ui/` — Low-level primitives (buttons, inputs).
  - `common/` — Shared domain-agnostic components (Logo, etc.).
- `routes/` — File-system routing (Pages/Layouts).
- `lib/` — Utilities, API clients, helpers.
- `hooks/` — Global reusable hooks.
- `types/` — Global shared types. **Do not cluster in `index.ts`**. Use named files (e.g. `ethereum.ts`, `api.ts`) for context.

## Performance & Quality
- **Performance**: Trust React Compiler. Don't proactively `useMemo` unless necessary.
- **Type Safety**: Strictly typed props. Zod for runtime data.
- **Images**: Use optimized formats.

## Quick Mantra
**"Typesafe Routing. Biome Formatting. Composition over Configuration."**
