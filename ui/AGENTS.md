# AGENTS.md — Next.js Project Standards

This file defines conventions for all contributors and AI agents working in this codebase.
Follow these rules without exception unless explicitly overridden in a sub-directory's own `AGENTS.md`.

---

## Project Philosophy

- **Business logic lives apart from UI.** Components render; hooks and utilities think.
- **Logic is abstracted, not inlined.** If it can be named and extracted, it should be.
- **Context is a first-class tool.** Use React Context for shared state — don't prop-drill.
- **Predictable structure beats clever structure.** Every file should be findable without a search.
- **Explicit over implicit.** Name things for what they do, not where they live.

---

## Naming Conventions

### Files & Directories → `snake_case` everywhere

```
components/
  user_profile/
    user_profile.tsx          ← component (visual only)
    use_user_profile.ts       ← logic hook
    user_profile.types.ts     ← types/interfaces
    index.ts                  ← barrel export

app/
  dashboard/
    page.tsx
    loading.tsx
    error.tsx

lib/
  format_currency.ts
  validate_email.ts

hooks/
  use_auth.ts
  use_debounce.ts

services/
  user_service.ts
  payment_service.ts

types/
  user.types.ts
  api.types.ts

context/
  auth_context.tsx
  theme_context.tsx
```

- No `camelCase`, `PascalCase`, or `kebab-case` for filenames — ever.
- Exception: Next.js special files (`page.tsx`, `layout.tsx`, `middleware.ts`, etc.) stay as-is — the framework requires them.

---

## Architecture: Separation of Concerns

### The Rule

> A component file must not contain business logic. It renders JSX and calls hooks — nothing else.

### Logic Abstraction

If a piece of logic can be named, extract it. Don't leave complex expressions inline in JSX or hooks.

```ts
// ❌ Bad — logic buried inline
const display_name = user.first_name
  ? `${user.first_name} ${user.last_name}`.trim()
  : user.email.split('@')[0]

// ✅ Good — extracted to lib/
// lib/format_display_name.ts
export function formatDisplayName(user: User): string {
  if (user.first_name) return `${user.first_name} ${user.last_name}`.trim()
  return user.email.split('@')[0]
}
```

The same applies inside hooks — if a block of logic can be pulled into a pure function in `lib/`, do it.

### What Counts as Business Logic

- Data fetching and mutation
- Validation
- Calculations and transformations
- Auth checks
- Side effects beyond pure UI (scroll position is UI; API calls are not)
- Complex conditional logic

### Component (`*.tsx`) — Visual Only

```tsx
// components/invoice_summary/invoice_summary.tsx
import { useInvoiceSummary } from './use_invoice_summary'
import type { InvoiceSummaryProps } from './invoice_summary.types'

export function InvoiceSummary({ invoice_id }: InvoiceSummaryProps) {
  const { total, line_items, is_loading, error } = useInvoiceSummary(invoice_id)

  if (is_loading) return <InvoiceSkeleton />
  if (error) return <ErrorMessage message={error.message} />

  return (
    <section>
      <LineItemList items={line_items} />
      <TotalDisplay total={total} />
    </section>
  )
}
```

### Hook (`use_*.ts`) — Logic Only

```ts
// components/invoice_summary/use_invoice_summary.ts
import { useQuery } from '@tanstack/react-query'
import { fetchInvoice } from '@/services/invoice_service'
import { formatLineItems } from '@/lib/format_line_items'

export function useInvoiceSummary(invoice_id: string) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['invoice', invoice_id],
    queryFn: () => fetchInvoice(invoice_id),
  })

  return {
    total: data?.total ?? 0,
    line_items: formatLineItems(data?.items ?? []),
    is_loading: isLoading,
    error,
  }
}
```

---

## Directory Structure

```
project-root/
├── app/                    # Next.js App Router pages & layouts
├── components/             # Shared UI components
│   └── [component_name]/
│       ├── index.ts
│       ├── [name].tsx
│       ├── use_[name].ts
│       └── [name].types.ts
├── hooks/                  # Global/shared hooks
├── context/                # React Context providers and consumers
├── lib/                    # Pure utility functions (no React)
├── services/               # External API calls, data access
├── types/                  # Shared TypeScript types
├── constants/              # App-wide constants
└── public/                 # Static assets
```

---

## Do This

- ✅ One component per file
- ✅ Export types from a dedicated `.types.ts` file
- ✅ Use barrel `index.ts` files for clean imports
- ✅ Use `snake_case` for all custom file and directory names
- ✅ Keep hooks pure and focused
- ✅ Extract named logic to `lib/` — no complex inline expressions
- ✅ Use React Context for shared cross-component state
- ✅ Fetch data in hooks or server components — not inside JSX
- ✅ Use `zod` (or equivalent) for runtime validation, in `lib/` or `services/`
- ✅ Type all function inputs and outputs explicitly
- ✅ Use `loading.tsx` / `error.tsx` for route-level states

## Don't Do This

- ❌ Business logic inside component files
- ❌ Complex logic inlined inside JSX or hooks — extract and name it
- ❌ Prop-drilling more than one level when Context fits
- ❌ One god-context holding all app state — split by domain
- ❌ `useState` + `useEffect` for data fetching — use a query library
- ❌ Inline `fetch()` calls inside components
- ❌ Deeply nested ternaries in JSX
- ❌ `any` type (use `unknown` and narrow it)
- ❌ `// @ts-ignore` without a comment explaining why
- ❌ Importing across feature boundaries without going through `index.ts`
- ❌ God-hooks that do 10 unrelated things
- ❌ Hard-coded strings for routes, keys, or config — use `constants/`
- ❌ Mixing server and client concerns in the same file without clear intent

---

## TypeScript

- Strict mode is on. Do not disable it.
- Prefer `interface` for object shapes, `type` for unions and aliases.
- No `enum` — use `as const` objects instead.

```ts
// Good
export const USER_ROLES = {
  ADMIN: 'admin',
  MEMBER: 'member',
} as const

export type UserRole = (typeof USER_ROLES)[keyof typeof USER_ROLES]
```

---

## Imports

Use path aliases, not relative hell:

```ts
// Good
import { formatCurrency } from '@/lib/format_currency'

// Bad
import { formatCurrency } from '../../../lib/format_currency'
```

Order: external packages → internal aliases → relative imports. Keep them grouped with a blank line between groups.

---

## ESLint

ESLint is enforced. All code must pass linting with zero errors before committing.

Key rules in effect:

- `@typescript-eslint/no-explicit-any` — no `any`, full stop
- `@typescript-eslint/explicit-function-return-type` — type your return values
- `react-hooks/rules-of-hooks` — hooks at the top level only
- `react-hooks/exhaustive-deps` — all hook dependencies declared
- `import/order` — imports must follow the defined group order (see Imports section)
- `no-console` — use a logger utility, not `console.log`

Do not disable ESLint rules with `eslint-disable` comments unless absolutely unavoidable, and always include a comment explaining why.

```ts
// ❌ Bad
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const data: any = response

// ✅ Good — narrow the type instead
const data: unknown = response
if (typeof data === 'object' && data !== null) { ... }
```

---

## Prettier

Prettier handles all formatting. Do not manually format code — let Prettier do it.

Config (`.prettierrc`):

```json
{
  "semi": false,
  "singleQuote": true,
  "trailingComma": "all",
  "printWidth": 100,
  "tabWidth": 2,
  "arrowParens": "always"
}
```

- Run `prettier --write` before committing, or configure your editor to format on save.
- Do not argue with Prettier. If it reformats something, that's the correct format.
- Agents: emit code that already conforms to these settings. Do not add semicolons, use single quotes for strings, and include trailing commas in multi-line structures.

---

## React Context

Use Context for state that multiple components across the tree need — auth, theme, locale, feature flags, and so on. Do not prop-drill more than one level when Context is the right fit.

### Structure

Each context lives in `context/` with three named exports: the context itself, a provider component, and a consumer hook.

```ts
// context/auth_context.tsx
'use client'

import { createContext, useContext, useState } from 'react'
import type { AuthContextValue, AuthProviderProps } from './auth_context.types'

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)

  // business logic lives here, not in components
  function logout() {
    setUser(null)
    // ... cleanup
  }

  return (
    <AuthContext.Provider value={{ user, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
```

### Rules

- Always guard the consumer hook with a null check and a descriptive error.
- Keep Context values lean — don't shove the entire app state into one context.
- Business logic inside the provider is fine; business logic inside a component consuming context is not.
- Do not reach for Context for purely local state — `useState` in a hook is still the right call there.

---

## Function & Component Declaration Style

Prefer arrow functions assigned to `const` over `function` declarations.

### Components & hooks

```ts
// ✅ Good
export const InvoiceSummary = ({ invoice_id }: InvoiceSummaryProps) => {
  // ...
}

export const useInvoiceSummary = (invoice_id: string) => {
  // ...
}
```

### Pages

Named const + default export on a separate line:

```ts
// ✅ Good
const DashboardPage = () => {
  // ...
}

export default DashboardPage
```

### Never

```ts
// ❌ Bad
export default function DashboardPage() { ... }
export default function MyComponent() { ... }
```

This applies to all components, hooks, and page files without exception.

---



- Default to **Server Components**. Add `'use client'` only when you need interactivity or browser APIs.
- Never put data fetching in a Client Component when a Server Component can do it.
- Keep Client Components as leaf nodes where possible.

---

## Agents: When Generating Code

1. Check this file first.
2. Match the naming and structure conventions exactly — no camelCase filenames, no logic in components.
3. When adding a component, generate all three files: `.tsx`, `use_*.ts`, `.types.ts`, `index.ts`.
4. When unsure where something belongs: `lib/` for pure functions, `services/` for I/O, `hooks/` if it needs React, `context/` if it's shared state across the tree.
5. Extract any non-trivial logic into a named function in `lib/` — do not leave it inline.
6. Do not consolidate files to save tokens. The structure exists for a reason.
7. Emit Prettier-compliant code: single quotes, no semicolons, trailing commas, 100-char line width.
8. Do not emit ESLint violations — no `any`, no unused variables, no missing deps in hooks.
