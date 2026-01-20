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
- **TypeScript**: `strict: true`. No `any`. Use `unknown` + narrowing.
- **Styling**: Tailwind CSS v4 (`@tailwindcss/vite`) + `clsx` + `tailwind-merge` + `cva`.

## Routing (TanStack Router)
- **File-Based Routing**: All routes live in `apps/web/src/routes`.
- **Typesafe**: Rely on `routeTree.gen.ts`. Do not manually edit it.
- **Loaders**: Use `loader` in routes for data fetching (parallelized by default).
- **Search Params**: Validate with Zod schemas in `validateSearch`.
- **Links**: Use type-safe `<Link>` component. Never `<a>` tags for internal nav.

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
- `routes/` — File-system routing (Pages/Layouts).
- `components/` — Reusable, shared UI components.
  - `ui/` — Low-level primitives.
- `lib/` — Utilities, API clients, helpers.
- `hooks/` — Global reusable hooks.

## Performance & Quality
- **Performance**: Trust React Compiler. Don't proactively `useMemo` unless necessary.
- **Type Safety**: Strictly typed props. Zod for runtime data.
- **Images**: Use optimized formats.

## Quick Mantra
**"Typesafe Routing. Biome Formatting. Composition over Configuration."**
