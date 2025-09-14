# 🛠️ Makefile for hashing Package
# Provides build automation and development tools

.PHONY: help build test clean lint fmt vet security coverage dev-setup

# Default target
help: ## 📚 Show this help message
	@echo "🛠️  Available targets for hashing:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Determine if we're in src structure or not
WORKING_DIR := $(shell if [ -f "src/go.mod" ]; then echo "src"; else echo "."; fi)

build: ## 🔨 Build the package
	@echo "🔨 Building hashing..."
	@cd $(WORKING_DIR) && go build -v ./...

test: ## 🧪 Run tests with coverage
	@echo "🧪 Running tests for hashing..."
	@mkdir -p coverage
	@cd $(WORKING_DIR) && go test -v -race -coverprofile=../coverage/coverage.out -covermode=atomic ./...
	@cd $(WORKING_DIR) && go tool cover -func=../coverage/coverage.out

test-short: ## ⚡ Run tests without coverage
	@echo "⚡ Running quick tests for hashing..."
	@cd $(WORKING_DIR) && go test -v -short ./...

clean: ## 🧹 Clean build artifacts and cache
	@echo "🧹 Cleaning hashing..."
	@cd $(WORKING_DIR) && go clean -cache -testcache -modcache
	@rm -rf coverage/

lint: ## 🔍 Run golangci-lint
	@echo "🔍 Linting hashing..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd $(WORKING_DIR) && golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run 'make dev-setup' first."; \
	fi

fmt: ## 🎨 Format code
	@echo "🎨 Formatting hashing..."
	@cd $(WORKING_DIR) && go fmt ./...

vet: ## 🔍 Run go vet
	@echo "🔍 Running go vet for hashing..."
	@cd $(WORKING_DIR) && go vet ./...

security: ## 🔒 Run security scans
	@echo "🔒 Running security scans for hashing..."
	@if command -v gosec >/dev/null 2>&1; then \
		cd $(WORKING_DIR) && gosec ./...; \
	else \
		echo "⚠️  gosec not installed. Run 'make dev-setup' first."; \
	fi
	@if command -v govulncheck >/dev/null 2>&1; then \
		cd $(WORKING_DIR) && govulncheck ./...; \
	else \
		echo "⚠️  govulncheck not installed. Run 'make dev-setup' first."; \
	fi

coverage: test ## 📊 Generate coverage report
	@echo "📊 Generating coverage report for hashing..."
	@cd $(WORKING_DIR) && go tool cover -html=../coverage/coverage.out -o ../coverage/coverage.html
	@echo "📊 Coverage report: coverage/coverage.html"

dev-setup: ## 🛠️ Install development tools
	@echo "🛠️ Installing development tools for hashing..."
	@cd $(WORKING_DIR) && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@cd $(WORKING_DIR) && go install github.com/securecodewarrior/gosec/cmd/gosec@latest
	@cd $(WORKING_DIR) && go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development tools installed"

all: clean fmt vet lint test security ## 🚀 Run complete CI pipeline
	@echo "🎉 All checks passed for hashing!"

# Dependencies
deps: ## 📥 Download dependencies
	@echo "📥 Downloading dependencies for hashing..."
	@cd $(WORKING_DIR) && go mod download
	@cd $(WORKING_DIR) && go mod tidy

update-deps: ## ⬆️ Update dependencies
	@echo "⬆️ Updating dependencies for hashing..."
	@cd $(WORKING_DIR) && go get -u ./...
	@cd $(WORKING_DIR) && go mod tidy
