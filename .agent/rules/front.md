---
trigger: glob
description: React + TypeScript frontend development rules.
globs: "*.ts", "*.tsx"
---

# React + TypeScript Rules
(React 19, Vite + Rolldown, TanStack Stack, Biome)

> **See also**: general.md for project philosophy and engineering principles.

We write **clean, maintainable, type-safe, SOLID-compliant** React code.

---

## 1. Tooling & Language Discipline

| Tool            | Purpose                                             |
| --------------- | --------------------------------------------------- |
| **Build**       | Vite 7 (Rolldown) + `@vitejs/plugin-react-swc`      |
| **Lint/Format** | Biome (`@biomejs/biome`). No ESLint/Prettier.       |
| **Styling**     | Tailwind CSS v4 + `clsx` + `tailwind-merge` + `cva` |
| **Forms**       | TanStack Form + Zod                                 |
| **Routing**     | TanStack Router (file-based)                        |
| **State**       | React Query for server state, Context for DI        |

### Dark Mode Pattern
**Always include dark mode variants** when styling components:
```tsx
// ❌ Bad: Light mode only
<div className="bg-white text-gray-900" />

// ✅ Good: Both themes
<div className="bg-white dark:bg-gray-900 text-gray-900 dark:text-white" />
```
- Use `dark:` prefix for dark mode overrides
- Theme class applied to `<html>` element
- Default: dark mode

---

## 2. TypeScript Standards

### Non-Negotiables
- **`strict: true`** — Always enabled.
- **No `any`** — Use `unknown` with type narrowing.
- **Explicit returns** — Exported functions must declare return types.

### Type Patterns
```typescript
// ❌ Bad: Inline union
interface User { role: "admin" | "user"; }

// ✅ Good: Extracted type
type UserRole = "admin" | "user";
interface User { role: UserRole; }

// ✅ Better: As-const enum
export const UserRole = {
  Admin: "admin",
  User: "user",
} as const;
export type UserRole = (typeof UserRole)[keyof typeof UserRole];
```

### Type vs Interface
- **`interface`** — Object shapes, component props (extendable).
- **`type`** — Unions, primitives, computed types.

---

## 3. Component Design

### SOLID Principles Applied
| Principle | React Application                                            |
| --------- | ------------------------------------------------------------ |
| **SRP**   | One component = one concern. Extract logic to hooks.         |
| **OCP**   | Extend via composition (slots, children), not boolean flags. |
| **LSP**   | Custom components accept standard HTML props.                |
| **ISP**   | Props should be minimal. Don't pass entire objects.          |
| **DIP**   | Components depend on props/hooks, not global imports.        |

### Best Practices
```tsx
// ✅ Good: Typed, composable, minimal props
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

### Rules
- **Function components only** — No class components.
- **No `React.FC`** — Use explicit prop typing.
- **Hooks start with `use`** — Always.

---

## 4. Routing (TanStack Router)

| Rule              | Description                                                      |
| ----------------- | ---------------------------------------------------------------- |
| **File-based**    | Routes live in `apps/web/src/routes/`                            |
| **Typesafe**      | Rely on generated `routeTree.gen.ts`                             |
| **Thin routes**   | Route files = routing config only. Import Page from `features/`. |
| **`<Link>` only** | Never use `<a>` for internal navigation.                         |

### Loaders & Search Params
```typescript
// Route with loader and typed search params
export const Route = createFileRoute('/timeline')({
  component: TimelinePage,
  validateSearch: z.object({ filter: z.string().optional() }),
  loader: async () => fetchTimeline(),
});
```

---

## 5. React 19 Specifics

| Feature      | Usage                                       |
| ------------ | ------------------------------------------- |
| **Actions**  | Use `useActionState` for mutations          |
| **Suspense** | Embrace `<Suspense>` boundaries             |
| **`use()`**  | Replace `useContext` with `use(Context)`    |
| **Refs**     | Pass `ref` as prop (no `forwardRef` needed) |

---

## 6. Error Handling & Testing

### Error Handling
- Wrap async operations in try/catch.
- Use error boundaries for component failures.
- Display user-friendly error messages, log details for debugging.

### Testing Philosophy
- **Unit tests** for hooks and utilities.
- **Component tests** with Testing Library.
- **E2E tests** for critical user flows (Playwright).

---

## 7. Directory Structure

```
src/
├── routes/       # Typesafe routing config
├── features/     # Feature vertical slices (auth, timeline)
├── components/   # ui/ (primitives), common/ (shared)
├── lib/          # webcrypto, api clients, utils
└── hooks/        # Global hooks
```

## 5. PR Checklist
- [ ] No `any`.
- [ ] Dark mode variants included.
- [ ] `pnpm audit` passes.
- [ ] No `console.log` of sensitive data.
- [ ] Zod validation on inputs.