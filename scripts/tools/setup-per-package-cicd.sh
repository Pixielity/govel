#!/bin/bash

##
# 📦 GoVel Framework - Per-Package CI/CD Setup Script
#
# This script creates individual CI/CD configurations for each package
# in the GoVel framework. It generates GitHub Actions workflows, linting
# configurations, and other necessary files for each package.
#
# Features:
# - 🎯 Individual CI workflows per package
# - 🔧 Package-specific configurations
# - 📊 Code coverage per package
# - 🔒 Security scanning per package
# - 🛠️ Build automation per package
#
# Usage:
#   ./setup-per-package-cicd.sh [options]
#
# Options:
#   --dry-run, -d     Show what would be created without creating files
#   --help, -h        Show this help message  
#   --verbose, -v     Enable verbose output
#   --force, -f       Overwrite existing files without confirmation
#
# Author: GoVel Framework Team
# Version: 1.0.0
# License: MIT
##

set -euo pipefail

# 🎨 Color definitions for beautiful output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly BOLD='\033[1m'
readonly NC='\033[0m'

# 📝 Script configuration
readonly SCRIPT_NAME="$(basename "$0")"
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "$SCRIPT_DIR" && pwd)"

# 🔧 Default settings
DRY_RUN=false
VERBOSE=false
FORCE=false

##
# 📝 Print formatted message with emoji and color support
##
print_message() {
    local type="$1"
    local message="$2"
    local emoji="${3:-}"
    
    case "$type" in
        "info")
            echo -e "${BLUE}${emoji:-ℹ️}  INFO:${NC} $message"
            ;;
        "success")
            echo -e "${GREEN}${emoji:-✅} SUCCESS:${NC} $message"
            ;;
        "warning")
            echo -e "${YELLOW}${emoji:-⚠️}  WARNING:${NC} $message"
            ;;
        "error")
            echo -e "${RED}${emoji:-❌} ERROR:${NC} $message" >&2
            ;;
        "debug")
            if [[ "$VERBOSE" == "true" ]]; then
                echo -e "${PURPLE}${emoji:-🔍} DEBUG:${NC} $message"
            fi
            ;;
        "header")
            echo -e "\n${BOLD}${CYAN}${emoji:-🚀} $message${NC}"
            echo -e "${CYAN}$(printf '=%.0s' {1..50})${NC}"
            ;;
    esac
}

##
# 🛠️ Display help information
##
show_help() {
    cat << EOF
${BOLD}${CYAN}📦 GoVel Per-Package CI/CD Setup${NC}

${BOLD}DESCRIPTION:${NC}
    Creates individual CI/CD configurations for each GoVel package.
    Generates GitHub Actions workflows, configuration files, and automation scripts.

${BOLD}USAGE:${NC}
    $SCRIPT_NAME [OPTIONS]

${BOLD}OPTIONS:${NC}
    -d, --dry-run         🔍 Show what would be created without creating files
    -v, --verbose         📝 Enable detailed output and debug information
    -f, --force           💪 Overwrite existing files without confirmation
    -h, --help            ❓ Show this help message and exit

${BOLD}EXAMPLES:${NC}
    $SCRIPT_NAME --dry-run           # Preview what will be created
    $SCRIPT_NAME                     # Interactive setup with confirmations
    $SCRIPT_NAME --force --verbose   # Force setup with detailed output

EOF
}

##
# 🎯 Parse command line arguments
##
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dry-run)
                DRY_RUN=true
                print_message "info" "Dry-run mode enabled - no files will be created" "🔍"
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                print_message "info" "Verbose mode enabled" "📝"
                shift
                ;;
            -f|--force)
                FORCE=true
                print_message "info" "Force mode enabled - existing files will be overwritten" "💪"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_message "error" "Unknown option: $1"
                print_message "info" "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

##
# 🔍 Find all packages with go.mod files
##
find_packages() {
    print_message "info" "Discovering GoVel packages..." "🔍" >&2
    
    local packages=()
    while IFS= read -r -d '' mod_file; do
        local pkg_dir
        pkg_dir=$(dirname "$mod_file")
        # Remove /src suffix if present
        pkg_dir=${pkg_dir%/src}
        packages+=("$pkg_dir")
    done < <(find "$PROJECT_ROOT/packages" -name "go.mod" -type f -print0)
    
    if [[ ${#packages[@]} -eq 0 ]]; then
        print_message "warning" "No packages found with go.mod files" >&2
        return 1
    fi
    
    print_message "success" "Found ${#packages[@]} packages:" "📦" >&2
    for pkg in "${packages[@]}"; do
        echo "  📁 $pkg" >&2
    done
    
    # Return packages as newline-separated string to avoid word splitting issues
    printf '%s\n' "${packages[@]}"
}

##
# 📄 Generate package-specific GitHub workflow
##
generate_package_workflow() {
    local pkg_path="$1"
    local pkg_name
    pkg_name=$(basename "$pkg_path")
    
    print_message "info" "Generating workflow for $pkg_name..." "⚙️"
    
    local workflow_dir="$pkg_path/.github/workflows"
    local workflow_file="$workflow_dir/ci.yml"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: $workflow_file" "📄"
        return 0
    fi
    
    # Create directory structure
    mkdir -p "$workflow_dir"
    
    # Check if file exists and handle overwrite logic
    if [[ -f "$workflow_file" && "$FORCE" != "true" ]]; then
        print_message "warning" "Workflow already exists: $workflow_file"
        echo -n "❓ Overwrite? [y/N]: "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_message "info" "Skipping $pkg_name workflow"
            return 0
        fi
    fi
    
    # Generate workflow content
    cat > "$workflow_file" << EOF
# 🎯 $pkg_name Package CI Pipeline
#
# This workflow provides continuous integration for the $pkg_name package.
# It runs tests, linting, security scans, and quality checks.

name: 🧪 $pkg_name CI

on:
  push:
    branches: [ main, develop ]
    paths:
      - '$pkg_path/**'
      - '.github/workflows/$(basename "$workflow_file")'
  
  pull_request:
    branches: [ main, develop ]
    paths:
      - '$pkg_path/**'
      - '.github/workflows/$(basename "$workflow_file")'
  
  workflow_dispatch:

# 🔒 Security: Ensure minimal permissions
permissions:
  contents: read
  security-events: write
  pull-requests: write

# 🚫 Cancel in-progress runs for same PR/branch
concurrency:
  group: \${{ github.workflow }}-\${{ github.ref }}
  cancel-in-progress: true

env:
  PACKAGE_PATH: $pkg_path
  GO_VERSION: '1.23'
  COVERAGE_THRESHOLD: 80

jobs:
  test:
    name: 🧪 Test $pkg_name
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: \${{ env.GO_VERSION }}
          cache: true

      - name: 📦 Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/go/pkg/mod
            ~/.cache/go-build
          key: \${{ runner.os }}-go-\${{ env.GO_VERSION }}-\${{ hashFiles('$pkg_path/**/go.sum', '$pkg_path/**/go.mod') }}

      - name: 📥 Download dependencies
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "📥 Downloading dependencies for $pkg_name..."
          # Handle both src and non-src directory structures
          if [[ -f "src/go.mod" ]]; then
            cd src && go mod download
          else
            go mod download
          fi

      - name: 🔨 Build package
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "🔨 Building $pkg_name..."
          if [[ -f "src/go.mod" ]]; then
            cd src && go build -v ./...
          else
            go build -v ./...
          fi

      - name: 🧪 Run tests
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "🧪 Running tests for $pkg_name..."
          mkdir -p coverage
          if [[ -f "src/go.mod" ]]; then
            cd src && go test -v -race -coverprofile=../coverage/coverage.out -covermode=atomic ./...
            cd .. && go tool cover -func=coverage/coverage.out
          else
            go test -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./...
            go tool cover -func=coverage/coverage.out
          fi

      - name: 🔍 Run go vet
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "🔍 Running go vet for $pkg_name..."
          if [[ -f "src/go.mod" ]]; then
            cd src && go vet ./...
          else
            go vet ./...
          fi

      - name: 🎨 Check formatting
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "🎨 Checking formatting for $pkg_name..."
          if [[ -f "src/go.mod" ]]; then
            unformatted=\$(cd src && gofmt -l .)
          else
            unformatted=\$(gofmt -l .)
          fi
          
          if [[ -n "\$unformatted" ]]; then
            echo "❌ Unformatted files:"
            echo "\$unformatted"
            exit 1
          fi

      - name: 📊 Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./\${{ env.PACKAGE_PATH }}/coverage/coverage.out
          flags: $pkg_name
          name: $pkg_name-coverage
          fail_ci_if_error: false

  security:
    name: 🔒 Security Scan
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: \${{ env.GO_VERSION }}

      - name: 🔍 Run govulncheck
        working-directory: \${{ env.PACKAGE_PATH }}
        run: |
          echo "🔍 Installing govulncheck..."
          go install golang.org/x/vuln/cmd/govulncheck@latest
          
          echo "🔍 Scanning $pkg_name for vulnerabilities..."
          if [[ -f "src/go.mod" ]]; then
            cd src && govulncheck ./...
          else
            govulncheck ./...
          fi

      - name: 🔒 Run gosec
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: '-fmt sarif -out gosec-$pkg_name.sarif \${{ env.PACKAGE_PATH }}/...'

      - name: 🛡️ Upload SARIF file
        if: always()
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: gosec-$pkg_name.sarif
EOF

    print_message "success" "Created workflow for $pkg_name" "✅"
}

##
# 📄 Generate package-specific configurations
##
generate_package_configs() {
    local pkg_path="$1"
    local pkg_name
    pkg_name=$(basename "$pkg_path")
    
    print_message "info" "Generating configurations for $pkg_name..." "🔧"
    
    # Generate package-specific .golangci.yml
    local golangci_file="$pkg_path/.golangci.yml"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: $golangci_file" "📄"
        return 0
    fi
    
    if [[ -f "$golangci_file" && "$FORCE" != "true" ]]; then
        print_message "debug" "Config already exists: $golangci_file" "🔍"
        return 0
    fi
    
    cat > "$golangci_file" << 'EOF'
# 🔧 GolangCI-Lint Configuration for Package
# This configuration provides comprehensive linting for this specific package

run:
  timeout: 5m
  tests: true
  modules-download-mode: readonly

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters:
  enable:
    - errcheck      # Check for unchecked errors
    - gosimple      # Simplify code
    - govet         # Go vet
    - ineffassign   # Detect ineffectual assignments
    - staticcheck   # Static analysis
    - typecheck     # Type checking
    - unused        # Unused code
    - gocritic      # Go critic
    - gocyclo       # Cyclomatic complexity
    - gofmt         # Format checking
    - goimports     # Import checking
    - gosec         # Security issues
    - misspell      # Misspellings
    - revive        # Revive linter
    - unconvert     # Unnecessary conversions

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  gocyclo:
    min-complexity: 15
  
  gosec:
    excludes:
      - G104 # Audit errors not checked (too noisy for some packages)
  
  revive:
    rules:
      - name: var-naming
      - name: package-comments
      - name: exported
      - name: error-return
      - name: error-naming
      - name: if-return

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
  exclude-use-default: false
EOF

    print_message "success" "Created config for $pkg_name" "✅"
}

##
# 📄 Generate package-specific Makefile
##
generate_package_makefile() {
    local pkg_path="$1"
    local pkg_name
    pkg_name=$(basename "$pkg_path")
    
    print_message "info" "Generating Makefile for $pkg_name..." "🛠️"
    
    local makefile="$pkg_path/Makefile"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: $makefile" "📄"
        return 0
    fi
    
    if [[ -f "$makefile" && "$FORCE" != "true" ]]; then
        print_message "debug" "Makefile already exists: $makefile" "🔍"
        return 0
    fi
    
    cat > "$makefile" << EOF
# 🛠️ Makefile for $pkg_name Package
# Provides build automation and development tools

.PHONY: help build test clean lint fmt vet security coverage dev-setup

# Default target
help: ## 📚 Show this help message
	@echo "🛠️  Available targets for $pkg_name:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*\$\$' \$(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\\n", \$\$1, \$\$2}'

# Determine if we're in src structure or not
WORKING_DIR := \$(shell if [ -f "src/go.mod" ]; then echo "src"; else echo "."; fi)

build: ## 🔨 Build the package
	@echo "🔨 Building $pkg_name..."
	@cd \$(WORKING_DIR) && go build -v ./...

test: ## 🧪 Run tests with coverage
	@echo "🧪 Running tests for $pkg_name..."
	@mkdir -p coverage
	@cd \$(WORKING_DIR) && go test -v -race -coverprofile=../coverage/coverage.out -covermode=atomic ./...
	@cd \$(WORKING_DIR) && go tool cover -func=../coverage/coverage.out

test-short: ## ⚡ Run tests without coverage
	@echo "⚡ Running quick tests for $pkg_name..."
	@cd \$(WORKING_DIR) && go test -v -short ./...

clean: ## 🧹 Clean build artifacts and cache
	@echo "🧹 Cleaning $pkg_name..."
	@cd \$(WORKING_DIR) && go clean -cache -testcache -modcache
	@rm -rf coverage/

lint: ## 🔍 Run golangci-lint
	@echo "🔍 Linting $pkg_name..."
	@if command -v golangci-lint >/dev/null 2>&1; then \\
		cd \$(WORKING_DIR) && golangci-lint run; \\
	else \\
		echo "⚠️  golangci-lint not installed. Run 'make dev-setup' first."; \\
	fi

fmt: ## 🎨 Format code
	@echo "🎨 Formatting $pkg_name..."
	@cd \$(WORKING_DIR) && go fmt ./...

vet: ## 🔍 Run go vet
	@echo "🔍 Running go vet for $pkg_name..."
	@cd \$(WORKING_DIR) && go vet ./...

security: ## 🔒 Run security scans
	@echo "🔒 Running security scans for $pkg_name..."
	@if command -v gosec >/dev/null 2>&1; then \\
		cd \$(WORKING_DIR) && gosec ./...; \\
	else \\
		echo "⚠️  gosec not installed. Run 'make dev-setup' first."; \\
	fi
	@if command -v govulncheck >/dev/null 2>&1; then \\
		cd \$(WORKING_DIR) && govulncheck ./...; \\
	else \\
		echo "⚠️  govulncheck not installed. Run 'make dev-setup' first."; \\
	fi

coverage: test ## 📊 Generate coverage report
	@echo "📊 Generating coverage report for $pkg_name..."
	@cd \$(WORKING_DIR) && go tool cover -html=../coverage/coverage.out -o ../coverage/coverage.html
	@echo "📊 Coverage report: coverage/coverage.html"

dev-setup: ## 🛠️ Install development tools
	@echo "🛠️ Installing development tools for $pkg_name..."
	@cd \$(WORKING_DIR) && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@cd \$(WORKING_DIR) && go install github.com/securecodewarrior/gosec/cmd/gosec@latest
	@cd \$(WORKING_DIR) && go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development tools installed"

all: clean fmt vet lint test security ## 🚀 Run complete CI pipeline
	@echo "🎉 All checks passed for $pkg_name!"

# Dependencies
deps: ## 📥 Download dependencies
	@echo "📥 Downloading dependencies for $pkg_name..."
	@cd \$(WORKING_DIR) && go mod download
	@cd \$(WORKING_DIR) && go mod tidy

update-deps: ## ⬆️ Update dependencies
	@echo "⬆️ Updating dependencies for $pkg_name..."
	@cd \$(WORKING_DIR) && go get -u ./...
	@cd \$(WORKING_DIR) && go mod tidy
EOF

    print_message "success" "Created Makefile for $pkg_name" "✅"
}

##
# 🎯 Main execution function
##
main() {
    print_message "header" "GoVel Per-Package CI/CD Setup" "📦"
    
    # Parse command line arguments
    parse_arguments "$@"
    
    # Find all packages and read into array
    local packages=()
    while IFS= read -r pkg; do
        [[ -n "$pkg" ]] && packages+=("$pkg")
    done < <(find_packages)
    
    print_message "header" "Setup Preview" "🔍"
    echo ""
    echo -e "${BOLD}The following will be created for each package:${NC}"
    echo "  📄 .github/workflows/ci.yml    # Individual CI pipeline"
    echo "  📄 .golangci.yml               # Package-specific linting"
    echo "  📄 Makefile                    # Build automation"
    echo ""
    
    if [[ "$DRY_RUN" != "true" ]]; then
        echo -n "❓ Do you want to proceed with per-package CI/CD setup? [y/N]: "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_message "info" "Setup cancelled by user"
            exit 0
        fi
        print_message "success" "Proceeding with setup..." "✅"
    fi
    
    # Process each package
    print_message "header" "Processing Packages" "⚙️"
    
    local created_count=0
    for pkg in "${packages[@]}"; do
        echo ""
        print_message "info" "Processing $(basename "$pkg")..." "📦"
        
        # Generate workflow
        generate_package_workflow "$pkg"
        
        # Generate configurations
        generate_package_configs "$pkg"
        
        # Generate Makefile
        generate_package_makefile "$pkg"
        
        ((created_count++))
    done
    
    # Summary
    print_message "header" "Setup Complete" "🎉"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "success" "Dry run completed - showed what would be created for $created_count packages" "✅"
        print_message "info" "Run without --dry-run to create the files"
    else
        print_message "success" "Per-package CI/CD setup completed for $created_count packages!" "🎉"
        echo ""
        echo -e "${BOLD}📚 Next Steps:${NC}"
        echo "  1. 📝 Review the generated files in each package"
        echo "  2. 🔧 Customize configurations as needed"
        echo "  3. 🛠️ Run 'make dev-setup' in each package to install tools"
        echo "  4. 🧪 Run 'make test' in each package to verify setup"
        echo "  5. 📤 Commit and push the changes to GitHub"
        echo ""
        echo -e "${BOLD}🛠️  Package Commands:${NC}"
        echo "  make help          # Show available commands"
        echo "  make all           # Run complete CI pipeline"
        echo "  make test          # Run tests with coverage"
        echo "  make lint          # Run linting"
        echo "  make security      # Run security scans"
    fi
}

# Execute main function
main "$@"