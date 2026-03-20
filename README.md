# PocketPanel

A lightweight admin panel built with Svelte 5 + Vite (frontend) and Go + Fiber (backend).

## Architecture

```
pocketpanel/
├── apps/
│   ├── web/              # Svelte 5 + Vite (Port 3000)
│   └── api/              # Go 1.26 + Fiber (Port 3001)
├── packages/
│   └── schema/           # JSON Schemas for shared types
├── package.json          # Root workspace config
├── pnpm-workspace.yaml   # pnpm workspace definition
└── Makefile              # Task orchestration
```

## Prerequisites

- Go 1.26
- Node.js 20+ 
- pnpm 10+

## Getting Started

```bash
# Install dependencies
make install

# Start development servers (runs both)
make dev

# Or run individually
make dev-api   # Go API on http://localhost:3001
make dev-web   # Svelte dev server on http://localhost:3000
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Runs both API and Web concurrently |
| `make dev-api` | Go API with hot reload (Air) |
| `make dev-web` | Vite dev server |
| `make build` | Production builds for both |
| `make test` | Run all tests |
| `make db-migrate` | Run SQLite migrations |
| `make clean` | Remove build artifacts |

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login with username/password
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/me` - Get current user (requires auth)

## Authentication

The app uses session-based authentication with HTTP-only cookies. On login, a session cookie is set which is automatically sent with subsequent requests. The logout endpoint destroys the session server-side.

## Environment Variables

### API (`apps/api/.env`)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `:3001` | API server port |
| `DATABASE_PATH` | `pocketpanel.db` | SQLite database file path |
| `ENV` | `development` | Environment (development/production) |
| `CORS_ORIGINS` | `http://localhost:3000` | Comma-separated allowed CORS origins |

## Development

- Frontend runs on port 3000 with API proxy to port 3001
- API runs on port 3001 with CORS configured for localhost:3000
- SQLite database stored at `apps/api/pocketpanel.db`

## Tech Stack

**Frontend:**
- Svelte 5 (Runes)
- Vite 6
- TypeScript
- Vitest for testing

**Backend:**
- Go 1.26
- Fiber v3
- GORM with SQLite
- Session-based authentication
