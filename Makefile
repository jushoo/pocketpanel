.PHONY: dev dev-api dev-web build test install clean db-migrate

# Ports: Web 3000, API 3001
WEB_PORT=3000
API_PORT=3001

dev:
	@echo "Starting development servers..."
	pnpm exec concurrently \
		-n "API,WEB" \
		-c "cyan,green" \
		"make dev-api" \
		"make dev-web"

dev-api:
	@echo "Starting Go API on port $(API_PORT)..."
	@if command -v air >/dev/null 2>&1; then \
		cd apps/api && air; \
	else \
		echo "Air not found, using go run (no hot reload)..."; \
		cd apps/api && go run cmd/server/main.go; \
	fi

dev-web:
	@echo "Starting Svelte dev server on port $(WEB_PORT)..."
	cd apps/web && pnpm dev --port $(WEB_PORT)

build:
	@echo "Building production..."
	cd apps/web && pnpm build
	cd apps/api && go build -o bin/server cmd/server/main.go

test:
	@echo "Running tests..."
	cd apps/web && pnpm test
	cd apps/api && go test ./...

install:
	@echo "Installing dependencies..."
	pnpm install
	cd apps/api && go mod tidy

db-migrate:
	@echo "Running database migrations..."
	cd apps/api && go run cmd/migrate/main.go

clean:
	@echo "Cleaning build artifacts..."
	rm -rf apps/web/dist
	rm -rf apps/api/bin
	rm -rf apps/api/tmp
	rm -f apps/api/*.db
