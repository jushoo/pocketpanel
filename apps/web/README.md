# PocketPanel Web

SvelteKit frontend for PocketPanel admin panel.

## Project Structure

```
src/
├── lib/              # Shared code (components, utilities, stores)
│   ├── components/   # Reusable UI components
│   └── utils.ts      # Helper functions
├── routes/           # SvelteKit file-based routing
│   ├── +layout.svelte    # Root layout with global styles
│   ├── +page.svelte      # Home/dashboard page
│   └── login/
│       └── +page.svelte  # Login page
├── app.css          # Global styles & Tailwind imports
├── app.d.ts         # TypeScript declarations
└── app.html         # HTML template
static/              # Static assets (favicon, images, etc.)
svelte.config.js     # SvelteKit configuration
vite.config.ts       # Vite configuration
tsconfig.json        # TypeScript configuration
tailwind.config.ts   # Tailwind CSS configuration
```

## Development

```bash
# Install dependencies (from repo root)
pnpm install

# Start dev server (from repo root)
make dev-web

# Or directly in this directory
pnpm dev
```

The dev server runs on http://localhost:3000

## Building

```bash
# Production build (from repo root)
make build

# Or directly in this directory
pnpm build
```

Static files are output to `build/` directory.

## Testing

```bash
# Run unit tests (from repo root)
make test

# Or directly in this directory
pnpm test
```

## Code Quality

```bash
# Check TypeScript and Svelte
pnpm check

# Lint code
pnpm lint

# Format code
pnpm format
```

## Tech Stack

- **Framework:** SvelteKit 2 with static adapter
- **UI Library:** Svelte 5 (Runes)
- **Styling:** Tailwind CSS 4
- **Icons:** Lucide Svelte
- **Testing:** Vitest with browser mode (Playwright)
- **Linting:** ESLint + Prettier

## Environment Variables

Create `.env` file in this directory if needed:

```
PUBLIC_API_URL=http://localhost:3001
```

Note: The API is proxied in development via Vite's dev server configuration.

## Project Recreation

This project was created with:

```bash
pnpm dlx sv@0.12.8 create --template minimal --types ts \
  --add prettier eslint vitest="usages:unit,component" \
  tailwindcss="plugins:typography" \
  sveltekit-adapter="adapter:static" \
  --install pnpm web
```
