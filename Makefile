.PHONY: build run clean install help dev

# Build Lumine
build:
	@echo "🔨 Building Lumine..."
	@go build -mod=mod -o lumine .
	@if [ -f lumine ]; then \
		echo "✅ Build complete! Binary: ./lumine"; \
	else \
		echo "⚠️  Build completed but binary not found"; \
	fi

# Run Lumine (using go run for development)
run:
	@echo "🚀 Starting Lumine..."
	@go run -mod=mod .

# Run with build
run-build: build
	@if [ -f lumine ]; then \
		echo "🚀 Starting Lumine..."; \
		./lumine; \
	else \
		echo "❌ Binary not found, using go run instead..."; \
		go run -mod=mod .; \
	fi

# Development mode (go run)
dev:
	@echo "🔧 Running in development mode..."
	@go run -mod=mod .

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f lumine
	@go clean
	@echo "✅ Clean complete!"

# Install to system
install: build
	@if [ -f lumine ]; then \
		echo "📦 Installing Lumine to /usr/local/bin..."; \
		sudo cp lumine /usr/local/bin/; \
		echo "✅ Installed! Run with: lumine"; \
	else \
		echo "❌ Binary not found. Run 'make build' first"; \
		exit 1; \
	fi

# Tidy dependencies
tidy:
	@echo "📦 Tidying dependencies..."
	@rm -rf vendor
	@go mod tidy
	@echo "✅ Dependencies updated!"

# Show help
help:
	@echo "Lumine - Modern Development Stack Manager"
	@echo ""
	@echo "Available commands:"
	@echo "  make build      - Build Lumine binary"
	@echo "  make run        - Run Lumine (development mode)"
	@echo "  make run-build  - Build then run"
	@echo "  make dev        - Run in development mode"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make tidy       - Update dependencies"
	@echo "  make help       - Show this help"
	@echo ""
	@echo "Quick start:"
	@echo "  make run        # Fastest way to start"
	@echo "  make dev        # Same as run"
	@echo "  make run-build  # Build first, then run"

# Default target
.DEFAULT_GOAL := help
