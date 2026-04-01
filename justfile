# Default recipe - show list
[private]
default:
    @just --list

# Set variables
export UI_PORT := "3000"
export API_PORT := "3001"

# Development: Run both API and UI concurrently
dev:
    @echo "Starting development servers..."
    pnpm exec concurrently \
        -n "API,UI" \
        -c "cyan,green" \
        "just dev-api" \
        "just dev-ui"

# Development API: Start Go API with hot reload
dev-api:
    @echo "Starting Go API on port {{API_PORT}}..."
    bash -c 'if command -v air >/dev/null 2>&1; then cd apps/api && air; else echo "Air not found, using go run (no hot reload)..." && cd apps/api && go run cmd/server/main.go; fi'

# Development UI: Start SolidStart dev server
dev-ui:
    @echo "Starting SolidStart dev server on port {{UI_PORT}}..."
    cd apps/ui && pnpm dev --port {{UI_PORT}}

# Build: Production build for both apps
build:
    @echo "Building production..."
    cd apps/ui && pnpm build
    cd apps/api && go build -o bin/server cmd/server/main.go

# Test: Run all tests
test:
    @echo "Running tests..."
    cd apps/ui && pnpm test
    cd apps/api && go test ./...

# Install: Install dependencies
install:
    @echo "Installing dependencies..."
    pnpm install
    cd apps/api && go mod tidy

# Database: Run migrations
db-migrate:
    @echo "Running database migrations..."
    cd apps/api && go run cmd/migrate/main.go

# Clean: Remove build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf apps/ui/dist
    rm -rf apps/ui/.output
    rm -rf apps/api/bin
    rm -rf apps/api/tmp
    rm -f apps/api/*.db
