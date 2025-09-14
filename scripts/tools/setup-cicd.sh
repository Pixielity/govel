#!/bin/bash

##
# 🚀 GoVel Framework - Complete CI/CD Setup Script
#
# This master script orchestrates the complete setup of GitHub Actions CI/CD
# pipeline for the GoVel framework. It runs both the main workflow generator
# and configuration file generator in the correct sequence.
#
# Features:
# - 🎯 Complete CI/CD pipeline setup
# - 📄 GitHub Actions workflows 
# - 🔧 Development configurations
# - 🤖 Automated dependency management
# - 🔒 Security scanning setup
# - 🛠️ Build automation
# - 📊 Code quality gates
#
# Usage:
#   ./setup-cicd.sh [options]
#
# Options:
#   --dry-run, -d     Show what would be created without creating files
#   --help, -h        Show this help message  
#   --verbose, -v     Enable verbose output
#   --force, -f       Overwrite existing files without confirmation
#   --workflows-only  Only generate GitHub workflows (skip configs)
#   --configs-only    Only generate configuration files (skip workflows)
#
# Examples:
#   ./setup-cicd.sh --dry-run       # Preview complete setup
#   ./setup-cicd.sh                 # Interactive setup with confirmations
#   ./setup-cicd.sh --force         # Overwrite existing files
#   ./setup-cicd.sh --workflows-only # Only GitHub workflows
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
WORKFLOWS_ONLY=false
CONFIGS_ONLY=false

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
${BOLD}${CYAN}🚀 GoVel Complete CI/CD Setup${NC}

${BOLD}DESCRIPTION:${NC}
    Complete CI/CD pipeline setup for the GoVel framework.
    Generates GitHub Actions workflows, configuration files, and automation scripts.

${BOLD}USAGE:${NC}
    $SCRIPT_NAME [OPTIONS]

${BOLD}OPTIONS:${NC}
    -d, --dry-run         🔍 Show what would be created without creating files
    -v, --verbose         📝 Enable detailed output and debug information
    -f, --force           💪 Overwrite existing files without confirmation
    -h, --help            ❓ Show this help message and exit
    --workflows-only      ⚙️ Only generate GitHub Actions workflows
    --configs-only        🔧 Only generate configuration files

${BOLD}EXAMPLES:${NC}
    $SCRIPT_NAME --dry-run           # Preview complete setup
    $SCRIPT_NAME                     # Interactive setup with confirmations
    $SCRIPT_NAME --force --verbose   # Force setup with detailed output
    $SCRIPT_NAME --workflows-only    # Only GitHub workflows

${BOLD}WHAT GETS CREATED:${NC}

📁 GitHub Actions Workflows:
    🎯 ci.yml                    # Main CI pipeline
    🔒 security.yml             # Security scanning
    🤖 dependency-update.yml    # Automated dependency updates

📁 Configuration Files:
    🔧 .golangci.yml             # Comprehensive linting
    📊 .codecov.yml              # Code coverage reporting
    🤖 .github/dependabot.yml   # Dependency automation
    🛠️ Makefile                 # Build automation

${BOLD}FEATURES:${NC}
    ✅ Multi-version Go testing (1.21, 1.22, 1.23)
    ✅ Cross-platform testing (Linux, macOS, Windows)
    ✅ Smart package change detection
    ✅ Comprehensive security scanning
    ✅ Automated dependency updates
    ✅ Code coverage with quality gates
    ✅ 50+ linting rules
    ✅ Build automation with Makefile
    ✅ Professional development workflow

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
            --workflows-only)
                WORKFLOWS_ONLY=true
                print_message "info" "Workflows-only mode - skipping configuration files" "⚙️"
                shift
                ;;
            --configs-only)
                CONFIGS_ONLY=true
                print_message "info" "Configs-only mode - skipping workflows" "🔧"
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
    
    # Validate conflicting options
    if [[ "$WORKFLOWS_ONLY" == "true" && "$CONFIGS_ONLY" == "true" ]]; then
        print_message "error" "Cannot use --workflows-only and --configs-only together"
        exit 1
    fi
}

##
# 🔍 Check prerequisites
##
check_prerequisites() {
    print_message "header" "Checking Prerequisites" "🔍"
    
    # Check if we're in a Git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        print_message "warning" "Not in a Git repository - GitHub workflows may not work properly" "⚠️"
    else
        print_message "success" "Git repository detected" "✅"
    fi
    
    # Check Go installation
    if command -v go >/dev/null 2>&1; then
        local go_version=$(go version | awk '{print $3}' | sed 's/go//')
        print_message "success" "Go $go_version detected" "✅"
    else
        print_message "warning" "Go not installed - required for development workflow" "⚠️"
    fi
    
    # Check if required scripts exist
    local workflow_script="$PROJECT_ROOT/generate-cicd.sh"
    local config_script="$PROJECT_ROOT/generate-cicd-configs.sh"
    
    if [[ ! -f "$workflow_script" ]]; then
        print_message "error" "Workflow generator script not found: $workflow_script" "❌"
        exit 1
    fi
    
    if [[ ! -f "$config_script" ]]; then
        print_message "error" "Configuration generator script not found: $config_script" "❌"
        exit 1
    fi
    
    # Check script permissions
    if [[ ! -x "$workflow_script" ]]; then
        print_message "info" "Making workflow script executable" "🔧"
        chmod +x "$workflow_script"
    fi
    
    if [[ ! -x "$config_script" ]]; then
        print_message "info" "Making configuration script executable" "🔧"
        chmod +x "$config_script"
    fi
    
    print_message "success" "All prerequisites checked" "✅"
}

##
# ⚙️ Generate GitHub Actions workflows
##
generate_workflows() {
    if [[ "$CONFIGS_ONLY" == "true" ]]; then
        print_message "info" "Skipping workflows (configs-only mode)" "⏩"
        return 0
    fi
    
    print_message "header" "Generating GitHub Actions Workflows" "⚙️"
    
    local args=()
    
    if [[ "$DRY_RUN" == "true" ]]; then
        args+=("--dry-run")
    fi
    
    if [[ "$VERBOSE" == "true" ]]; then
        args+=("--verbose")
    fi
    
    if [[ "$FORCE" == "true" ]]; then
        args+=("--force")
    fi
    
    # Run the workflow generator
    print_message "info" "Executing workflow generator..." "🔄"
    if [[ ${#args[@]} -eq 0 ]]; then
        if "$PROJECT_ROOT/generate-cicd.sh"; then
            print_message "success" "GitHub Actions workflows generated successfully" "✅"
        else
            print_message "error" "Failed to generate GitHub Actions workflows" "❌"
            exit 1
        fi
    else
        if "$PROJECT_ROOT/generate-cicd.sh" "${args[@]}"; then
            print_message "success" "GitHub Actions workflows generated successfully" "✅"
        else
            print_message "error" "Failed to generate GitHub Actions workflows" "❌"
            exit 1
        fi
    fi
}

##
# 🔧 Generate configuration files
##
generate_configs() {
    if [[ "$WORKFLOWS_ONLY" == "true" ]]; then
        print_message "info" "Skipping configuration files (workflows-only mode)" "⏩"
        return 0
    fi
    
    print_message "header" "Generating Configuration Files" "🔧"
    
    local args=()
    
    if [[ "$DRY_RUN" == "true" ]]; then
        args+=("--dry-run")
    fi
    
    if [[ "$VERBOSE" == "true" ]]; then
        args+=("--verbose")
    fi
    
    if [[ "$FORCE" == "true" ]]; then
        args+=("--force")
    fi
    
    # Run the configuration generator
    print_message "info" "Executing configuration generator..." "🔄"
    if [[ ${#args[@]} -eq 0 ]]; then
        if "$PROJECT_ROOT/generate-cicd-configs.sh"; then
            print_message "success" "Configuration files generated successfully" "✅"
        else
            print_message "error" "Failed to generate configuration files" "❌"
            exit 1
        fi
    else
        if "$PROJECT_ROOT/generate-cicd-configs.sh" "${args[@]}"; then
            print_message "success" "Configuration files generated successfully" "✅"
        else
            print_message "error" "Failed to generate configuration files" "❌"
            exit 1
        fi
    fi
}

##
# 📋 Display final summary
##
show_final_summary() {
    print_message "header" "Setup Complete" "🎉"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "success" "Dry run completed - no files were created" "✅"
        print_message "info" "Run without --dry-run to create the files"
        return 0
    fi
    
    echo ""
    echo -e "${BOLD}📋 Files Created:${NC}"
    
    if [[ "$CONFIGS_ONLY" != "true" ]]; then
        echo -e "  ${CYAN}GitHub Actions Workflows:${NC}"
        echo "    📄 .github/workflows/ci.yml"
        echo "    📄 .github/workflows/security.yml" 
        echo "    📄 .github/workflows/dependency-update.yml"
        echo ""
    fi
    
    if [[ "$WORKFLOWS_ONLY" != "true" ]]; then
        echo -e "  ${CYAN}Configuration Files:${NC}"
        echo "    📄 .golangci.yml"
        echo "    📄 .codecov.yml"
        echo "    📄 .github/dependabot.yml"
        echo "    📄 Makefile"
        echo ""
    fi
    
    echo -e "${BOLD}🚀 Next Steps:${NC}"
    echo "  1. 📝 Review the generated files"
    echo "  2. 🔧 Customize configurations as needed"
    echo "  3. 🛠️ Run 'make dev-setup' to install development tools"
    echo "  4. 🧪 Run 'make test' to verify everything works"
    echo "  5. 📤 Commit and push the changes to GitHub"
    echo "  6. ✅ Enable GitHub Actions in your repository settings"
    echo ""
    
    echo -e "${BOLD}📚 Useful Commands:${NC}"
    echo "  make help          # Show all available targets"
    echo "  make build         # Build all packages"
    echo "  make test          # Run tests with coverage"
    echo "  make lint          # Run code linting"
    echo "  make security      # Run security scans"
    echo "  make all           # Complete build pipeline"
    echo ""
    
    print_message "success" "GoVel CI/CD setup completed successfully!" "🎉"
}

##
# 🎯 Show setup preview
##
show_setup_preview() {
    print_message "header" "Setup Preview" "🔍"
    
    echo ""
    echo -e "${BOLD}The following will be created:${NC}"
    echo ""
    
    if [[ "$CONFIGS_ONLY" != "true" ]]; then
        echo -e "${CYAN}📁 GitHub Actions Workflows:${NC}"
        echo "  ├── 🎯 ci.yml                    # Main CI pipeline with testing"
        echo "  ├── 🔒 security.yml             # Security scanning (gosec, govulncheck)"
        echo "  └── 🤖 dependency-update.yml    # Automated dependency updates"
        echo ""
    fi
    
    if [[ "$WORKFLOWS_ONLY" != "true" ]]; then
        echo -e "${CYAN}📁 Configuration Files:${NC}"
        echo "  ├── 🔧 .golangci.yml             # 50+ linting rules"
        echo "  ├── 📊 .codecov.yml              # Coverage reporting & gates"
        echo "  ├── 🤖 .github/dependabot.yml   # Weekly dependency updates"
        echo "  └── 🛠️ Makefile                  # Build automation (help, build, test, etc.)"
        echo ""
    fi
    
    echo -e "${BOLD}✨ Features Included:${NC}"
    echo "  • 🧪 Multi-version Go testing (1.21, 1.22, 1.23)"
    echo "  • 🖥️ Cross-platform testing (Linux, macOS, Windows)"
    echo "  • 📦 Smart package change detection"
    echo "  • 🔒 Comprehensive security scanning"
    echo "  • 📊 Code coverage with quality gates (80% threshold)"
    echo "  • 🎨 Code formatting and linting"
    echo "  • 🤖 Automated dependency management"
    echo "  • 🚀 Professional development workflow"
    echo ""
}

##
# 🏁 Main execution function
##
main() {
    print_message "header" "GoVel CI/CD Complete Setup" "🚀"
    print_message "info" "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
    print_message "info" "Project root: $PROJECT_ROOT"
    
    # Parse arguments
    parse_arguments "$@"
    
    # Check prerequisites
    check_prerequisites
    
    # Show setup preview
    show_setup_preview
    
    # Confirm with user unless force mode
    if [[ "$FORCE" != "true" && "$DRY_RUN" != "true" ]]; then
        echo ""
        echo -e "${YELLOW}❓ Do you want to proceed with the CI/CD setup? [y/N]:${NC} "
        read -r response
        
        case "$response" in
            [yY][eE][sS]|[yY])
                print_message "success" "Proceeding with setup..." "✅"
                ;;
            *)
                print_message "info" "Setup cancelled by user" "🚫"
                exit 0
                ;;
        esac
    fi
    
    # Generate workflows
    generate_workflows
    
    # Generate configuration files  
    generate_configs
    
    # Show final summary
    show_final_summary
}

# 🏁 Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi