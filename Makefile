.PHONY: build clean install test run dev help deps fmt vet lint docker-build release

# Variables
BINARY_NAME=lumine
VERSION?=dev
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags="-s -w -X main.version=$(VERSION)"

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help
	@echo "$(CYAN)Lumine Development Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

# Development
dev: ## Run in development mode with hot reload
	@echo "$(CYAN)Running in development mode...$(NC)"
	@$(GO) run .

run: build ## Build and run
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

watch: ## Watch for changes and rebuild (requires entr)
	@echo "$(CYAN)Watching for changes...$(NC)"
	@find . -name '*.go' | entr -r make run

# Build
build: deps ## Build the binary
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-cleanup: ## Build cleanup tool
	@echo "$(CYAN)Building cleanup tool...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup ./cmd/cleanup
	@echo "$(GREEN)✓ Cleanup tool built: $(BUILD_DIR)/lumine-cleanup$(NC)"

build-all: ## Build for all platforms
	@echo "$(CYAN)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "  → Linux AMD64..."
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup-linux-amd64 ./cmd/cleanup
	@echo "  → Linux ARM64..."
	@GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup-linux-arm64 ./cmd/cleanup
	@echo "  → macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup-darwin-amd64 ./cmd/cleanup
	@echo "  → macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup-darwin-arm64 ./cmd/cleanup
	@echo "  → Windows AMD64..."
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/lumine-cleanup-windows-amd64.exe ./cmd/cleanup
	@echo "$(GREEN)✓ All builds complete!$(NC)"

# Installation
install: build ## Install the binary to system
	@echo "$(CYAN)Installing to $(INSTALL_DIR)...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "$(GREEN)✓ Installation complete!$(NC)"

install-cleanup: build-cleanup ## Install cleanup tool
	@echo "$(CYAN)Installing cleanup tool to $(INSTALL_DIR)...$(NC)"
	@sudo cp $(BUILD_DIR)/lumine-cleanup $(INSTALL_DIR)/
	@echo "$(GREEN)✓ Cleanup tool installed!$(NC)"

install-all: build build-cleanup ## Install both lumine and cleanup tool
	@echo "$(CYAN)Installing all tools to $(INSTALL_DIR)...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo cp $(BUILD_DIR)/lumine-cleanup $(INSTALL_DIR)/
	@echo "$(GREEN)✓ All tools installed!$(NC)"

uninstall: ## Uninstall the binary from system
	@echo "$(YELLOW)Uninstalling...$(NC)"
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo rm -f $(INSTALL_DIR)/lumine-cleanup
	@echo "$(GREEN)✓ Uninstall complete!$(NC)"

# Dependencies
deps: ## Download and tidy dependencies
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies ready!$(NC)"

deps-update: ## Update all dependencies
	@echo "$(CYAN)Updating dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated!$(NC)"

# Testing
test: ## Run tests
	@echo "$(CYAN)Running tests...$(NC)"
	@$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@$(GO) test -cover -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

test-race: ## Run tests with race detector
	@echo "$(CYAN)Running tests with race detector...$(NC)"
	@$(GO) test -race -v ./...

benchmark: ## Run benchmarks
	@echo "$(CYAN)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

# Code Quality
fmt: ## Format code
	@echo "$(CYAN)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Format complete!$(NC)"

vet: ## Run go vet
	@echo "$(CYAN)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✓ Vet complete!$(NC)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(CYAN)Running linter...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)✓ Lint complete!$(NC)"

check: fmt vet ## Run fmt and vet
	@echo "$(GREEN)✓ Code check complete!$(NC)"

# Docker
docker-validate: ## Validate Docker installation
	@echo "$(CYAN)Validating Docker...$(NC)"
	@docker --version || (echo "$(RED)✗ Docker not installed$(NC)" && exit 1)
	@docker ps > /dev/null 2>&1 || (echo "$(RED)✗ Docker not running$(NC)" && exit 1)
	@echo "$(GREEN)✓ Docker is ready!$(NC)"

docker-build: ## Build Docker image
	@echo "$(CYAN)Building Docker image...$(NC)"
	@docker build -t lumine:$(VERSION) .
	@echo "$(GREEN)✓ Docker build complete!$(NC)"

docker-clean: ## Clean Docker resources
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	@docker system prune -f
	@echo "$(GREEN)✓ Docker cleanup complete!$(NC)"

# Database
db-setup: ## Setup all database services
	@echo "$(CYAN)Setting up databases...$(NC)"
	@docker compose -f docker-compose.db.yml up -d
	@echo "$(GREEN)✓ Databases are running!$(NC)"
	@echo ""
	@echo "$(CYAN)Database Access:$(NC)"
	@echo "  MySQL:        localhost:3306 (root/root)"
	@echo "  PostgreSQL:   localhost:5432 (postgres/postgres)"
	@echo "  MongoDB:      localhost:27017"
	@echo "  Redis:        localhost:6379"
	@echo ""
	@echo "$(CYAN)Admin Panels:$(NC)"
	@echo "  phpMyAdmin:   http://localhost:8080"
	@echo "  Adminer:      http://localhost:8081"
	@echo "  Mongo Express: http://localhost:8082"
	@echo "  Redis Commander: http://localhost:8083"

db-stop: ## Stop all database services
	@echo "$(YELLOW)Stopping databases...$(NC)"
	@docker compose -f docker-compose.db.yml down
	@echo "$(GREEN)✓ Databases stopped!$(NC)"

db-restart: ## Restart all database services
	@echo "$(CYAN)Restarting databases...$(NC)"
	@docker compose -f docker-compose.db.yml restart
	@echo "$(GREEN)✓ Databases restarted!$(NC)"

db-logs: ## Show database logs
	@docker compose -f docker-compose.db.yml logs -f

db-clean: ## Remove all database data (WARNING: destructive)
	@echo "$(RED)⚠️  This will delete all database data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose -f docker-compose.db.yml down -v; \
		echo "$(GREEN)✓ Database data removed!$(NC)"; \
	fi

db-reset: ## Reset databases (stop, remove, and restart fresh)
	@echo "$(CYAN)Resetting databases...$(NC)"
	@docker compose -f docker-compose.db.yml down -v
	@docker compose -f docker-compose.db.yml up -d
	@echo "$(GREEN)✓ Databases reset complete!$(NC)"

# Container Management
containers-list: ## List all Lumine containers
	@echo "$(CYAN)Lumine Containers:$(NC)"
	@docker ps -a --filter "name=lumine-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

containers-stop: ## Stop all Lumine containers
	@echo "$(YELLOW)Stopping all Lumine containers...$(NC)"
	@docker ps -q --filter "name=lumine-" | xargs -r docker stop
	@echo "$(GREEN)✓ All containers stopped!$(NC)"

containers-remove: ## Remove all Lumine containers
	@echo "$(RED)⚠️  This will remove all Lumine containers!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker ps -aq --filter "name=lumine-" | xargs -r docker rm -f; \
		echo "$(GREEN)✓ All containers removed!$(NC)"; \
	fi

containers-clean: ## Remove all stopped Lumine containers
	@echo "$(CYAN)Removing stopped containers...$(NC)"
	@docker ps -aq --filter "name=lumine-" --filter "status=exited" | xargs -r docker rm
	@echo "$(GREEN)✓ Stopped containers removed!$(NC)"

containers-prune: ## Remove all Lumine containers and volumes (DESTRUCTIVE!)
	@echo "$(RED)⚠️  WARNING: This will remove ALL Lumine containers and data!$(NC)"
	@echo "$(RED)This action cannot be undone!$(NC)"
	@read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker ps -aq --filter "name=lumine-" | xargs -r docker rm -f; \
		docker volume ls -q --filter "name=lumine_" | xargs -r docker volume rm; \
		docker network rm lumine 2>/dev/null || true; \
		echo "$(GREEN)✓ Complete cleanup done!$(NC)"; \
	else \
		echo "$(YELLOW)Cleanup cancelled.$(NC)"; \
	fi

# Network Management
network-create: ## Create Lumine network
	@echo "$(CYAN)Creating Lumine network...$(NC)"
	@docker network create lumine 2>/dev/null || echo "$(YELLOW)Network already exists$(NC)"
	@echo "$(GREEN)✓ Network ready!$(NC)"

network-remove: ## Remove Lumine network
	@echo "$(YELLOW)Removing Lumine network...$(NC)"
	@docker network rm lumine 2>/dev/null || echo "$(YELLOW)Network doesn't exist$(NC)"
	@echo "$(GREEN)✓ Network removed!$(NC)"

network-inspect: ## Inspect Lumine network
	@docker network inspect lumine 2>/dev/null || echo "$(RED)Network doesn't exist$(NC)"

# Volume Management
volumes-list: ## List all Lumine volumes
	@echo "$(CYAN)Lumine Volumes:$(NC)"
	@docker volume ls --filter "name=lumine_" --format "table {{.Name}}\t{{.Driver}}\t{{.Mountpoint}}"

volumes-remove: ## Remove all Lumine volumes (WARNING: data loss!)
	@echo "$(RED)⚠️  This will delete all database data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker volume ls -q --filter "name=lumine_" | xargs -r docker volume rm; \
		echo "$(GREEN)✓ All volumes removed!$(NC)"; \
	fi

volumes-prune: ## Remove unused Lumine volumes
	@echo "$(CYAN)Removing unused volumes...$(NC)"
	@docker volume prune -f --filter "label=com.docker.compose.project=lumine"
	@echo "$(GREEN)✓ Unused volumes removed!$(NC)"

# Cleanup
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete!$(NC)"

clean-all: clean db-clean docker-clean ## Clean everything including Docker
	@echo "$(GREEN)✓ Full cleanup complete!$(NC)"

clean-containers: containers-prune ## Alias for containers-prune
	@echo "$(GREEN)✓ Containers cleaned!$(NC)"

cleanup-interactive: build-cleanup ## Run interactive cleanup tool
	@echo "$(CYAN)Starting interactive cleanup...$(NC)"
	@$(BUILD_DIR)/lumine-cleanup

cleanup-stop: build-cleanup ## Stop all containers using cleanup tool
	@echo "1" | $(BUILD_DIR)/lumine-cleanup

cleanup-remove: build-cleanup ## Remove containers using cleanup tool
	@echo "2" | $(BUILD_DIR)/lumine-cleanup

cleanup-nuclear: build-cleanup ## Nuclear cleanup using cleanup tool
	@echo "4" | $(BUILD_DIR)/lumine-cleanup

clean-everything: ## Nuclear option - remove EVERYTHING (containers, volumes, networks, builds)
	@echo "$(RED)⚠️  NUCLEAR OPTION: This will remove EVERYTHING!$(NC)"
	@echo "$(RED)- All Lumine containers$(NC)"
	@echo "$(RED)- All Lumine volumes (database data)$(NC)"
	@echo "$(RED)- Lumine network$(NC)"
	@echo "$(RED)- Build artifacts$(NC)"
	@echo "$(RED)- Docker cache$(NC)"
	@echo ""
	@read -p "Type 'DESTROY' to confirm: " confirm; \
	if [ "$$confirm" = "DESTROY" ]; then \
		$(MAKE) containers-prune; \
		$(MAKE) volumes-remove; \
		$(MAKE) network-remove; \
		$(MAKE) clean; \
		docker system prune -af; \
		echo "$(GREEN)✓ Everything destroyed!$(NC)"; \
	else \
		echo "$(YELLOW)Destruction cancelled.$(NC)"; \
	fi

# Release
release: clean build-all ## Create release artifacts
	@echo "$(CYAN)Creating release artifacts...$(NC)"
	@cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "$(GREEN)✓ Release artifacts ready in $(BUILD_DIR)/$(NC)"

# Info
info: ## Show project information
	@echo "$(CYAN)Lumine Project Information$(NC)"
	@echo ""
	@echo "  Version:      $(VERSION)"
	@echo "  Go Version:   $$(go version | cut -d' ' -f3)"
	@echo "  Build Dir:    $(BUILD_DIR)"
	@echo "  Install Dir:  $(INSTALL_DIR)"
	@echo ""
	@echo "$(CYAN)Docker Status:$(NC)"
	@docker --version 2>/dev/null || echo "  $(RED)Not installed$(NC)"
	@docker compose version 2>/dev/null || echo "  $(RED)Compose not installed$(NC)"
	@echo ""

.DEFAULT_GOAL := help
