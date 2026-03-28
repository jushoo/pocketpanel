# Agent Guidelines for PocketPanel

This is a monorepo with a SolidStart frontend (Port 3000) and Go backend (Port 3001).

## Build & Development Commands

### Root Level
```bash
make dev              # Start both dev servers
make build            # Build production
make test             # Run all tests
make install          # Install dependencies
make db-migrate       # Database migrations
make clean            # Clean build artifacts
```

### Individual Apps

**Web (SolidStart):**
```bash
cd apps/web
pnpm dev              # Dev server
pnpm build            # Production build
pnpm check            # TypeScript check
pnpm lint             # ESLint + Prettier check
pnpm format           # Auto-format with Prettier
pnpm test             # Run all tests
pnpm test:unit -- --run --reporter=verbose <pattern>  # Single test
```

**API (Go):**
```bash
cd apps/api
air                   # Development with hot reload
go build -o bin/server cmd/server/main.go  # Build
go test ./...         # Run all tests
go test -v ./internal/handlers/...  # Specific package
go test -run TestName ./...         # Single test
```

## Code Style

### TypeScript/Solid Import Ordering
```typescript
import { redirect } from '@solidjs/router';
import { Button } from '~/components/ui/button';
import { describe, it, expect } from 'vitest';
import type { RouteDefinition } from '@solidjs/router';
```

### TypeScript/Solid Naming
- **Components**: PascalCase (`Button.tsx`)
- **Utilities**: camelCase (`formatDate.ts`)
- **Routes**: kebab-case files (`index.tsx`, `about.tsx`)
- **Props/State**: Use `createSignal()` and `createStore()` for state, destructure props

```tsx
import { createSignal } from 'solid-js';

function Button(props: { label: string }) {
  const [count, setCount] = createSignal(0);
  return <button>{props.label}</button>;
}
```

### TypeScript/Solid Error Handling
Use standard try/catch with proper error responses:

```typescript
try {
  const data = await fetchData();
} catch (error) {
  throw new Response('User-friendly message', { status: 500 });
}
```

### TypeScript/Solid Testing
```typescript
import { describe, it, expect } from 'vitest';

describe('function', () => {
  it('should work', () => {
    expect(value).toBe(expected);
  });
});
```

### Go Import Ordering
```go
import (
    "os"
    "strings"
    
    "github.com/gofiber/fiber/v3"
    "gorm.io/gorm"
    
    "pocketpanel/api/internal/models"
)
```

### Go Naming
- **Packages**: lowercase, no underscores (`handlers`, `middleware`)
- **Exported**: PascalCase (`AuthHandler`, `NewServer`)
- **Unexported**: camelCase (`getEnv`, `validateToken`)
- **Interfaces**: `-er` suffix (`Handler`, `Reader`)

### Go Error Handling
Return errors with context, don't panic:

```go
func (h *Handler) Process(id string) (*Result, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }
    data, err := h.repo.Get(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get data: %w", err)
    }
    return data, nil
}
```

## UI Styling (Web)

### Component Usage
Install components via CLI:
```bash
cd apps/web && npx shadcn@latest add <component-name>
```

### Styling Rules
- **NEVER modify CSS theme** in `src/app.css`
- **Use theme CSS variables**: `bg-background`, `text-foreground`, `text-muted-foreground`
- **NO hardcoded colors**: Don't use `bg-[#0d0d0d]` or `text-white/90`
- **Trust shadcn defaults**: Don't override component styles unnecessarily

### Design Principles
- **Clean** - Minimal visual noise, generous whitespace
- **Minimal** - Only essential elements
- **Basic** - Simple layouts, straightforward interactions
- **Ollama-style** - Centered content, subtle borders, muted text hierarchy

### Good Example
```tsx
<div class="min-h-screen flex items-center justify-center bg-background">
  <div class="w-full max-w-sm px-6">
    <h1 class="text-xl font-medium text-foreground">Title</h1>
    <p class="text-sm text-muted-foreground">Subtitle</p>
    <Button class="w-full">Action</Button>
  </div>
</div>
```

## Environment Configuration

**Web (`apps/web/.env`):**
- `API_URL` - Backend URL (default: `http://localhost:3001`)

**API (`apps/api/.env`):**
- `PORT` - API port (default: `:3001`)
- `DATABASE_PATH` - SQLite file path
- `ENV` - `development` or `production`
- `CORS_ORIGINS` - Comma-separated allowed origins
- `SERVERS_PATH` - Path for server files (default: `/opt/pocketpanel/servers`)

## Authentication Flow

1. Login posts to `/api/v1/auth/login`
2. Backend sets HTTP-only `session_id` cookie
3. Frontend proxies cookie on subsequent requests
4. Protected routes check session via `/api/v1/me`
5. Logout calls `/api/v1/auth/logout` and clears cookie

## Before Committing

1. Run type checks: `pnpm check` (web), `go build ./...` (api)
2. Run linting: `pnpm lint` (web), `go vet ./...` (api)
3. Run tests: `pnpm test` (web), `go test ./...` (api)
4. Format code: `pnpm format` (web), `gofmt -w .` (api)
