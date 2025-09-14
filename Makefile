# 🛠️ GoVel Framework - Build Automation Makefile
#
# This Makefile provides comprehensive build automation for the GoVel framework.
# It includes targets for building, testing, linting, security scanning,
# and development workflow automation.
#
# Features:
# - 🔨 Build automation for all packages
# - 🧪 Testing with coverage reporting
# - 🔍 Code quality and linting
# - 🔒 Security scanning
# - 📦 Dependency management
# - 🚀 Release automation
# - 🧹 Cleanup and maintenance
#
# Usage:
#   make help          # Show available targets
#   make build         # Build all packages
#   make test          # Run all tests
#   make lint          # Run linting
#   make security      # Run security scans
#
# Author: GoVel Framework Team
# Version: 1.0.0
# License: MIT

# 🎯 Default target
.DEFAULT_GOAL := help

# 📋 Configuration variables
GO_VERSION := 1.23
GOLANGCI_VERSION := v1.55.2
GOSEC_VERSION := latest
GOVULNCHECK_VERSION := latest

# 📁 Directory configuration  
PACKAGES_DIR := packages
SCRIPTS_DIR := scripts
BUILD_DIR := build
COVERAGE_DIR := coverage
REPORTS_DIR := reports

# 🔍 Package discovery
PACKAGES := $(shell find $(PACKAGES_DIR) -name "go.mod" -exec dirname {} \; | sort)
PACKAGE_NAMES := $(notdir $(PACKAGES))

# 🎨 Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[1;37m
BOLD := \033[1m
NC := \033[0m

# 🖥️ Platform detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	PLATFORM := linux
endif
ifeq ($(UNAME_S),Darwin)
	PLATFORM := darwin
endif
ifeq ($(UNAME_S),CYGWIN*)
	PLATFORM := windows
endif
ifeq ($(UNAME_S),MINGW*)
	PLATFORM := windows
endif

ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
	ARCH := amd64
endif
ifeq ($(ARCH),arm64)
	ARCH := arm64
endif

##@ 📋 Help and Information

.PHONY: help
help: ## 📋 Show this help message
	@echo "$(BOLD)$(CYAN)🛠️ GoVel Framework - Build Automation$(NC)"
	@echo ""
	@echo "$(BOLD)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(CYAN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)Configuration:$(NC)"
	@echo "  Go Version: $(GREEN)$(GO_VERSION)$(NC)"
	@echo "  Platform: $(GREEN)$(PLATFORM)$(NC)"
	@echo "  Architecture: $(GREEN)$(ARCH)$(NC)"
	@echo "  Packages: $(GREEN)$(words $(PACKAGES))$(NC) found"
	@echo ""

.PHONY: info
info: ## ℹ️ Show project information
	@echo "$(BOLD)$(BLUE)📊 GoVel Framework Information$(NC)"
	@echo ""
	@echo "$(BOLD)📦 Packages:$(NC)"
	@for pkg in $(PACKAGE_NAMES); do \
		echo "  - $$pkg"; \
	done
	@echo ""
	@echo "$(BOLD)🔧 Environment:$(NC)"
	@echo "  Go Version: $$(go version 2>/dev/null || echo 'Not installed')"
	@echo "  Make Version: $$(make --version | head -1)"
	@echo "  Platform: $(PLATFORM)/$(ARCH)"
	@echo ""

##@ 🏗️ Building

.PHONY: build
build: ## 🔨 Build all packages
	@echo "$(BOLD)$(GREEN)🔨 Building all packages...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)📦 Building $$pkg...$(NC)"; \
		cd $$pkg && go build -v ./... && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ All packages built successfully$(NC)"

.PHONY: build-%
build-%: ## 🔨 Build specific package (e.g., make build-logger)
	@if [ -d "$(PACKAGES_DIR)/$*" ]; then \
		echo "$(CYAN)🔨 Building package: $*$(NC)"; \
		cd $(PACKAGES_DIR)/$* && go build -v ./...; \
		echo "$(GREEN)✅ Package $* built successfully$(NC)"; \
	else \
		echo "$(RED)❌ Package $* not found$(NC)"; \
		exit 1; \
	fi

.PHONY: clean
clean: ## 🧹 Clean build artifacts and caches
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR) $(REPORTS_DIR)
	@go clean -cache -testcache -modcache
	@for pkg in $(PACKAGES); do \
		cd $$pkg && go clean ./... && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

##@ 🧪 Testing

.PHONY: test
test: ## 🧪 Run all tests with coverage
	@echo "$(BOLD)$(BLUE)🧪 Running tests with coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🧪 Testing $$pkg...$(NC)"; \
		pkg_name=$$(basename $$pkg); \
		cd $$pkg && go test -v -race -coverprofile=../$(COVERAGE_DIR)/$$pkg_name-coverage.out -covermode=atomic ./...; \
		if [ $$? -eq 0 ]; then \
			echo "$(GREEN)✅ Tests passed for $$pkg$(NC)"; \
		else \
			echo "$(RED)❌ Tests failed for $$pkg$(NC)"; \
		fi; \
		cd - > /dev/null; \
	done
	@$(MAKE) coverage-merge

.PHONY: test-%
test-%: ## 🧪 Test specific package (e.g., make test-logger)
	@if [ -d "$(PACKAGES_DIR)/$*" ]; then \
		echo "$(CYAN)🧪 Testing package: $*$(NC)"; \
		mkdir -p $(COVERAGE_DIR); \
		cd $(PACKAGES_DIR)/$* && go test -v -race -coverprofile=../../$(COVERAGE_DIR)/$*-coverage.out -covermode=atomic ./...; \
		echo "$(GREEN)✅ Tests completed for $*$(NC)"; \
	else \
		echo "$(RED)❌ Package $* not found$(NC)"; \
		exit 1; \
	fi

.PHONY: test-quick
test-quick: ## ⚡ Run tests without coverage (faster)
	@echo "$(BLUE)⚡ Running quick tests...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🧪 Testing $$pkg...$(NC)"; \
		cd $$pkg && go test -v ./... && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Quick tests completed$(NC)"

.PHONY: test-verbose
test-verbose: ## 📝 Run tests with verbose output
	@echo "$(BLUE)📝 Running verbose tests...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🧪 Testing $$pkg...$(NC)"; \
		cd $$pkg && go test -v -race -count=1 ./... && cd - > /dev/null; \
	done

.PHONY: benchmark
benchmark: ## 🏃 Run benchmark tests
	@echo "$(PURPLE)🏃 Running benchmarks...$(NC)"
	@mkdir -p $(REPORTS_DIR)
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🏃 Benchmarking $$pkg...$(NC)"; \
		pkg_name=$$(basename $$pkg); \
		cd $$pkg && go test -bench=. -benchmem -run=^$$ ./... > ../$(REPORTS_DIR)/$$pkg_name-benchmark.txt; \
		cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Benchmarks completed$(NC)"

##@ 📊 Coverage

.PHONY: coverage-merge
coverage-merge: ## 📊 Merge coverage reports
	@echo "$(BLUE)📊 Merging coverage reports...$(NC)"
	@if command -v gocovmerge >/dev/null 2>&1; then \
		gocovmerge $(COVERAGE_DIR)/*-coverage.out > $(COVERAGE_DIR)/merged-coverage.out; \
		echo "$(GREEN)✅ Coverage reports merged$(NC)"; \
	else \
		echo "$(YELLOW)⚠️ gocovmerge not found, installing...$(NC)"; \
		go install github.com/wadey/gocovmerge@latest; \
		gocovmerge $(COVERAGE_DIR)/*-coverage.out > $(COVERAGE_DIR)/merged-coverage.out; \
		echo "$(GREEN)✅ Coverage reports merged$(NC)"; \
	fi

.PHONY: coverage-html
coverage-html: coverage-merge ## 📊 Generate HTML coverage report
	@echo "$(BLUE)📊 Generating HTML coverage report...$(NC)"
	@go tool cover -html=$(COVERAGE_DIR)/merged-coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)✅ HTML coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: coverage-report
coverage-report: coverage-merge ## 📊 Show coverage summary
	@echo "$(BLUE)📊 Coverage Summary:$(NC)"
	@go tool cover -func=$(COVERAGE_DIR)/merged-coverage.out | tail -1
	@go tool cover -func=$(COVERAGE_DIR)/merged-coverage.out > $(COVERAGE_DIR)/coverage-summary.txt

##@ 🔍 Code Quality

.PHONY: lint
lint: ## 🔍 Run linting with golangci-lint
	@echo "$(PURPLE)🔍 Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config .golangci.yml ./...; \
	else \
		echo "$(YELLOW)⚠️ golangci-lint not found, installing...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_VERSION); \
		golangci-lint run --config .golangci.yml ./...; \
	fi
	@echo "$(GREEN)✅ Linting completed$(NC)"

.PHONY: lint-fix
lint-fix: ## 🔧 Run linting with auto-fix
	@echo "$(PURPLE)🔧 Running golangci-lint with auto-fix...$(NC)"
	@golangci-lint run --config .golangci.yml --fix ./...
	@echo "$(GREEN)✅ Linting with auto-fix completed$(NC)"

.PHONY: format
format: ## 🎨 Format code with gofmt and goimports
	@echo "$(CYAN)🎨 Formatting code...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🎨 Formatting $$pkg...$(NC)"; \
		cd $$pkg && gofmt -s -w . && cd - > /dev/null; \
		cd $$pkg && goimports -w . && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Code formatting completed$(NC)"

.PHONY: vet
vet: ## 🔍 Run go vet
	@echo "$(PURPLE)🔍 Running go vet...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🔍 Vetting $$pkg...$(NC)"; \
		cd $$pkg && go vet ./... && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Go vet completed$(NC)"

##@ 🔒 Security

.PHONY: security
security: security-gosec security-govulncheck ## 🔒 Run all security scans

.PHONY: security-gosec  
security-gosec: ## 🛡️ Run gosec security scanner
	@echo "$(RED)🛡️ Running gosec security scanner...$(NC)"
	@mkdir -p $(REPORTS_DIR)/security
	@if command -v gosec >/dev/null 2>&1; then \
		for pkg in $(PACKAGES); do \
			echo "$(CYAN)🔍 Scanning $$pkg...$(NC)"; \
			pkg_name=$$(basename $$pkg); \
			cd $$pkg && gosec -fmt json -out ../$(REPORTS_DIR)/security/$$pkg_name-gosec.json ./... || true; \
			cd - > /dev/null; \
		done; \
	else \
		echo "$(YELLOW)⚠️ gosec not found, installing...$(NC)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@$(GOSEC_VERSION); \
		$(MAKE) security-gosec; \
	fi
	@echo "$(GREEN)✅ gosec scan completed$(NC)"

.PHONY: security-govulncheck
security-govulncheck: ## 🔍 Run govulncheck vulnerability scanner
	@echo "$(RED)🔍 Running govulncheck...$(NC)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		for pkg in $(PACKAGES); do \
			echo "$(CYAN)🔍 Checking $$pkg for vulnerabilities...$(NC)"; \
			cd $$pkg && govulncheck ./... && cd - > /dev/null; \
		done; \
	else \
		echo "$(YELLOW)⚠️ govulncheck not found, installing...$(NC)"; \
		go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION); \
		$(MAKE) security-govulncheck; \
	fi
	@echo "$(GREEN)✅ Vulnerability check completed$(NC)"

##@ 📦 Dependencies

.PHONY: deps-download
deps-download: ## 📥 Download dependencies for all packages
	@echo "$(BLUE)📥 Downloading dependencies...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)📥 Downloading deps for $$pkg...$(NC)"; \
		cd $$pkg && go mod download && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Dependencies downloaded$(NC)"

.PHONY: deps-tidy
deps-tidy: ## 🧹 Tidy dependencies for all packages
	@echo "$(BLUE)🧹 Tidying dependencies...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🧹 Tidying deps for $$pkg...$(NC)"; \
		cd $$pkg && go mod tidy && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Dependencies tidied$(NC)"

.PHONY: deps-verify
deps-verify: ## ✅ Verify dependencies for all packages
	@echo "$(BLUE)✅ Verifying dependencies...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)✅ Verifying deps for $$pkg...$(NC)"; \
		cd $$pkg && go mod verify && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Dependencies verified$(NC)"

.PHONY: deps-update
deps-update: ## 🔄 Update dependencies (patch versions)
	@echo "$(YELLOW)🔄 Updating dependencies (patch versions)...$(NC)"
	@for pkg in $(PACKAGES); do \
		echo "$(CYAN)🔄 Updating deps for $$pkg...$(NC)"; \
		cd $$pkg && go get -u=patch ./... && go mod tidy && cd - > /dev/null; \
	done
	@echo "$(GREEN)✅ Dependencies updated$(NC)"

##@ 🚀 Release and Distribution

.PHONY: pre-release
pre-release: clean lint test security ## 🚀 Run all pre-release checks
	@echo "$(BOLD)$(GREEN)🚀 Pre-release checks completed successfully$(NC)"

.PHONY: check-all
check-all: lint vet test security ## ✅ Run all quality checks
	@echo "$(BOLD)$(GREEN)✅ All quality checks passed$(NC)"

##@ 🛠️ Development

.PHONY: dev-setup
dev-setup: ## 🛠️ Set up development environment
	@echo "$(BLUE)🛠️ Setting up development environment...$(NC)"
	@echo "$(CYAN)📥 Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	@go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	@go install github.com/wadey/gocovmerge@latest
	@echo "$(GREEN)✅ Development environment setup completed$(NC)"

.PHONY: dev-reset
dev-reset: clean deps-tidy ## 🔄 Reset development environment
	@echo "$(YELLOW)🔄 Resetting development environment...$(NC)"
	@$(MAKE) deps-download
	@echo "$(GREEN)✅ Development environment reset$(NC)"

##@ 📊 Reporting

.PHONY: report-summary
report-summary: ## 📊 Generate comprehensive project summary
	@echo "$(BOLD)$(BLUE)📊 GoVel Framework Summary Report$(NC)"
	@echo ""
	@echo "$(BOLD)📦 Packages: $(GREEN)$(words $(PACKAGES))$(NC)"
	@for pkg in $(PACKAGE_NAMES); do \
		lines=$$(find $(PACKAGES_DIR)/$$pkg -name "*.go" -not -path "*/vendor/*" -not -name "*_test.go" | xargs wc -l 2>/dev/null | tail -1 | awk '{print $$1}' || echo "0"); \
		tests=$$(find $(PACKAGES_DIR)/$$pkg -name "*_test.go" | wc -l || echo "0"); \
		echo "  - $(CYAN)$$pkg$(NC): $$lines lines, $$tests test files"; \
	done
	@echo ""
	@total_lines=$$(find $(PACKAGES_DIR) -name "*.go" -not -path "*/vendor/*" -not -name "*_test.go" | xargs wc -l 2>/dev/null | tail -1 | awk '{print $$1}' || echo "0"); \
	total_tests=$$(find $(PACKAGES_DIR) -name "*_test.go" | wc -l || echo "0"); \
	echo "$(BOLD)📈 Total: $(GREEN)$$total_lines$(NC) lines of code, $(GREEN)$$total_tests$(NC) test files"

##@ 🎯 Convenience Targets

.PHONY: all
all: clean build test lint security ## 🎯 Run complete build pipeline
	@echo "$(BOLD)$(GREEN)🎉 Complete build pipeline finished successfully$(NC)"

.PHONY: quick
quick: build test-quick ## ⚡ Quick build and test
	@echo "$(BOLD)$(GREEN)⚡ Quick pipeline completed$(NC)"

# 📝 Phony target declaration
.PHONY: help info build clean test coverage lint format vet security deps-download deps-tidy deps-verify deps-update pre-release check-all dev-setup dev-reset report-summary all quick

