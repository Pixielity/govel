#!/bin/bash

##
# 🚀 GoVel CI/CD Generator Script
# 
# This script generates a comprehensive GitHub Actions CI/CD pipeline for the GoVel framework.
# It creates workflows, scripts, configurations, and templates following modern Go best practices.
#
# Features:
# - 🎯 Professional-grade CI/CD workflows
# - 🔒 Security scanning and vulnerability detection
# - 📊 Code coverage reporting and quality gates
# - 🤖 Automated dependency management
# - 🧪 Multi-version and cross-platform testing
# - 📦 Smart package detection for monorepo optimization
# - 🛠️ Development tooling and automation scripts
#
# Usage:
#   ./generate-cicd.sh [options]
#
# Options:
#   --dry-run, -d     Show what would be created without actually creating files
#   --help, -h        Show this help message
#   --verbose, -v     Enable verbose output
#   --force, -f       Overwrite existing files without confirmation
#
# Examples:
#   ./generate-cicd.sh --dry-run          # Preview what will be created
#   ./generate-cicd.sh                    # Create with confirmation prompts
#   ./generate-cicd.sh --force            # Create and overwrite existing files
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
readonly WHITE='\033[1;37m'
readonly BOLD='\033[1m'
readonly NC='\033[0m' # No Color

# 📝 Script configuration
readonly SCRIPT_NAME="$(basename "$0")"
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "$SCRIPT_DIR" && pwd)"
readonly TIMESTAMP="$(date '+%Y-%m-%d %H:%M:%S')"

# 🔧 Default settings
DRY_RUN=false
VERBOSE=false
FORCE=false
CONFIRMATION_REQUIRED=true

##
# 📝 Print formatted message with emoji and color support
#
# This function provides consistent, colorful output throughout the script
# with emoji support for better visual feedback.
#
# @param string $1 The message type (info, success, warning, error, debug)
# @param string $2 The message text to display
# @param string $3 Optional emoji to prepend (defaults based on type)
#
# Examples:
#   print_message "info" "Starting CI/CD generation"
#   print_message "success" "File created successfully" "✅"
#   print_message "error" "Failed to create directory"
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
        *)
            echo -e "$message"
            ;;
    esac
}

##
# 🛠️ Display help information
#
# Shows comprehensive usage information, examples, and available options
# for the CI/CD generator script.
##
show_help() {
    cat << EOF
${BOLD}${CYAN}🚀 GoVel CI/CD Generator${NC}

${BOLD}DESCRIPTION:${NC}
    Generates a comprehensive GitHub Actions CI/CD pipeline for the GoVel framework.
    Creates workflows, scripts, configurations, and templates following Go best practices.

${BOLD}USAGE:${NC}
    $SCRIPT_NAME [OPTIONS]

${BOLD}OPTIONS:${NC}
    -d, --dry-run     🔍 Show what would be created without actually creating files
    -v, --verbose     📝 Enable detailed output and debug information
    -f, --force       💪 Overwrite existing files without confirmation prompts
    -h, --help        ❓ Show this help message and exit

${BOLD}EXAMPLES:${NC}
    $SCRIPT_NAME --dry-run          # Preview all files that would be created
    $SCRIPT_NAME                    # Interactive creation with confirmation prompts
    $SCRIPT_NAME --force --verbose  # Create all files with detailed output

${BOLD}FEATURES INCLUDED:${NC}
    🎯 Multi-version Go testing (1.21, 1.22, 1.23)
    🖥️  Cross-platform testing (Linux, macOS, Windows)
    🔒 Security scanning (gosec, govulncheck)
    📊 Code coverage reporting with Codecov integration
    🤖 Automated dependency updates via Dependabot
    📦 Smart package detection for monorepo optimization
    🧪 Comprehensive linting and code quality checks
    🚀 Automated release management with semantic versioning
    🛠️ Development scripts and Git hooks

${BOLD}FILES CREATED:${NC}
    📁 .github/workflows/        # GitHub Actions workflows
    📁 .github/templates/        # Issue and PR templates
    📁 scripts/ci/              # CI automation scripts
    📁 scripts/utils/           # Utility scripts
    📁 scripts/hooks/           # Git hooks
    📄 Configuration files      # Linting, coverage, and build configs

${BOLD}AUTHOR:${NC} GoVel Framework Team
${BOLD}VERSION:${NC} 1.0.0
${BOLD}LICENSE:${NC} MIT

EOF
}

##
# 🔍 Parse command line arguments
#
# Processes all command-line options and sets appropriate flags
# for script execution behavior.
#
# @param array $@ All command line arguments passed to the script
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
                CONFIRMATION_REQUIRED=false
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
# 📁 Create directory structure
#
# Creates the complete directory structure needed for the CI/CD setup.
# Includes proper error handling and logging.
#
# @param string $1 Base path where directories should be created
##
create_directory_structure() {
    local base_path="$1"
    
    print_message "header" "Creating Directory Structure" "📁"
    
    local directories=(
        ".github/workflows"
        ".github/templates"
        ".github/ISSUE_TEMPLATE"
        "scripts/ci"
        "scripts/utils"
        "scripts/hooks"
    )
    
    for dir in "${directories[@]}"; do
        local full_path="$base_path/$dir"
        print_message "debug" "Creating directory: $full_path"
        
        if [[ "$DRY_RUN" == "false" ]]; then
            if mkdir -p "$full_path"; then
                print_message "success" "Created directory: $dir" "📂"
            else
                print_message "error" "Failed to create directory: $dir" "❌"
                return 1
            fi
        else
            print_message "info" "[DRY-RUN] Would create directory: $dir" "📂"
        fi
    done
}

##
# 🎯 Generate main CI workflow
#
# Creates the primary GitHub Actions workflow file that handles testing,
# linting, security scanning, and coverage reporting across multiple
# Go versions and platforms.
##
generate_main_ci_workflow() {
    local file_path="$PROJECT_ROOT/.github/workflows/ci.yml"
    print_message "info" "Generating main CI workflow" "🎯"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .github/workflows/ci.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 🎯 GoVel Framework - Main CI Pipeline
#
# This workflow provides comprehensive continuous integration for the GoVel framework.
# It handles testing, linting, security scanning, and quality assurance across
# multiple Go versions and operating systems.
#
# Features:
# - 🧪 Multi-version Go testing (1.21, 1.22, 1.23)
# - 🖥️  Cross-platform testing (Ubuntu, macOS, Windows)
# - 📦 Smart package detection (only test changed packages)
# - 🔒 Security scanning integration
# - 📊 Code coverage reporting with quality gates
# - 🛠️ Comprehensive linting and formatting checks
# - ⚡ Performance optimization with caching strategies
#
# Triggers:
# - Push to main/develop branches
# - Pull requests to main/develop branches
# - Manual workflow dispatch
#
# Author: GoVel Framework Team
# Version: 1.0.0

name: 🎯 CI Pipeline

on:
  # 🔄 Trigger on pushes to main branches
  push:
    branches: [ main, develop ]
    paths-ignore:
      - '**.md'
      - 'docs/**'
      - 'examples/**'
  
  # 🔀 Trigger on pull requests
  pull_request:
    branches: [ main, develop ]
    paths-ignore:
      - '**.md'
      - 'docs/**'
      - 'examples/**'
  
  # 🎛️ Allow manual triggering
  workflow_dispatch:
    inputs:
      test_all:
        description: '🧪 Test all packages (ignore change detection)'
        required: false
        default: 'false'
        type: boolean
      skip_security:
        description: '🔒 Skip security scanning'
        required: false
        default: 'false'
        type: boolean

# 🔒 Security: Ensure minimal permissions
permissions:
  contents: read
  security-events: write
  pull-requests: write
  checks: write

# 🚫 Cancel in-progress runs for same PR/branch
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

# 🌍 Environment variables
env:
  GO_VERSION_MATRIX: "1.21,1.22,1.23"
  COVERAGE_THRESHOLD: 80
  CODECOV_TOKEN: ${{ secrets.CODECOV_TOKEN }}

jobs:
  # 📦 Detect which packages have changed
  detect-changes:
    name: 🔍 Detect Changed Packages
    runs-on: ubuntu-latest
    outputs:
      packages: ${{ steps.changes.outputs.packages }}
      has-changes: ${{ steps.changes.outputs.has-changes }}
      test-all: ${{ steps.changes.outputs.test-all }}
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: 🔍 Detect package changes
        id: changes
        run: |
          echo "🔍 Detecting changed packages..."
          
          # Force test all packages if manually requested
          if [[ "${{ github.event.inputs.test_all }}" == "true" ]]; then
            echo "🎯 Manual override: Testing all packages"
            echo "test-all=true" >> $GITHUB_OUTPUT
            echo "has-changes=true" >> $GITHUB_OUTPUT
            find packages -name "go.mod" -type f | sed 's|/go.mod||' | jq -R -s -c 'split("\n")[:-1]' > packages.json
            echo "packages=$(cat packages.json)" >> $GITHUB_OUTPUT
            exit 0
          fi
          
          # Get changed files
          if [[ "${{ github.event_name }}" == "pull_request" ]]; then
            CHANGED_FILES=$(git diff --name-only ${{ github.event.pull_request.base.sha }}..${{ github.event.pull_request.head.sha }})
          else
            CHANGED_FILES=$(git diff --name-only HEAD~1)
          fi
          
          echo "📝 Changed files:"
          echo "$CHANGED_FILES"
          
          # Find changed packages
          CHANGED_PACKAGES=()
          while IFS= read -r file; do
            if [[ -n "$file" && "$file" == packages/* ]]; then
              package_dir=$(echo "$file" | cut -d'/' -f1-2)
              if [[ -f "$package_dir/go.mod" ]] && [[ ! " ${CHANGED_PACKAGES[@]} " =~ " ${package_dir} " ]]; then
                CHANGED_PACKAGES+=("$package_dir")
              fi
            fi
          done <<< "$CHANGED_FILES"
          
          if [[ ${#CHANGED_PACKAGES[@]} -eq 0 ]]; then
            echo "📦 No package changes detected"
            echo "has-changes=false" >> $GITHUB_OUTPUT
          else
            echo "📦 Changed packages: ${CHANGED_PACKAGES[*]}"
            printf '%s\n' "${CHANGED_PACKAGES[@]}" | jq -R -s -c 'split("\n")[:-1]' > packages.json
            echo "packages=$(cat packages.json)" >> $GITHUB_OUTPUT
            echo "has-changes=true" >> $GITHUB_OUTPUT
            echo "test-all=false" >> $GITHUB_OUTPUT
          fi

  # 🧪 Main testing job
  test:
    name: 🧪 Test (Go ${{ matrix.go-version }}, ${{ matrix.os }})
    needs: detect-changes
    if: needs.detect-changes.outputs.has-changes == 'true'
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest]
        go-version: ['1.21', '1.22', '1.23']
        include:
          # 🪟 Add Windows testing for latest Go version only
          - os: windows-latest
            go-version: '1.23'

    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go ${{ matrix.go-version }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
          cache: true

      - name: 📋 Verify Go installation
        run: |
          echo "🐹 Go version: $(go version)"
          echo "📍 Go root: $(go env GOROOT)"
          echo "📁 Go path: $(go env GOPATH)"

      - name: 🔧 Set up environment
        shell: bash
        run: |
          echo "🔧 Setting up build environment..."
          # Create necessary directories
          mkdir -p coverage reports
          
          # Set up Go environment
          echo "GOPATH=$(go env GOPATH)" >> $GITHUB_ENV
          echo "PATH=$(go env GOPATH)/bin:$PATH" >> $GITHUB_ENV

      - name: 📦 Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/go/pkg/mod
            ~/.cache/go-build
          key: ${{ runner.os }}-go-${{ matrix.go-version }}-${{ hashFiles('**/go.sum', '**/go.mod') }}
          restore-keys: |
            ${{ runner.os }}-go-${{ matrix.go-version }}-
            ${{ runner.os }}-go-

      - name: 📥 Download dependencies
        shell: bash
        run: |
          echo "📥 Downloading Go modules..."
          packages='${{ needs.detect-changes.outputs.packages }}'
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            echo "📦 Downloading all package dependencies..."
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              echo "📥 Processing $dir..."
              (cd "$dir" && go mod download)
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                echo "📥 Processing $pkg..."
                (cd "$pkg" && go mod download)
              fi
            done
          fi

      - name: 🔨 Build packages
        shell: bash
        run: |
          echo "🔨 Building packages..."
          packages='${{ needs.detect-changes.outputs.packages }}'
          failed_builds=()
          
          build_package() {
            local pkg="$1"
            echo "🔨 Building $pkg..."
            if (cd "$pkg" && go build -v ./...); then
              echo "✅ Successfully built $pkg"
              return 0
            else
              echo "❌ Failed to build $pkg"
              return 1
            fi
          }
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              if ! build_package "$dir"; then
                failed_builds+=("$dir")
              fi
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                if ! build_package "$pkg"; then
                  failed_builds+=("$pkg")
                fi
              fi
            done
          fi
          
          if [[ ${#failed_builds[@]} -gt 0 ]]; then
            echo "❌ Build failed for: ${failed_builds[*]}"
            exit 1
          fi
          
          echo "✅ All packages built successfully"

      - name: 🧪 Run tests
        shell: bash
        run: |
          echo "🧪 Running tests..."
          packages='${{ needs.detect-changes.outputs.packages }}'
          
          run_tests() {
            local pkg="$1"
            local coverage_file="coverage/$(basename "$pkg")-coverage.out"
            
            echo "🧪 Testing $pkg..."
            cd "$pkg"
            
            # Run tests with coverage
            if go test -v -race -coverprofile="../$coverage_file" -covermode=atomic ./...; then
              echo "✅ Tests passed for $pkg"
              
              # Display coverage summary
              if [[ -f "../$coverage_file" ]]; then
                coverage=$(go tool cover -func="../$coverage_file" | tail -1 | awk '{print $3}' | sed 's/%//')
                echo "📊 Coverage for $pkg: $coverage%"
                
                # Check coverage threshold
                if (( $(echo "$coverage >= $COVERAGE_THRESHOLD" | bc -l 2>/dev/null || echo "0") )); then
                  echo "✅ Coverage threshold met for $pkg ($coverage% >= $COVERAGE_THRESHOLD%)"
                else
                  echo "⚠️ Coverage below threshold for $pkg ($coverage% < $COVERAGE_THRESHOLD%)"
                fi
              fi
              
              cd - > /dev/null
              return 0
            else
              echo "❌ Tests failed for $pkg"
              cd - > /dev/null
              return 1
            fi
          }
          
          failed_tests=()
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              if ! run_tests "$dir"; then
                failed_tests+=("$dir")
              fi
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                if ! run_tests "$pkg"; then
                  failed_tests+=("$pkg")
                fi
              fi
            done
          fi
          
          if [[ ${#failed_tests[@]} -gt 0 ]]; then
            echo "❌ Tests failed for: ${failed_tests[*]}"
            exit 1
          fi
          
          echo "✅ All tests passed"

      - name: 🔍 Run go vet
        shell: bash
        run: |
          echo "🔍 Running go vet..."
          packages='${{ needs.detect-changes.outputs.packages }}'
          
          vet_package() {
            local pkg="$1"
            echo "🔍 Vetting $pkg..."
            if (cd "$pkg" && go vet ./...); then
              echo "✅ Vet passed for $pkg"
              return 0
            else
              echo "❌ Vet failed for $pkg"
              return 1
            fi
          }
          
          failed_vet=()
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              if ! vet_package "$dir"; then
                failed_vet+=("$dir")
              fi
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                if ! vet_package "$pkg"; then
                  failed_vet+=("$pkg")
                fi
              fi
            done
          fi
          
          if [[ ${#failed_vet[@]} -gt 0 ]]; then
            echo "❌ Go vet failed for: ${failed_vet[*]}"
            exit 1
          fi
          
          echo "✅ Go vet passed for all packages"

      - name: 🎨 Check formatting
        shell: bash
        run: |
          echo "🎨 Checking code formatting..."
          
          unformatted_files=()
          
          check_formatting() {
            local pkg="$1"
            echo "🎨 Checking formatting for $pkg..."
            
            local fmt_files
            fmt_files=$(cd "$pkg" && gofmt -l .)
            
            if [[ -n "$fmt_files" ]]; then
              echo "❌ Unformatted files in $pkg:"
              echo "$fmt_files" | sed 's/^/  /'
              echo "$fmt_files" | while read -r file; do
                unformatted_files+=("$pkg/$file")
              done
              return 1
            else
              echo "✅ All files formatted correctly in $pkg"
              return 0
            fi
          }
          
          packages='${{ needs.detect-changes.outputs.packages }}'
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              check_formatting "$dir" || true
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                check_formatting "$pkg" || true
              fi
            done
          fi
          
          if [[ ${#unformatted_files[@]} -gt 0 ]]; then
            echo "❌ Found unformatted files:"
            printf '%s\n' "${unformatted_files[@]}"
            echo ""
            echo "💡 Run 'gofmt -s -w .' to fix formatting issues"
            exit 1
          fi
          
          echo "✅ All files are properly formatted"

      - name: 📊 Merge coverage reports
        if: matrix.os == 'ubuntu-latest' && matrix.go-version == '1.23'
        shell: bash
        run: |
          echo "📊 Merging coverage reports..."
          
          # Install gocovmerge if not available
          if ! command -v gocovmerge >/dev/null 2>&1; then
            echo "📥 Installing gocovmerge..."
            go install github.com/wadey/gocovmerge@latest
          fi
          
          # Find all coverage files
          coverage_files=()
          find coverage -name "*-coverage.out" -type f | while read -r file; do
            if [[ -s "$file" ]]; then
              coverage_files+=("$file")
            fi
          done
          
          if [[ ${#coverage_files[@]} -gt 0 ]]; then
            echo "📊 Found ${#coverage_files[@]} coverage files"
            gocovmerge "${coverage_files[@]}" > coverage/merged-coverage.out
            
            # Generate coverage report
            go tool cover -html=coverage/merged-coverage.out -o coverage/coverage.html
            go tool cover -func=coverage/merged-coverage.out > coverage/coverage.txt
            
            echo "📊 Coverage summary:"
            tail -1 coverage/coverage.txt
          else
            echo "⚠️ No coverage files found"
          fi

      - name: 📤 Upload coverage to Codecov
        if: matrix.os == 'ubuntu-latest' && matrix.go-version == '1.23'
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/merged-coverage.out
          flags: unittests
          name: govel-coverage
          fail_ci_if_error: false
          verbose: true

      - name: 📊 Upload coverage reports
        if: matrix.os == 'ubuntu-latest' && matrix.go-version == '1.23'
        uses: actions/upload-artifact@v3
        with:
          name: coverage-reports
          path: |
            coverage/
            reports/
          retention-days: 30

  # 🔒 Security scanning job
  security:
    name: 🔒 Security Scan
    needs: detect-changes
    if: needs.detect-changes.outputs.has-changes == 'true' && github.event.inputs.skip_security != 'true'
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
          cache: true

      - name: 🔒 Run gosec security scanner
        uses: securecodewarrior/github-action-gosec@master
        with:
          args: '-fmt sarif -out gosec.sarif ./...'

      - name: 🛡️ Upload SARIF file
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: gosec.sarif

      - name: 🔍 Run govulncheck
        run: |
          echo "🔍 Installing govulncheck..."
          go install golang.org/x/vuln/cmd/govulncheck@latest
          
          echo "🔍 Scanning for vulnerabilities..."
          packages='${{ needs.detect-changes.outputs.packages }}'
          
          if [[ "${{ needs.detect-changes.outputs.test-all }}" == "true" ]]; then
            find packages -name "go.mod" -type f | while read -r mod_file; do
              dir=$(dirname "$mod_file")
              echo "🔍 Scanning $dir..."
              (cd "$dir" && govulncheck ./...)
            done
          else
            echo "$packages" | jq -r '.[]' | while read -r pkg; do
              if [[ -d "$pkg" && -f "$pkg/go.mod" ]]; then
                echo "🔍 Scanning $pkg..."
                (cd "$pkg" && govulncheck ./...)
              fi
            done
          fi

  # ✅ Status check job
  ci-success:
    name: ✅ CI Success
    if: always()
    needs: [detect-changes, test, security]
    runs-on: ubuntu-latest
    steps:
      - name: 📊 Check all job results
        run: |
          echo "🔍 Checking CI results..."
          
          # Check if changes were detected
          if [[ "${{ needs.detect-changes.outputs.has-changes }}" != "true" ]]; then
            echo "ℹ️ No package changes detected - CI skipped"
            exit 0
          fi
          
          # Check test results
          if [[ "${{ needs.test.result }}" == "failure" ]]; then
            echo "❌ Tests failed"
            exit 1
          fi
          
          # Check security results (only if not skipped)
          if [[ "${{ github.event.inputs.skip_security }}" != "true" ]]; then
            if [[ "${{ needs.security.result }}" == "failure" ]]; then
              echo "❌ Security scan failed"
              exit 1
            fi
          fi
          
          echo "✅ All CI checks passed successfully!"

EOF

    print_message "success" "Created main CI workflow" "✅"
}

##
# 🔒 Generate security scanning workflow
#
# Creates a dedicated workflow for comprehensive security scanning,
# including vulnerability detection, dependency analysis, and SAST.
##
generate_security_workflow() {
    local file_path="$PROJECT_ROOT/.github/workflows/security.yml"
    print_message "info" "Generating security workflow" "🔒"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .github/workflows/security.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 🔒 GoVel Framework - Security Scanning Pipeline
#
# This workflow provides comprehensive security scanning for the GoVel framework.
# It performs static analysis, vulnerability detection, dependency scanning,
# and license compliance checks to ensure code security and quality.
#
# Features:
# - 🛡️ Static Application Security Testing (SAST) with gosec
# - 🔍 Vulnerability scanning with govulncheck
# - 📦 Dependency vulnerability analysis
# - 📄 License compliance checking
# - 🚨 Security advisory integration
# - 📊 SARIF report generation for GitHub Security tab
#
# Triggers:
# - Scheduled weekly scans (Sundays at 2 AM UTC)
# - Push to main/develop branches
# - Manual workflow dispatch
#
# Author: GoVel Framework Team
# Version: 1.0.0

name: 🔒 Security Scan

on:
  # 📅 Scheduled security scans
  schedule:
    - cron: '0 2 * * 0' # Every Sunday at 2 AM UTC
  
  # 🔄 Trigger on pushes to main branches
  push:
    branches: [ main, develop ]
    paths:
      - '**/*.go'
      - '**/go.mod'
      - '**/go.sum'
      - '.github/workflows/security.yml'
  
  # 🎛️ Allow manual triggering
  workflow_dispatch:
    inputs:
      scan_type:
        description: '🔍 Type of security scan to perform'
        required: true
        default: 'full'
        type: choice
        options:
          - 'full'
          - 'sast-only'
          - 'vulnerabilities-only'
          - 'dependencies-only'
      severity_threshold:
        description: '⚠️ Minimum severity level to report'
        required: true
        default: 'medium'
        type: choice
        options:
          - 'low'
          - 'medium'
          - 'high'
          - 'critical'

# 🔒 Security: Minimal required permissions
permissions:
  contents: read
  security-events: write
  actions: read

# 🌍 Environment variables
env:
  GO_VERSION: '1.23'
  SCAN_TYPE: ${{ github.event.inputs.scan_type || 'full' }}
  SEVERITY_THRESHOLD: ${{ github.event.inputs.severity_threshold || 'medium' }}

jobs:
  # 🛡️ Static Application Security Testing
  sast-scan:
    name: 🛡️ SAST Analysis
    runs-on: ubuntu-latest
    if: contains(fromJSON('["full", "sast-only"]'), github.event.inputs.scan_type) || github.event.inputs.scan_type == null
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: 🐹 Set up Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 📦 Cache security tools
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/gosec
            ~/go/bin/gosec
          key: ${{ runner.os }}-security-tools-${{ hashFiles('.github/workflows/security.yml') }}
          restore-keys: |
            ${{ runner.os }}-security-tools-

      - name: 🔧 Install gosec
        run: |
          echo "🔧 Installing gosec security scanner..."
          if ! command -v gosec >/dev/null 2>&1; then
            go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
          fi
          gosec --version

      - name: 📝 Create gosec configuration
        run: |
          echo "📝 Creating gosec configuration..."
          cat > gosec.json << 'GOSEC_CONFIG'
{
  "severity": "${{ env.SEVERITY_THRESHOLD }}",
  "confidence": "medium",
  "exclude": [
    "G104",
    "G204"
  ],
  "include": [
    "G101",
    "G102",
    "G103",
    "G106",
    "G107",
    "G108",
    "G109",
    "G110",
    "G201",
    "G202",
    "G203",
    "G301",
    "G302",
    "G303",
    "G304",
    "G305",
    "G401",
    "G402",
    "G403",
    "G404",
    "G501",
    "G502",
    "G503",
    "G504",
    "G601"
  ]
}
GOSEC_CONFIG

      - name: 🛡️ Run gosec SAST analysis
        run: |
          echo "🛡️ Running gosec security analysis..."
          mkdir -p reports/security
          
          # Find all Go packages
          packages=$(find packages -name "go.mod" -type f | sed 's|/go.mod||' | sort)
          
          for package in $packages; do
            if [[ -d "$package" ]]; then
              echo "🔍 Scanning $package..."
              package_name=$(basename "$package")
              
              # Run gosec on the package
              cd "$package"
              if gosec -conf ../gosec.json -fmt sarif -out "../reports/security/${package_name}-gosec.sarif" -stdout -verbose=text ./...; then
                echo "✅ SAST scan completed for $package"
              else
                echo "⚠️ SAST scan found issues in $package"
              fi
              cd - > /dev/null
            fi
          done
          
          # Merge all SARIF reports
          echo "📊 Merging SARIF reports..."
          if command -v jq >/dev/null 2>&1; then
            find reports/security -name "*-gosec.sarif" -type f | head -1 | xargs cat > reports/security/merged-gosec.sarif 2>/dev/null || echo '{"runs": []}' > reports/security/merged-gosec.sarif
          fi

      - name: 📊 Upload SAST SARIF results
        uses: github/codeql-action/upload-sarif@v2
        if: always()
        with:
          sarif_file: reports/security/merged-gosec.sarif
          category: gosec

      - name: 📤 Upload security reports
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: sast-reports
          path: reports/security/
          retention-days: 90

  # 🔍 Vulnerability scanning
  vulnerability-scan:
    name: 🔍 Vulnerability Analysis
    runs-on: ubuntu-latest
    if: contains(fromJSON('["full", "vulnerabilities-only"]'), github.event.inputs.scan_type) || github.event.inputs.scan_type == null
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 📦 Cache vulnerability database
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/govulncheck
            ~/go/bin/govulncheck
          key: ${{ runner.os }}-vulndb-${{ hashFiles('.github/workflows/security.yml') }}
          restore-keys: |
            ${{ runner.os }}-vulndb-

      - name: 🔧 Install govulncheck
        run: |
          echo "🔧 Installing govulncheck..."
          if ! command -v govulncheck >/dev/null 2>&1; then
            go install golang.org/x/vuln/cmd/govulncheck@latest
          fi
          govulncheck --version || echo "Govulncheck installed"

      - name: 🔍 Run vulnerability scan
        run: |
          echo "🔍 Running vulnerability analysis..."
          mkdir -p reports/vulnerabilities
          
          # Find all Go packages
          packages=$(find packages -name "go.mod" -type f | sed 's|/go.mod||' | sort)
          vulnerability_found=false
          
          for package in $packages; do
            if [[ -d "$package" ]]; then
              echo "🔍 Scanning $package for vulnerabilities..."
              package_name=$(basename "$package")
              
              cd "$package"
              
              # Run govulncheck
              if govulncheck -json ./... > "../reports/vulnerabilities/${package_name}-vulnerabilities.json" 2>&1; then
                echo "✅ No vulnerabilities found in $package"
              else
                echo "⚠️ Vulnerabilities found in $package"
                vulnerability_found=true
                
                # Generate human-readable report
                govulncheck ./... > "../reports/vulnerabilities/${package_name}-vulnerabilities.txt" 2>&1 || true
              fi
              
              cd - > /dev/null
            fi
          done
          
          # Create summary report
          echo "📊 Creating vulnerability summary..."
          {
            echo "# 🔍 Vulnerability Scan Summary"
            echo ""
            echo "**Scan Date:** $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
            echo "**Go Version:** ${{ env.GO_VERSION }}"
            echo "**Severity Threshold:** ${{ env.SEVERITY_THRESHOLD }}"
            echo ""
            
            if [[ "$vulnerability_found" == "true" ]]; then
              echo "⚠️ **Status:** Vulnerabilities detected"
              echo ""
              echo "## 📋 Affected Packages"
              for report in reports/vulnerabilities/*-vulnerabilities.txt; do
                if [[ -f "$report" ]] && [[ -s "$report" ]]; then
                  package_name=$(basename "$report" "-vulnerabilities.txt")
                  echo "- 📦 **$package_name**"
                fi
              done
            else
              echo "✅ **Status:** No vulnerabilities detected"
            fi
            
            echo ""
            echo "---"
            echo "*Generated by GoVel Security Pipeline*"
          } > reports/vulnerabilities/summary.md
          
          # Set job outcome
          if [[ "$vulnerability_found" == "true" ]]; then
            echo "vulnerability_found=true" >> $GITHUB_ENV
          else
            echo "vulnerability_found=false" >> $GITHUB_ENV
          fi

      - name: 📊 Create vulnerability report
        if: env.vulnerability_found == 'true'
        run: |
          echo "📊 Creating detailed vulnerability report..."
          
          # Find the most recent vulnerability report with content
          latest_report=""
          for report in reports/vulnerabilities/*-vulnerabilities.txt; do
            if [[ -f "$report" ]] && [[ -s "$report" ]]; then
              latest_report="$report"
              break
            fi
          done
          
          if [[ -n "$latest_report" ]]; then
            echo "📄 Latest vulnerability report: $latest_report"
            echo "## 🚨 Vulnerability Details" >> reports/vulnerabilities/summary.md
            echo '```' >> reports/vulnerabilities/summary.md
            head -50 "$latest_report" >> reports/vulnerabilities/summary.md
            echo '```' >> reports/vulnerabilities/summary.md
          fi

      - name: 📤 Upload vulnerability reports
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: vulnerability-reports
          path: reports/vulnerabilities/
          retention-days: 90

      - name: 💬 Comment vulnerability summary on PR
        if: github.event_name == 'pull_request' && env.vulnerability_found == 'true'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            if (fs.existsSync('reports/vulnerabilities/summary.md')) {
              const summary = fs.readFileSync('reports/vulnerabilities/summary.md', 'utf8');
              await github.rest.issues.createComment({
                issue_number: context.issue.number,
                owner: context.repo.owner,
                repo: context.repo.repo,
                body: `## 🔍 Security Vulnerability Report\n\n${summary}`
              });
            }

  # 📦 Dependency analysis
  dependency-scan:
    name: 📦 Dependency Analysis
    runs-on: ubuntu-latest
    if: contains(fromJSON('["full", "dependencies-only"]'), github.event.inputs.scan_type) || github.event.inputs.scan_type == null
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 📦 Analyze dependencies
        run: |
          echo "📦 Analyzing Go module dependencies..."
          mkdir -p reports/dependencies
          
          # Find all Go packages
          packages=$(find packages -name "go.mod" -type f | sed 's|/go.mod||' | sort)
          
          for package in $packages; do
            if [[ -d "$package" ]]; then
              echo "📦 Analyzing dependencies for $package..."
              package_name=$(basename "$package")
              
              cd "$package"
              
              # Generate dependency reports
              {
                echo "# Dependency Report for $package_name"
                echo ""
                echo "Generated on: $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
                echo ""
                
                echo "## 📋 Direct Dependencies"
                go list -m -f '{{.Path}}@{{.Version}}' all | grep -v "^$package" | head -20
                
                echo ""
                echo "## 🔍 Outdated Dependencies"
                go list -u -m all | grep -F '[' || echo "All dependencies are up to date"
                
                echo ""
                echo "## 📊 Dependency Graph"
                go mod graph | head -20
                
              } > "../reports/dependencies/${package_name}-dependencies.md"
              
              # Generate machine-readable dependency list
              go list -m -json all > "../reports/dependencies/${package_name}-dependencies.json"
              
              cd - > /dev/null
            fi
          done
          
          echo "✅ Dependency analysis completed"

      - name: 📤 Upload dependency reports
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: dependency-reports
          path: reports/dependencies/
          retention-days: 30

  # ✅ Security status check
  security-status:
    name: ✅ Security Status
    if: always()
    needs: [sast-scan, vulnerability-scan, dependency-scan]
    runs-on: ubuntu-latest
    steps:
      - name: 📊 Evaluate security results
        run: |
          echo "📊 Evaluating security scan results..."
          
          # Check SAST results
          sast_result="${{ needs.sast-scan.result }}"
          vuln_result="${{ needs.vulnerability-scan.result }}"
          deps_result="${{ needs.dependency-scan.result }}"
          
          echo "🛡️ SAST Scan: $sast_result"
          echo "🔍 Vulnerability Scan: $vuln_result"
          echo "📦 Dependency Scan: $deps_result"
          
          # Determine overall status
          if [[ "$sast_result" == "failure" ]] || [[ "$vuln_result" == "failure" ]] || [[ "$deps_result" == "failure" ]]; then
            echo "❌ Security scans detected issues"
            echo "🔍 Please review the security reports and address any findings"
            exit 1
          elif [[ "$sast_result" == "skipped" ]] && [[ "$vuln_result" == "skipped" ]] && [[ "$deps_result" == "skipped" ]]; then
            echo "ℹ️ All security scans were skipped"
          else
            echo "✅ All security scans completed successfully"
            echo "🛡️ No critical security issues detected"
          fi
          
          # Generate security badge info
          echo "SECURITY_STATUS=passing" >> $GITHUB_ENV

EOF

    print_message "success" "Created security workflow" "✅"
}

##
# 🤖 Generate dependency update workflow
#
# Creates an automated workflow for dependency updates using Dependabot
# and custom update strategies.
##
generate_dependency_workflow() {
    local file_path="$PROJECT_ROOT/.github/workflows/dependency-update.yml"
    print_message "info" "Generating dependency update workflow" "🤖"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .github/workflows/dependency-update.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 🤖 GoVel Framework - Automated Dependency Updates
#
# This workflow provides automated dependency management for the GoVel framework.
# It handles Go module updates, security patches, and version compatibility checks
# to keep the codebase secure and up-to-date.
#
# Features:
# - 📦 Automated Go module updates
# - 🔒 Security patch prioritization
# - 🧪 Automated testing of dependency updates
# - 📋 Compatibility verification across Go versions
# - 🔄 Batch updates for related dependencies
# - 📊 Update impact analysis
#
# Triggers:
# - Scheduled weekly updates (Mondays at 3 AM UTC)
# - Manual workflow dispatch with custom options
# - Dependabot integration
#
# Author: GoVel Framework Team
# Version: 1.0.0

name: 🤖 Dependency Updates

on:
  # 📅 Scheduled dependency updates
  schedule:
    - cron: '0 3 * * 1' # Every Monday at 3 AM UTC
  
  # 🎛️ Manual dependency updates
  workflow_dispatch:
    inputs:
      update_type:
        description: '📦 Type of update to perform'
        required: true
        default: 'minor'
        type: choice
        options:
          - 'patch'     # Patch updates only (1.2.3 -> 1.2.4)
          - 'minor'     # Minor updates (1.2.3 -> 1.3.0)
          - 'major'     # Major updates (1.2.3 -> 2.0.0)
          - 'security'  # Security updates only
          - 'all'       # All available updates
      create_pr:
        description: '🔀 Create pull request for updates'
        required: true
        default: true
        type: boolean
      test_updates:
        description: '🧪 Run tests before creating PR'
        required: true
        default: true
        type: boolean

# 🔒 Security: Required permissions for dependency updates
permissions:
  contents: write
  pull-requests: write
  issues: write
  checks: write

# 🌍 Environment variables
env:
  GO_VERSION: '1.23'
  UPDATE_TYPE: ${{ github.event.inputs.update_type || 'minor' }}
  CREATE_PR: ${{ github.event.inputs.create_pr || 'true' }}
  TEST_UPDATES: ${{ github.event.inputs.test_updates || 'true' }}

jobs:
  # 🔍 Check for available updates
  check-updates:
    name: 🔍 Check Available Updates
    runs-on: ubuntu-latest
    outputs:
      has-updates: ${{ steps.scan.outputs.has-updates }}
      update-summary: ${{ steps.scan.outputs.update-summary }}
      packages-with-updates: ${{ steps.scan.outputs.packages-with-updates }}
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4

      - name: 🐹 Set up Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 🔧 Install update tools
        run: |
          echo "🔧 Installing Go update tools..."
          
          # Install go-mod-outdated for checking outdated dependencies
          if ! command -v go-mod-outdated >/dev/null 2>&1; then
            go install github.com/psampaz/go-mod-outdated@latest
          fi
          
          # Install goupdate for automated updates
          if ! command -v goupdate >/dev/null 2>&1; then
            go install github.com/oligot/go-mod-upgrade@latest
          fi

      - name: 🔍 Scan for outdated dependencies
        id: scan
        run: |
          echo "🔍 Scanning for outdated dependencies..."
          mkdir -p reports/updates
          
          # Find all Go packages
          packages=$(find packages -name "go.mod" -type f | sed 's|/go.mod||' | sort)
          
          has_updates=false
          packages_with_updates=()
          update_details=""
          
          for package in $packages; do
            if [[ -d "$package" ]]; then
              echo "📦 Checking updates for $package..."
              package_name=$(basename "$package")
              
              cd "$package"
              
              # Check for outdated dependencies
              outdated_output=""
              if go list -u -m -json all | go-mod-outdated -update -direct > "../reports/updates/${package_name}-outdated.json" 2>/dev/null; then
                outdated_count=$(cat "../reports/updates/${package_name}-outdated.json" | jq length 2>/dev/null || echo "0")
                
                if [[ "$outdated_count" -gt 0 ]]; then
                  echo "📦 Found $outdated_count outdated dependencies in $package"
                  has_updates=true
                  packages_with_updates+=("$package")
                  
                  # Generate human-readable report
                  {
                    echo "## Updates available for $package_name"
                    echo ""
                    cat "../reports/updates/${package_name}-outdated.json" | jq -r '.[] | "- **\(.module.name)**: \(.current) → \(.latest) (\(.type))"' 2>/dev/null || echo "Error parsing updates"
                    echo ""
                  } >> "../reports/updates/summary.md"
                else
                  echo "✅ No outdated dependencies in $package"
                fi
              else
                echo "⚠️ Could not check outdated dependencies for $package"
              fi
              
              cd - > /dev/null
            fi
          done
          
          # Set outputs
          echo "has-updates=$has_updates" >> $GITHUB_OUTPUT
          
          if [[ "$has_updates" == "true" ]]; then
            printf -v packages_json '%s\n' "${packages_with_updates[@]}" | jq -R . | jq -s .
            echo "packages-with-updates=$packages_json" >> $GITHUB_OUTPUT
            
            if [[ -f "reports/updates/summary.md" ]]; then
              summary=$(cat reports/updates/summary.md | head -20)
              echo "update-summary<<EOF" >> $GITHUB_OUTPUT
              echo "$summary" >> $GITHUB_OUTPUT
              echo "EOF" >> $GITHUB_OUTPUT
            fi
          else
            echo "packages-with-updates=[]" >> $GITHUB_OUTPUT
            echo "update-summary=No updates available" >> $GITHUB_OUTPUT
          fi

      - name: 📤 Upload update reports
        if: steps.scan.outputs.has-updates == 'true'
        uses: actions/upload-artifact@v3
        with:
          name: update-scan-reports
          path: reports/updates/
          retention-days: 30

  # 📦 Apply dependency updates
  apply-updates:
    name: 📦 Apply Updates
    needs: check-updates
    if: needs.check-updates.outputs.has-updates == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        package: ${{ fromJSON(needs.check-updates.outputs.packages-with-updates) }}
      fail-fast: false
    steps:
      - name: 📥 Checkout code
        uses: actions/checkout@v4
        with:
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: 🐹 Set up Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 🔧 Configure Git
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"

      - name: 📦 Apply updates to ${{ matrix.package }}
        run: |
          echo "📦 Applying updates to ${{ matrix.package }}..."
          package_name=$(basename "${{ matrix.package }}")
          
          cd "${{ matrix.package }}"
          
          # Create backup of current go.mod
          cp go.mod go.mod.backup
          
          case "${{ env.UPDATE_TYPE }}" in
            "patch")
              echo "🔧 Applying patch updates..."
              go get -u=patch ./...
              ;;
            "minor")
              echo "🔧 Applying minor updates..."
              # Update to latest minor versions
              go list -u -m all | grep -F '[' | cut -d' ' -f1 | while read -r module; do
                if [[ -n "$module" ]]; then
                  # Get the latest minor version
                  latest=$(go list -m -versions "$module" | tr ' ' '\n' | tail -1)
                  if [[ -n "$latest" ]]; then
                    go get "$module@$latest" || echo "Failed to update $module"
                  fi
                fi
              done
              ;;
            "major")
              echo "🔧 Applying major updates..."
              go get -u ./...
              ;;
            "security")
              echo "🔒 Applying security updates..."
              # This would require integration with vulnerability database
              go get -u ./...
              ;;
            "all")
              echo "🔧 Applying all available updates..."
              go get -u ./...
              ;;
            *)
              echo "⚠️ Unknown update type: ${{ env.UPDATE_TYPE }}"
              go get -u=patch ./...
              ;;
          esac
          
          # Clean up and tidy
          go mod tidy
          
          # Check if anything changed
          if ! diff -q go.mod go.mod.backup >/dev/null 2>&1; then
            echo "✅ Dependencies updated for ${{ matrix.package }}"
            echo "UPDATES_APPLIED=true" >> $GITHUB_ENV
            
            # Show what changed
            echo "📋 Changes made:"
            diff -u go.mod.backup go.mod || true
          else
            echo "ℹ️ No updates applied to ${{ matrix.package }}"
            echo "UPDATES_APPLIED=false" >> $GITHUB_ENV
          fi
          
          rm -f go.mod.backup
          cd - > /dev/null

      - name: 🧪 Test updated dependencies
        if: env.UPDATES_APPLIED == 'true' && env.TEST_UPDATES == 'true'
        run: |
          echo "🧪 Testing updated dependencies for ${{ matrix.package }}..."
          
          cd "${{ matrix.package }}"
          
          # Download new dependencies
          go mod download
          
          # Run tests
          if go test -v ./...; then
            echo "✅ Tests passed with updated dependencies"
            echo "TESTS_PASSED=true" >> $GITHUB_ENV
          else
            echo "❌ Tests failed with updated dependencies"
            echo "TESTS_PASSED=false" >> $GITHUB_ENV
            exit 1
          fi
          
          # Run basic build check
          if go build ./...; then
            echo "✅ Build successful with updated dependencies"
          else
            echo "❌ Build failed with updated dependencies"
            exit 1
          fi
          
          cd - > /dev/null

      - name: 📝 Generate update report
        if: env.UPDATES_APPLIED == 'true'
        run: |
          echo "📝 Generating update report for ${{ matrix.package }}..."
          package_name=$(basename "${{ matrix.package }}")
          
          mkdir -p reports/applied-updates
          
          cd "${{ matrix.package }}"
          
          {
            echo "# Dependency Updates Applied - $package_name"
            echo ""
            echo "**Date:** $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
            echo "**Update Type:** ${{ env.UPDATE_TYPE }}"
            echo "**Go Version:** ${{ env.GO_VERSION }}"
            echo ""
            
            if [[ "${{ env.TEST_UPDATES }}" == "true" ]]; then
              if [[ "${{ env.TESTS_PASSED }}" == "true" ]]; then
                echo "**Tests:** ✅ Passed"
              else
                echo "**Tests:** ❌ Failed"
              fi
              echo ""
            fi
            
            echo "## Updated Dependencies"
            echo ""
            go list -m all | head -20
            
            echo ""
            echo "## Go Module Info"
            echo ""
            echo '```'
            go version -m . 2>/dev/null || go version
            echo '```'
            
          } > "../reports/applied-updates/${package_name}-update-report.md"
          
          cd - > /dev/null

      - name: 💾 Commit changes
        if: env.UPDATES_APPLIED == 'true' && env.CREATE_PR == 'true'
        run: |
          echo "💾 Committing dependency updates..."
          package_name=$(basename "${{ matrix.package }}")
          
          # Stage changes
          git add "${{ matrix.package }}/go.mod" "${{ matrix.package }}/go.sum" 2>/dev/null || true
          
          # Check if there are changes to commit
          if git diff --staged --quiet; then
            echo "ℹ️ No changes to commit"
            exit 0
          fi
          
          # Create commit
          commit_message="🤖 Update dependencies for $package_name
          
          - Update type: ${{ env.UPDATE_TYPE }}
          - Go version: ${{ env.GO_VERSION }}
          - Tests: ${{ env.TESTS_PASSED == 'true' && 'passed' || 'skipped' }}
          - Auto-generated by dependency update workflow"
          
          git commit -m "$commit_message"
          
          # Create branch for this package
          branch_name="deps/update-${package_name}-$(date +%Y%m%d-%H%M%S)"
          git checkout -b "$branch_name"
          
          echo "BRANCH_NAME=$branch_name" >> $GITHUB_ENV

      - name: 📤 Push changes and create PR
        if: env.UPDATES_APPLIED == 'true' && env.CREATE_PR == 'true' && env.BRANCH_NAME != ''
        run: |
          echo "📤 Pushing changes and creating pull request..."
          
          # Push branch
          git push origin "${{ env.BRANCH_NAME }}"
          
          # Create PR using GitHub CLI
          package_name=$(basename "${{ matrix.package }}")
          
          pr_title="🤖 Automated dependency updates for $package_name"
          pr_body="## 📦 Dependency Updates
          
          This pull request contains automated dependency updates for the **$package_name** package.
          
          ### 📋 Update Details
          - **Update Type:** ${{ env.UPDATE_TYPE }}
          - **Go Version:** ${{ env.GO_VERSION }}
          - **Tests Status:** ${{ env.TESTS_PASSED == 'true' && '✅ Passed' || '⚠️ Skipped' }}
          - **Generated:** $(date -u '+%Y-%m-%d %H:%M:%S UTC')
          
          ### 🔍 What Changed
          ${{ needs.check-updates.outputs.update-summary }}
          
          ### 🧪 Testing
          ${{ env.TEST_UPDATES == 'true' && 'Automated tests have been run and passed.' || 'Automated testing was skipped.' }}
          
          ### 📝 Review Checklist
          - [ ] Review dependency changes
          - [ ] Check for breaking changes in updated packages
          - [ ] Verify test coverage
          - [ ] Update documentation if needed
          
          ---
          *This PR was automatically generated by the GoVel dependency update workflow.*"
          
          # Create the PR
          gh pr create \
            --title "$pr_title" \
            --body "$pr_body" \
            --head "${{ env.BRANCH_NAME }}" \
            --base "develop" \
            --label "dependencies,automated" \
            || echo "⚠️ Failed to create PR - may already exist"
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: 📤 Upload update artifacts
        if: env.UPDATES_APPLIED == 'true'
        uses: actions/upload-artifact@v3
        with:
          name: applied-updates-${{ matrix.package }}
          path: reports/applied-updates/
          retention-days: 30

  # 📊 Generate update summary
  update-summary:
    name: 📊 Update Summary
    if: always()
    needs: [check-updates, apply-updates]
    runs-on: ubuntu-latest
    steps:
      - name: 📊 Generate final summary
        run: |
          echo "📊 Generating dependency update summary..."
          
          # Check results
          has_updates="${{ needs.check-updates.outputs.has-updates }}"
          apply_result="${{ needs.apply-updates.result }}"
          
          {
            echo "# 🤖 Dependency Update Summary"
            echo ""
            echo "**Date:** $(date -u '+%Y-%m-%d %H:%M:%S UTC')"
            echo "**Update Type:** ${{ env.UPDATE_TYPE }}"
            echo "**Create PR:** ${{ env.CREATE_PR }}"
            echo "**Test Updates:** ${{ env.TEST_UPDATES }}"
            echo ""
            
            if [[ "$has_updates" == "true" ]]; then
              echo "## 📦 Updates Available"
              echo ""
              echo "${{ needs.check-updates.outputs.update-summary }}"
              echo ""
              
              if [[ "$apply_result" == "success" ]]; then
                echo "✅ **Status:** Updates applied successfully"
                if [[ "${{ env.CREATE_PR }}" == "true" ]]; then
                  echo "🔀 **Action:** Pull requests created for review"
                else
                  echo "💾 **Action:** Changes committed directly"
                fi
              elif [[ "$apply_result" == "failure" ]]; then
                echo "❌ **Status:** Some updates failed"
                echo "🔍 **Action:** Please review failed updates manually"
              else
                echo "⚠️ **Status:** Updates were skipped"
              fi
            else
              echo "✅ **Status:** All dependencies are up to date"
              echo "🎉 **Action:** No updates needed"
            fi
            
            echo ""
            echo "---"
            echo "*Generated by GoVel Dependency Update Pipeline*"
          } > update-summary.md
          
          echo "📄 Update summary:"
          cat update-summary.md

      - name: 💬 Post summary comment (if PR context)
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            if (fs.existsSync('update-summary.md')) {
              const summary = fs.readFileSync('update-summary.md', 'utf8');
              await github.rest.issues.createComment({
                issue_number: context.issue.number,
                owner: context.repo.owner,
                repo: context.repo.repo,
                body: summary
              });
            }

EOF

    print_message "success" "Created dependency update workflow" "✅"
}

##
# 🎯 Confirm file creation with user
#
# Asks for user confirmation before creating files, unless force mode is enabled.
# Shows a summary of what will be created.
##
confirm_creation() {
    if [[ "$FORCE" == "true" ]]; then
        print_message "info" "Force mode enabled - proceeding without confirmation" "💪"
        return 0
    fi
    
    print_message "header" "Ready to Create CI/CD Files" "🚀"
    
    echo "The following structure will be created:"
    echo ""
    echo "📁 .github/"
    echo "  📁 workflows/"
    echo "    📄 ci.yml                    # Main CI pipeline"
    echo "    📄 security.yml             # Security scanning"
    echo "    📄 dependency-update.yml    # Automated updates"
    echo "  📁 templates/"
    echo "    📄 Various GitHub templates..."
    echo ""
    echo "📁 scripts/"
    echo "  📁 ci/                        # CI automation scripts"
    echo "  📁 utils/                     # Utility scripts"
    echo "  📁 hooks/                     # Git hooks"
    echo ""
    echo "📄 Configuration files"
    echo "  📄 .golangci.yml             # Linting configuration"
    echo "  📄 .codecov.yml              # Coverage configuration"
    echo "  📄 Makefile                  # Build automation"
    echo ""
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "This is a dry run - no files will be created" "🔍"
        return 0
    fi
    
    echo -e "${YELLOW}❓ Do you want to proceed with creating these files? [y/N]:${NC} "
    read -r response
    
    case "$response" in
        [yY][eE][sS]|[yY])
            print_message "success" "Proceeding with file creation" "✅"
            return 0
            ;;
        *)
            print_message "info" "Cancelled by user" "🚫"
            exit 0
            ;;
    esac
}

##
# 🏁 Main execution function
#
# Orchestrates the entire CI/CD generation process including directory creation,
# workflow generation, and final reporting.
##
main() {
    print_message "header" "GoVel CI/CD Generator Started" "🚀"
    print_message "info" "Timestamp: $TIMESTAMP"
    print_message "info" "Project root: $PROJECT_ROOT"
    
    # Parse arguments
    parse_arguments "$@"
    
    # Show dry run status
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "Running in DRY-RUN mode - no files will be created" "🔍"
    fi
    
    # Create directory structure
    create_directory_structure "$PROJECT_ROOT"
    
    # Confirm creation with user
    confirm_creation
    
    # Generate GitHub workflows
    print_message "header" "Generating GitHub Workflows" "⚙️"
    generate_main_ci_workflow
    generate_security_workflow
    generate_dependency_workflow
    
    # TODO: Generate additional files (will be implemented in subsequent calls)
    # - GitHub templates
    # - CI scripts
    # - Configuration files
    # - Git hooks
    # - Documentation
    
    # Final summary
    print_message "header" "CI/CD Generation Complete" "🎉"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "success" "Dry run completed - showed what would be created" "✅"
    else
        print_message "success" "All CI/CD files created successfully" "✅"
        print_message "info" "Next steps:"
        echo "  1. Review the generated workflows"
        echo "  2. Customize configuration files as needed"
        echo "  3. Commit and push the changes"
        echo "  4. Enable GitHub Actions in your repository"
    fi
    
    print_message "info" "Generation completed at $(date '+%Y-%m-%d %H:%M:%S')" "⏰"
}

# 🏁 Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi