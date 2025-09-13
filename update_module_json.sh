#!/bin/bash

# GoVel Module.json Generator Script
# This script regenerates all module.json files with a standardized format

set -e  # Exit on any error

# Configuration
PROJECT_ROOT="/Users/akouta/Projects/govel"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Package metadata lookup functions
get_package_description() {
    local pkg="$1"
    case "$pkg" in
        "application") echo "Core application framework with lifecycle management and service container integration for GoVel" ;;
        "config") echo "Configuration management system with support for multiple formats and environment-based loading for GoVel" ;;
        "container") echo "Dependency injection container with service binding and resolution capabilities for GoVel" ;;
        "cookie") echo "Laravel-compatible cookie management package for GoVel framework" ;;
        "encryption") echo "Encryption and decryption services with multiple cipher support for GoVel framework" ;;
        "hashing") echo "Password hashing and verification utilities with multiple algorithms for GoVel framework" ;;
        "logger") echo "Logging infrastructure with multiple channels and formatters for GoVel framework" ;;
        "pipeline") echo "Pipeline processing system for middleware and data transformation in GoVel framework" ;;
        "support") echo "Core utility functions and helpers for GoVel framework development" ;;
        "types") echo "Type definitions and interfaces for GoVel framework components" ;;
        "new/bus") echo "Command bus pattern implementation with queuing support for GoVel framework" ;;
        "new/cache") echo "Caching system with multiple drivers and Laravel compatibility for GoVel framework" ;;
        "new/filesystem") echo "File system operations and storage abstraction for GoVel framework" ;;
        "new/health_check") echo "Health check monitoring system for GoVel applications" ;;
        "new/middleware") echo "HTTP middleware collection for GoVel web applications" ;;
        "new/redis") echo "Redis client and utilities for GoVel framework" ;;
        "new/routing") echo "HTTP routing system with Laravel-compatible features for GoVel framework" ;;
        "new/session") echo "Session management with multiple storage drivers for GoVel framework" ;;
        "new/translation") echo "Internationalization and localization system for GoVel framework" ;;
        "new/webserver") echo "Web server abstraction with multiple adapter support for GoVel framework" ;;
        "new/exceptions") echo "Exception handling and error management for GoVel framework" ;;
        "new/ignition") echo "Debug and error page utilities for GoVel framework development" ;;
        *) echo "GoVel framework package" ;;
    esac
}

get_package_category() {
    local pkg="$1"
    case "$pkg" in
        "application"|"config"|"container"|"pipeline"|"types") echo "core" ;;
        "cookie"|"new/middleware"|"new/routing"|"new/session"|"new/webserver") echo "web" ;;
        "encryption"|"hashing") echo "security" ;;
        "logger"|"new/cache"|"new/filesystem"|"new/redis"|"new/exceptions") echo "infrastructure" ;;
        "support") echo "utility" ;;
        "new/bus") echo "architecture" ;;
        "new/health_check") echo "monitoring" ;;
        "new/translation") echo "localization" ;;
        "new/ignition") echo "development" ;;
        *) echo "utility" ;;
    esac
}

get_package_keywords() {
    local pkg="$1"
    case "$pkg" in
        "application") echo "\"application\", \"lifecycle\", \"bootstrap\", \"service-provider\"" ;;
        "config") echo "\"config\", \"configuration\", \"environment\", \"settings\"" ;;
        "container") echo "\"container\", \"dependency-injection\", \"service-binding\", \"ioc\"" ;;
        "cookie") echo "\"cookie\", \"http\", \"web\", \"csrf\", \"security\", \"encryption\", \"samesite\"" ;;
        "encryption") echo "\"encryption\", \"decryption\", \"security\", \"aes\", \"cipher\"" ;;
        "hashing") echo "\"hashing\", \"password\", \"bcrypt\", \"argon2\", \"security\"" ;;
        "logger") echo "\"logging\", \"logger\", \"infrastructure\", \"monitoring\"" ;;
        "pipeline") echo "\"pipeline\", \"middleware\", \"processing\", \"chain\"" ;;
        "support") echo "\"utility\", \"helpers\", \"support\", \"common\"" ;;
        "types") echo "\"types\", \"interfaces\", \"definitions\", \"contracts\"" ;;
        "new/bus") echo "\"bus\", \"command\", \"queue\", \"architecture\", \"pattern\"" ;;
        "new/cache") echo "\"cache\", \"caching\", \"storage\", \"performance\"" ;;
        "new/filesystem") echo "\"filesystem\", \"storage\", \"file\", \"disk\"" ;;
        "new/health_check") echo "\"health\", \"monitoring\", \"check\", \"status\"" ;;
        "new/middleware") echo "\"middleware\", \"http\", \"web\", \"request\"" ;;
        "new/redis") echo "\"redis\", \"cache\", \"storage\", \"nosql\"" ;;
        "new/routing") echo "\"routing\", \"http\", \"web\", \"router\"" ;;
        "new/session") echo "\"session\", \"web\", \"state\", \"storage\"" ;;
        "new/translation") echo "\"translation\", \"i18n\", \"localization\", \"language\"" ;;
        "new/webserver") echo "\"webserver\", \"http\", \"server\", \"web\", \"fiber\"" ;;
        "new/exceptions") echo "\"exceptions\", \"error\", \"handling\", \"debug\"" ;;
        "new/ignition") echo "\"ignition\", \"debug\", \"error-page\", \"development\"" ;;
        *) echo "\"utility\"" ;;
    esac
}

get_package_providers() {
    local pkg="$1"
    case "$pkg" in
        "application") echo "\"ApplicationServiceProvider\"" ;;
        "config") echo "\"ConfigServiceProvider\"" ;;
        "container") echo "\"ContainerServiceProvider\"" ;;
        "cookie") echo "\"CookieServiceProvider\"" ;;
        "encryption") echo "\"EncryptionServiceProvider\"" ;;
        "hashing") echo "\"HashingServiceProvider\"" ;;
        "logger") echo "\"LoggerServiceProvider\"" ;;
        "pipeline") echo "\"PipelineServiceProvider\"" ;;
        "new/bus") echo "\"BusServiceProvider\"" ;;
        "new/cache") echo "\"CacheServiceProvider\"" ;;
        "new/filesystem") echo "\"FilesystemServiceProvider\"" ;;
        "new/health_check") echo "\"HealthCheckServiceProvider\"" ;;
        "new/middleware") echo "\"MiddlewareServiceProvider\"" ;;
        "new/redis") echo "\"RedisServiceProvider\"" ;;
        "new/routing") echo "\"RoutingServiceProvider\"" ;;
        "new/session") echo "\"SessionServiceProvider\"" ;;
        "new/translation") echo "\"TranslationServiceProvider\"" ;;
        "new/webserver") echo "\"WebserverServiceProvider\"" ;;
        *) echo "" ;;
    esac
}

# Function to find all packages
find_packages() {
    local packages=()
    
    # Find direct packages
    for dir in "$PROJECT_ROOT/packages"/*; do
        if [ -d "$dir" ] && [ "$(basename "$dir")" != "new" ]; then
            packages+=("$(basename "$dir")")
        fi
    done
    
    # Find new/* packages
    if [ -d "$PROJECT_ROOT/packages/new" ]; then
        for dir in "$PROJECT_ROOT/packages/new"/*; do
            if [ -d "$dir" ]; then
                packages+=("new/$(basename "$dir")")
            fi
        done
    fi
    
    printf '%s\n' "${packages[@]}"
}

# Function to analyze package dependencies
analyze_dependencies() {
    local package_path="$1"
    local dependencies=()
    
    # Find all Go files and extract govel imports
    if [ -d "$package_path/src" ]; then
        local imports=$(find "$package_path/src" -name "*.go" -type f -exec grep -h "^[[:space:]]*\"govel/" {} \; 2>/dev/null | \
                       sed 's/^[[:space:]]*//' | \
                       sed 's/"//g' | \
                       sort -u)
        
        for import in $imports; do
            # Extract package name from import path
            if [[ $import =~ ^govel/(.+)$ ]]; then
                local dep_pkg="${BASH_REMATCH[1]}"
                
                # Skip self-references and subpackages
                local current_pkg=$(basename "$package_path")
                if [[ "$package_path" == *"/new/"* ]]; then
                    current_pkg="new/$(basename "$package_path")"
                fi
                
                if [[ "$dep_pkg" != "$current_pkg" ]] && [[ "$dep_pkg" != "$current_pkg"/* ]]; then
                    # Convert to @govel format
                    local govel_dep="@govel/${dep_pkg//\//-}"
                    if [[ "$dep_pkg" == "new/"* ]]; then
                        govel_dep="@govel/${dep_pkg#new/}"
                    else
                        govel_dep="@govel/${dep_pkg}"
                    fi
                    dependencies+=("$govel_dep")
                fi
            fi
        done
    fi
    
    # Remove duplicates and return
    printf '%s\n' "${dependencies[@]}" | sort -u
}

# Function to generate dependencies JSON
generate_dependencies_json() {
    local deps=("$@")
    local json="{"
    local first=true
    
    for dep in "${deps[@]}"; do
        if [ "$first" = false ]; then
            json+=", "
        fi
        json+="\"$dep\": \"^1.0.0\""
        first=false
    done
    
    json+="}"
    echo "$json"
}

# Function to create module.json for a package
create_module_json() {
    local package_name="$1"
    local package_path="$2"
    
    print_status "Generating module.json for $package_name..."
    
    # Get package metadata using functions
    local description=$(get_package_description "$package_name")
    local category=$(get_package_category "$package_name")
    local keywords=$(get_package_keywords "$package_name")
    local providers=$(get_package_providers "$package_name")
    
    # Analyze dependencies
    local deps_array=($(analyze_dependencies "$package_path"))
    local dependencies_json=$(generate_dependencies_json "${deps_array[@]}")
    
    # Handle providers array
    local providers_array=""
    if [ -n "$providers" ]; then
        providers_array="[$providers]"
    else
        providers_array="[]"
    fi
    
    # Create module.json content
    cat > "$package_path/module.json" << EOF
{
  "name": "@govel/$package_name",
  "version": "1.0.0",
  "description": "$description",
  "author": "GoVel Framework Team",
  "license": "MIT",
  "keywords": [
    "go", 
    "golang", 
    $keywords,
    "laravel", 
    "govel", 
    "framework", 
    "module"
  ],
  "repository": {
    "type": "git",
    "url": "https://github.com/govel-framework/govel.git",
    "directory": "packages/$package_name"
  },
  "bugs": {
    "url": "https://github.com/govel-framework/govel/issues"
  },
  "homepage": "https://github.com/govel-framework/govel/tree/main/packages/$package_name#readme",
  "dependencies": $dependencies_json,
  "scripts": {
    "test": "go test -v ./...",
    "test:coverage": "go test -v -cover ./...",
    "test:race": "go test -v -race ./...",
    "build": "go build ./...",
    "lint": "golangci-lint run",
    "fmt": "go fmt ./...",
    "vet": "go vet ./...",
    "mod:tidy": "go mod tidy",
    "mod:verify": "go mod verify",
    "clean": "go clean -cache -testcache -modcache"
  },
  "hooks": {
    "pre-install": [],
    "post-install": [
      "go mod tidy",
      "go mod download"
    ],
    "pre-update": [],
    "post-update": [
      "go mod tidy",
      "go mod download"
    ],
    "pre-build": [
      "go fmt ./...",
      "go vet ./..."
    ],
    "post-build": [],
    "pre-test": [
      "go mod verify"
    ],
    "post-test": [],
    "pre-publish": [
      "go test ./...",
      "go fmt ./...",
      "go vet ./...",
      "golangci-lint run"
    ],
    "post-publish": []
  },
  "engines": {
    "go": ">=1.19"
  },
  "files": [
    "src/",
    "README.md",
    "LICENSE",
    "go.mod",
    "go.sum"
  ],
  "govel": {
    "type": "package",
    "category": "$category",
    "providers": $providers_array
  }
}
EOF

    echo "  ✅ Created module.json for $package_name"
}

# Function to backup and remove existing module.json files
backup_and_remove_existing() {
    print_status "Backing up and removing existing module.json files..."
    
    local backup_dir="$PROJECT_ROOT/module_json_backup_$(date +%s)"
    mkdir -p "$backup_dir"
    
    local count=0
    while IFS= read -r -d '' file; do
        local relative_path="${file#$PROJECT_ROOT/packages/}"
        local backup_path="$backup_dir/$relative_path"
        
        # Create backup directory structure
        mkdir -p "$(dirname "$backup_path")"
        
        # Copy to backup
        cp "$file" "$backup_path"
        
        # Remove original
        rm "$file"
        
        ((count++))
    done < <(find "$PROJECT_ROOT/packages" -name "module.json" -type f -print0)
    
    print_success "Backed up $count module.json files to $backup_dir"
}

# Function to show help
show_help() {
    cat << EOF
GoVel Module.json Generator Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -h, --help      Show this help message
    --backup-only   Only backup existing files without regenerating
    --no-backup     Skip backup of existing files
    --dry-run       Show what would be generated without making changes

DESCRIPTION:
    Regenerates all module.json files in the GoVel project with standardized format.
    Automatically detects dependencies by analyzing Go import statements.

FEATURES:
    • Standardized module.json format across all packages
    • Automatic dependency detection from Go imports
    • Package-specific metadata (descriptions, keywords, categories)
    • Backup of existing files before regeneration
    • Support for both core and new/* packages

EXAMPLES:
    $0                  # Generate all module.json files
    $0 --dry-run        # Preview what would be generated
    $0 --no-backup      # Generate without backing up existing files

EOF
}

# Main function
main() {
    print_status "🔄 GoVel Module.json Generator"
    print_status "Standardizing module.json files across all packages"
    echo ""
    
    # Change to project root
    cd "$PROJECT_ROOT" || {
        print_error "Cannot change to project root: $PROJECT_ROOT"
        exit 1
    }
    
    # Find all packages
    local packages=($(find_packages))
    print_status "Found ${#packages[@]} packages to process"
    
    # Backup existing files unless --no-backup is specified
    if [ "$SKIP_BACKUP" != "true" ]; then
        backup_and_remove_existing
        echo ""
    fi
    
    # Generate new module.json files
    print_status "Generating new module.json files..."
    echo ""
    
    for package in "${packages[@]}"; do
        local package_path="$PROJECT_ROOT/packages/$package"
        
        if [ -d "$package_path" ]; then
            create_module_json "$package" "$package_path"
        else
            print_warning "Package directory not found: $package_path"
        fi
    done
    
    echo ""
    print_success "🎉 Successfully generated ${#packages[@]} module.json files!"
    echo ""
    
    print_status "✅ What was accomplished:"
    echo "  • Standardized format across all packages"
    echo "  • Automatic dependency detection from Go imports" 
    echo "  • Package-specific metadata and categorization"
    echo "  • Consistent repository and bug tracking URLs"
    echo "  • Standardized npm-style scripts"
    echo ""
    
    print_status "🔧 Next steps:"
    echo "  1. Review generated module.json files"
    echo "  2. Adjust package descriptions or keywords if needed"
    echo "  3. Commit the changes to version control"
    echo ""
    
    print_status "📁 Packages processed:"
    for package in "${packages[@]}"; do
        echo "  • $package"
    done
}

# Command line argument handling
SKIP_BACKUP=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --backup-only)
            backup_and_remove_existing
            exit 0
            ;;
        --no-backup)
            SKIP_BACKUP=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            print_status "🔍 DRY RUN MODE - No files will be modified"
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
done

# Run main function
if [ "$DRY_RUN" = "true" ]; then
    print_status "This would generate module.json files for the following packages:"
    packages=($(find_packages))
    for package in "${packages[@]}"; do
        echo "  • $package"
    done
else
    main "$@"
fi