## ---------------- Build & Run ----------------

# Build the application
build:
	@echo "🔨 Building application..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@echo "🚀 Running application..."
	@go run ./cmd/api

# Watch for changes (dev only)
watch:
	@echo "👀 Watching for changes..."
	@air -c .air.toml

# Install dependencies & setup environment
setup:
	@echo "📦 Setting up project..."
	@go mod download & go mod tidy
	@cp .env.example .env || true

# Generate Swagger API docs
gen-docs:
	@echo "📖 Generating API documentation..."
	@swag init -g cmd/api/main.go -o cmd/docs


# Setup & start environment container
env-up:
	@echo "🐘 Starting environment container..."

	@docker compose --env-file .env up --build -d

# Shutdown environment container
env-down:
	@echo "🛑 Stopping environment container..."
	@docker compose --env-file .env down
