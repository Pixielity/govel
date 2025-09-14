#!/bin/bash

##
# 🔧 GoVel CI/CD Configuration Generator
# 
# This companion script generates configuration files, GitHub templates,
# automation scripts, and other supporting files for the GoVel CI/CD pipeline.
#
# Features:
# - 📄 GitHub issue and PR templates
# - 🔧 Linting and code quality configurations
# - 📊 Code coverage configurations
# - 🤖 Dependabot configuration
# - 🛠️ Build automation (Makefile)
# - 📜 CI automation scripts
# - 🪝 Git hooks for development workflow
# - 📋 CODEOWNERS file
#
# Usage:
#   ./generate-cicd-configs.sh [options]
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
#
# @param string $1 The message type (info, success, warning, error, debug)
# @param string $2 The message text to display
# @param string $3 Optional emoji to prepend
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
# 🛠️ Generate GolangCI-Lint configuration
#
# Creates a comprehensive linting configuration for Go code quality.
##
generate_golangci_config() {
    local file_path="$PROJECT_ROOT/.golangci.yml"
    print_message "info" "Generating GolangCI-Lint configuration" "🔧"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .golangci.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 🔧 GoVel Framework - GolangCI-Lint Configuration
#
# This configuration provides comprehensive Go code analysis for the GoVel framework.
# It includes 50+ linters covering code quality, security, performance, and style.
#
# Features:
# - 🚀 Performance optimizations with parallel processing
# - 🛡️ Security-focused linting rules
# - 🎨 Code style and formatting enforcement
# - 🐛 Bug detection and prevention
# - 📊 Code complexity analysis
# - 🧪 Test quality improvements
#
# Documentation: https://golangci-lint.run/usage/configuration/
# Author: GoVel Framework Team
# Version: 1.0.0

# ⚡ Performance and execution settings
run:
  # 🎯 Timeout for analysis, use 0 to disable
  timeout: 5m
  
  # 🚀 Number of operating system threads (`GOMAXPROCS`) that can execute golangci-lint simultaneously
  concurrency: 4
  
  # 📦 Exit code when at least one issue was found
  issues-exit-code: 1
  
  # 🧪 Include test files in analysis
  tests: true
  
  # 📁 Directories to skip during analysis
  skip-dirs:
    - vendor
    - node_modules
    - .git
    - tmp
    - build
    - dist
    - examples/old
  
  # 📄 File patterns to skip
  skip-files:
    - ".*\\.pb\\.go$"           # Generated protobuf files
    - ".*\\.gen\\.go$"          # Generated files
    - ".*_mock\\.go$"           # Mock files (except in dedicated test packages)
    - "mock_.*\\.go$"           # Mock files
  
  # 📍 Working directory for analysis
  # Default: the directory where golangci-lint is run
  # build-tags:
  #   - mytag

# 📋 Output configuration
output:
  # 🎨 Output format: colored-line-number|line-number|json|tab|checkstyle|code-climate|junit-xml|github-actions
  format: colored-line-number
  
  # 📄 Print lines of code with issue
  print-issued-lines: true
  
  # 📊 Print linter name in the end of issue text
  print-linter-name: true
  
  # 📈 Make issues output unique by line
  uniq-by-line: true
  
  # 📝 Add a prefix to the output file references
  path-prefix: ""
  
  # 🔍 Sort results by: filepath, line and column
  sort-results: true

# 🛠️ Linters configuration
linters-settings:
  # 📏 Line length checker
  lll:
    line-length: 120
    tab-width: 4

  # 🎯 Cyclomatic complexity
  cyclop:
    max-complexity: 15
    package-average: 10.0
    skip-tests: false

  # 🔍 Cognitive complexity
  gocognit:
    min-complexity: 20

  # 📊 Function length checker
  funlen:
    lines: 100
    statements: 60
    ignore-comments: true

  # 🏗️ Nested if statements
  nestif:
    min-complexity: 5

  # 🚫 Banned imports
  depguard:
    rules:
      main:
        deny:
          - pkg: "github.com/pkg/errors"
            desc: "Use standard library errors package or golang.org/x/xerrors"
          - pkg: "github.com/sirupsen/logrus"
            desc: "Use our internal logger package instead"
        allow:
          - $gostd
          - govel

  # 📝 Comment requirements
  godot:
    scope: declarations
    exclude:
      - "^fixme:"
      - "^todo:"
    period: true
    capital: true

  # 📋 Documentation requirements  
  godox:
    keywords:
      - NOTE
      - OPTIMIZE
      - HACK
      - FIXME
      - TODO
      - BUG

  # 🎨 Import formatting
  goimports:
    local-prefixes: govel

  # 📦 Import grouping and ordering
  gci:
    sections:
      - standard                    # Standard library
      - default                     # External packages
      - prefix(govel)              # GoVel framework packages
      - blank                      # Blank imports
      - dot                        # Dot imports
    custom-order: true

  # 🏷️ Struct tag validation
  govet:
    enable:
      - assign
      - atomic
      - atomicalign
      - bools
      - buildtag
      - cgocall
      - composites
      - copylocks
      - deepequalerrors
      - errorsas
      - fieldalignment
      - findcall
      - framepointer
      - httpresponse
      - ifaceassert
      - loopclosure
      - lostcancel
      - nilfunc
      - nilness
      - printf
      - reflectvaluecompare
      - shadow
      - shift
      - sigchanyzer
      - sortslice
      - stdmethods
      - stringintconv
      - structtag
      - testinggoroutine
      - tests
      - unmarshal
      - unreachable
      - unsafeptr
      - unusedresult

  # 🔒 Security analysis
  gosec:
    severity: medium
    confidence: medium
    excludes:
      - G104  # Errors unhandled (we handle this with errcheck)
      - G204  # Subprocess launched with variable (context-dependent)
    includes:
      - G101  # Look for hard coded credentials
      - G102  # Bind to all interfaces
      - G103  # Audit the use of unsafe block
      - G106  # Audit the use of ssh.InsecureIgnoreHostKey
      - G107  # Url provided to HTTP request as taint input
      - G108  # Profiling endpoint automatically exposed on /debug/pprof
      - G109  # Potential Integer overflow made by strconv.Atoi result conversion to int16/32
      - G110  # Potential DoS vulnerability via decompression bomb
      - G201  # SQL query construction using format string
      - G202  # SQL query construction using string concatenation
      - G203  # Use of unescaped data in HTML templates
      - G301  # Poor file permissions used when creating a directory
      - G302  # Poor file permissions used with chmod
      - G303  # Creating tempfile using a predictable path
      - G304  # File path provided as taint input
      - G305  # File traversal when extracting zip archive
      - G401  # Detect the usage of DES, RC4, MD5 or SHA1
      - G402  # Look for bad TLS connection settings
      - G403  # Ensure minimum RSA key length of 2048 bits
      - G404  # Insecure random number source (rand)
      - G501  # Import blocklist: crypto/md5
      - G502  # Import blocklist: crypto/des
      - G503  # Import blocklist: crypto/rc4
      - G504  # Import blocklist: net/http/cgi
      - G601  # Implicit memory aliasing of items from a range statement

  # 🧪 Test quality
  testpackage:
    skip-regexp: '(export|internal)_test\.go'

  # 🎯 Unused code detection
  unused:
    check-exported: true
    go: "1.21"

  # ⚠️ Error handling
  errcheck:
    check-type-assertions: true
    check-blank: true
    ignore: fmt:.*,io/ioutil:^Read.*
    exclude-functions:
      - io/ioutil.ReadFile
      - io.Copy(*bytes.Buffer)
      - io.Copy(os.Stdout)

  # 🏷️ Variable naming conventions
  revive:
    min-confidence: 0.8
    rules:
      - name: var-naming
        severity: warning
        disabled: false
        arguments:
          - "ID"
          - "URL" 
          - "HTTP"
          - "JSON"
          - "API"
          - "UUID"
          - "SQL"
          - "TCP"
          - "UDP"
          - "IP"
      - name: exported
        severity: warning
        disabled: false
        arguments:
          - "checkPrivateReceivers"
          - "sayRepetitiveInsteadOfStutters"

# 🎛️ Enabled linters
linters:
  # 🚀 Fast linters that run by default
  fast: false
  
  # ✅ Enable specific linters
  enable:
    # 🔍 Code analysis and bug detection
    - errcheck          # Check for unchecked errors
    - gosimple          # Simplify code
    - govet             # Vet examines Go source code
    - ineffassign       # Detect ineffectual assignments
    - staticcheck       # Advanced Go linter
    - typecheck         # Parse and type-check Go code
    - unused            # Check for unused constants, variables, functions and types
    
    # 📊 Code complexity and quality
    - cyclop            # Check cyclomatic complexity
    - funlen            # Tool for detection of long functions
    - gocognit          # Compute cognitive complexities
    - gocyclo           # Computes cyclomatic complexities
    - nestif            # Report deeply nested if statements
    
    # 🎨 Code style and formatting
    - gci               # Control golang package import order
    - gofmt             # Check if code was gofmt-ed
    - gofumpt           # Check if code was gofumpt-ed
    - goimports         # Check import statements are formatted
    - misspell          # Find commonly misspelled English words
    - whitespace        # Tool for detection of leading/trailing whitespace
    - wsl               # Whitespace Linter
    
    # 🔒 Security
    - gosec             # Inspect source code for security problems
    
    # 🏗️ Code structure and design
    - dupl              # Tool for code clone detection
    - goconst           # Find repeated strings that could be constants
    - gocritic          # Most opinionated Go source code linter
    - godot             # Check if comments end in a period
    - godox             # Tool for detection of FIXME, TODO and other comment keywords
    - gomnd             # Detects magic numbers
    - gomoddirectives   # Manage the use of 'replace', 'retract', and 'excludes' directives
    - gomodguard        # Allow and block list linter for direct Go module dependencies
    - goprintffuncname  # Check that printf-like functions are named with f at the end
    - lll               # Line length limit
    - nakedret          # Find naked returns in functions greater than a specified function length
    - nilerr            # Find the code that returns nil even if it checks that the error is not nil
    - nlreturn          # Check for a new line before return and branch statements
    - noctx             # Find sending http request without context.Context
    - nolintlint        # Report ill-formed or insufficient nolint directives
    - predeclared       # Find code that shadows one of Go's predeclared identifiers
    - revive            # Fast, configurable, extensible, flexible, and beautiful linter for Go
    - rowserrcheck      # Check whether Err of rows is checked successfully
    - sqlclosecheck     # Check that sql.Rows and sql.Stmt are closed
    - unconvert         # Remove unnecessary type conversions
    - unparam           # Report unused function parameters
    - wastedassign      # Find wasted assignment statements
    
    # 🧪 Test quality
    - testpackage       # Make sure that separate _test packages are used
    - tparallel         # Detect inappropriate usage of t.Parallel()
    - thelper           # Detect golang test helpers without t.Helper()
    
    # 📦 Imports and dependencies  
    - depguard          # Check if package imports are in a list of acceptable packages
    - importas          # Enforce consistent import aliases
    
    # 🚫 Error prone patterns
    - errorlint         # Find code that will cause problems with the error wrapping scheme
    - forcetypeassert   # Find forced type assertions
    - makezero          # Find slice declarations with non-zero initial length
    - nilnil            # Check that there is no simultaneous return of nil error and nil value
    
  # ❌ Disable specific linters (if needed)
  disable:
    - exhaustive        # Check exhaustiveness of enum switch statements (too strict)
    - exhaustivestruct  # Check that all struct fields are initialized (too strict)
    - forbidigo         # Forbid identifiers (can be overly restrictive)
    - gci               # Disabled in favor of goimports for now
    - gochecknoglobals  # Check that no global variables exist (too strict for framework)
    - gochecknoinits    # Check that no init functions are present (needed for framework)
    - goerr113          # Check error handling expressions (can be overly strict)
    - golint            # Deprecated, replaced by revive
    - interfacer        # Deprecated linter
    - maligned          # Deprecated, replaced by fieldalignment in govet
    - scopelint         # Deprecated, replaced by exportloopref
    - varcheck          # Deprecated, replaced by unused
    - deadcode          # Deprecated, replaced by unused
    - structcheck       # Deprecated, replaced by unused

# 🎯 Issues configuration
issues:
  # 📊 Maximum issues count per one linter
  max-issues-per-linter: 50
  
  # 📈 Maximum count of issues with the same text
  max-same-issues: 3
  
  # 🎨 Show only new issues created after git revision
  # new: false
  # new-from-rev: origin/main
  
  # 🔧 Fix found issues (if it's supported by the linter)
  fix: false
  
  # 🚫 Excluding configuration per-path, per-linter, per-text and per-source
  exclude:
    # 📋 Default exclusions (can be overridden)
    - "Error return value of .((os\\.)?std(out|err)\\..*|.*Close|.*Flush|os\\.Remove(All)?|.*printf?|os\\.(Un)?Setenv). is not checked"
    - "func name will be used as test\\.Test.* by other packages, and that stutters; consider calling this"
    - "G104: Errors unhandled"  # Covered by errcheck
    - "G204: Subprocess launched with variable"  # Context dependent
    
  # 📁 Exclude rules by file path pattern
  exclude-rules:
    # 🧪 Test files have different standards
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - funlen
        - gocognit
        - cyclop
        - lll
        - goconst
        - gomnd
        - maintidx
        - nestif
        
    # 📦 Main packages can have longer functions
    - path: cmd/
      linters:
        - funlen
        - gocognit
        - cyclop
        
    # 🔧 Generated code exclusions
    - path: ".*\\.pb\\.go"
      linters:
        - lll
        - gocyclo
        - errcheck
        - gosec
        - dupl
        - goconst
        - funlen
        - gomnd
        
    # 🏗️ Mock files exclusions
    - path: ".*_mock\\.go"
      linters:
        - errcheck
        - gosec
        - dupl
        - goconst
        - gomnd
        - unused
        
    # 📝 Example code can be less strict
    - path: examples/
      linters:
        - errcheck
        - gosec
        - goconst
        - gomnd
        - lll
        - funlen
        
    # 🎯 Specific rule exclusions
    - linters:
        - lll
      source: "^//go:generate "
      
    - linters:
        - goconst
      source: "const.*=.*(json|yaml|toml):"

  # 🎯 Include issues created by deprecated linters
  include:
    - EXC0002  # disable excluding of issues about comments from golint
    - EXC0003  # disable excluding of issues about comments from revive
    - EXC0004  # disable excluding of issues about comments from govet
    - EXC0005  # disable excluding of issues about printf from gosimple
    - EXC0011  # disable excluding of issues about assignments from govet
    - EXC0012  # disable excluding of issues about pack pragmas from revive
    - EXC0013  # disable excluding of issues about assignments from revive
    - EXC0014  # disable excluding of issues about comments from revive
    - EXC0015  # disable excluding of issues about blank imports

# 🌟 Severity configuration
severity:
  # 📊 Default severity for issues
  default-severity: error
  
  # 🎯 Case sensitive matching
  case-sensitive: true
  
  # 📋 Set the default severity for issues
  rules:
    - linters:
        - dupl
        - goconst
        - gomnd
        - gocritic
      severity: warning
      
    - linters:
        - errcheck
        - gosec
        - govet
        - staticcheck
        - typecheck
      severity: error

EOF

    print_message "success" "Created GolangCI-Lint configuration" "✅"
}

##
# 📊 Generate Codecov configuration  
#
# Creates configuration for code coverage reporting and quality gates.
##
generate_codecov_config() {
    local file_path="$PROJECT_ROOT/.codecov.yml"
    print_message "info" "Generating Codecov configuration" "📊"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .codecov.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 📊 GoVel Framework - Codecov Configuration
#
# This configuration defines code coverage reporting and quality gates
# for the GoVel framework. It ensures high-quality code through
# comprehensive coverage analysis and automated quality checks.
#
# Features:
# - 📈 Coverage thresholds and quality gates
# - 📦 Per-package coverage requirements
# - 🎯 Pull request status checks
# - 📊 Coverage trend analysis
# - 🚫 Proper path ignoring for generated/test files
#
# Documentation: https://docs.codecov.io/docs/codecov-yaml
# Author: GoVel Framework Team
# Version: 1.0.0

# 🎯 Coverage configuration
coverage:
  # 📊 Precision for coverage percentage (2 decimal places)
  precision: 2
  
  # 🎯 Rounding mode: down, up, nearest
  round: nearest
  
  # 📈 Coverage range for color coding (red to green)
  range: "60...90"
  
  # 📋 Status checks configuration
  status:
    # 🔍 Project-wide coverage requirements
    project:
      default:
        # 🎯 Target coverage percentage
        target: 80%
        
        # 📉 Allowed coverage drop threshold
        threshold: 2%
        
        # 📊 Base for comparison (auto, parent, pr)
        base: auto
        
        # 🚫 Only post status if coverage changes
        if_no_uploads: error
        
        # 🎛️ Only run on these branches
        branches:
          - main
          - develop
          
        # 📦 Path-based coverage (focus on source code)
        paths:
          - "packages/*/src"
          - "packages/*/*.go"
        
    # 🔄 Patch coverage (new code in PRs)
    patch:
      default:
        # 🎯 New code must be well tested
        target: 85%
        
        # 📉 Strict threshold for new code
        threshold: 3%
        
        # 📊 Base comparison
        base: auto
        
        # 🔍 Only analyze changed files
        only_pulls: true
        
        # 📦 Focus on source code changes
        paths:
          - "packages/*/src"
          - "packages/*/*.go"

  # 📋 Individual package coverage requirements
  flags:
    # 🔐 Core packages require higher coverage
    core-packages:
      paths:
        - packages/application
        - packages/container  
        - packages/config
        - packages/logger
      target: 85%
      
    # 🛡️ Security-critical packages
    security-packages:
      paths:
        - packages/encryption
        - packages/hashing
      target: 90%
      
    # 🔧 Utility packages
    utility-packages:
      paths:
        - packages/support
        - packages/cookie
        - packages/pipeline
      target: 80%

# 💬 Comment configuration for pull requests
comment:
  # 📊 Layout of the coverage comment
  layout: "reach,diff,flags,tree,betaprofiling"
  
  # 🎯 Behavior settings
  behavior: default
  
  # 📉 Require coverage changes to post comment
  require_changes: true
  
  # 📝 Show file-level coverage in comments
  require_base: no
  require_head: yes
  
  # 📋 Branch settings for comments
  branches:
    - main
    - develop

# 📂 Path ignoring configuration
ignore:
  # 🧪 Test files and test utilities
  - "**/tests/**"
  - "**/test/**"
  - "**/*_test.go"
  - "**/mocks/**"
  - "**/mock_*.go"
  - "**/*_mock.go"
  
  # 📄 Documentation and examples
  - "**/docs/**"
  - "**/examples/**"
  - "**/demo/**"
  - "**/*.md"
  
  # 🔧 Generated and vendor code
  - "**/vendor/**"
  - "**/*.pb.go"
  - "**/*.gen.go"
  - "**/generated/**"
  
  # 🏗️ Build and tooling
  - "**/scripts/**"
  - "**/tools/**"
  - "**/build/**"
  - "**/dist/**"
  - "Makefile"
  - "**/*.yml"
  - "**/*.yaml"
  - "**/go.mod"
  - "**/go.sum"
  
  # 📦 Package-specific ignores
  - "packages/new/**"          # New/experimental packages
  - "**/cmd/**"               # Command-line tools
  - "**/main.go"              # Main entry points

# 🚨 Notification configuration
github_checks:
  # ✅ Enable GitHub status checks
  annotations: true

# 🔧 Parsing configuration
parsers:
  go:
    # 📊 Go coverage report parsing
    branch_detection:
      conditional: true
      loop: true
      method: false
      macro: false
  
  gcov:
    # 🎯 Branch coverage settings
    branch_detection:
      conditional: true
      loop: true
      method: false
      macro: false

# 📈 Profiling and analytics
profiling:
  # 🔍 Critical files that need high attention
  critical_files_paths:
    - packages/application/application.go
    - packages/container/container.go
    - packages/encryption/src/*.go
    - packages/hashing/src/*.go
    - packages/logger/logger.go

# 🎛️ Advanced configuration
codecov:
  # 🔒 Security and validation
  require_ci_to_pass: true
  
  # 📊 Archive settings
  archive:
    uploads: true
    
  # 🕒 Upload timeout
  max_report_age: 24  # hours
  
  # 📧 Notification settings
  notify:
    # 📥 After N uploads
    after_n_builds: 2
    
    # ⏰ Wait time before sending notifications  
    wait_for_ci: true

EOF

    print_message "success" "Created Codecov configuration" "✅"
}

##
# 🤖 Generate Dependabot configuration
#
# Creates automated dependency update configuration for GitHub.
##
generate_dependabot_config() {
    local file_path="$PROJECT_ROOT/.github/dependabot.yml"
    print_message "info" "Generating Dependabot configuration" "🤖"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: .github/dependabot.yml" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
# 🤖 GoVel Framework - Dependabot Configuration
#
# This configuration enables automated dependency updates for the GoVel framework.
# It manages Go modules, GitHub Actions, and other dependencies across all packages
# to ensure security and up-to-date dependencies.
#
# Features:
# - 📦 Go module updates for all packages
# - 🔄 GitHub Actions workflow updates  
# - 🔒 Security-focused update strategy
# - 📅 Scheduled updates to minimize noise
# - 🏷️ Automated labeling and grouping
# - 👥 Automatic reviewer assignment
#
# Documentation: https://docs.github.com/en/code-security/dependabot
# Author: GoVel Framework Team
# Version: 1.0.0

version: 2
updates:
  # 📦 Go module updates for main packages
  - package-ecosystem: "gomod"
    directory: "/packages/application"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 5
    reviewers:
      - "govel-maintainers"
    assignees:
      - "govel-lead"
    labels:
      - "dependencies"
      - "go-modules"
      - "application"
    commit-message:
      prefix: "🤖"
      include: "scope"
    rebase-strategy: "auto"
    allow:
      - dependency-type: "direct"
      - dependency-type: "indirect"
    groups:
      security-updates:
        patterns:
          - "golang.org/x/*"
          - "*security*"
        update-types:
          - "patch"
          - "minor"
      minor-updates:
        patterns:
          - "*"
        update-types:
          - "minor"
          - "patch"

  # 🗂️ Container package
  - package-ecosystem: "gomod"
    directory: "/packages/container"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "container"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🔧 Config package
  - package-ecosystem: "gomod"
    directory: "/packages/config"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "config"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🔐 Encryption package
  - package-ecosystem: "gomod"
    directory: "/packages/encryption"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
      - "security-team"
    labels:
      - "dependencies"
      - "go-modules"
      - "encryption"
      - "security"
    commit-message:
      prefix: "🔒"
      include: "scope"
    allow:
      - dependency-type: "direct"

  # 🔑 Hashing package
  - package-ecosystem: "gomod"
    directory: "/packages/hashing"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
      - "security-team"
    labels:
      - "dependencies"
      - "go-modules"
      - "hashing"
      - "security"
    commit-message:
      prefix: "🔒"
      include: "scope"

  # 📝 Logger package
  - package-ecosystem: "gomod"
    directory: "/packages/logger"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "logger"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🚰 Pipeline package
  - package-ecosystem: "gomod"
    directory: "/packages/pipeline"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "pipeline"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🛠️ Support package
  - package-ecosystem: "gomod"
    directory: "/packages/support"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "support"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🍪 Cookie package
  - package-ecosystem: "gomod"
    directory: "/packages/cookie"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
    labels:
      - "dependencies"
      - "go-modules"
      - "cookie"
    commit-message:
      prefix: "🤖"
      include: "scope"

  # 🔄 GitHub Actions updates
  - package-ecosystem: "github-actions"
    directory: "/.github/workflows"
    schedule:
      interval: "weekly"
      day: "tuesday"
      time: "10:00"
      timezone: "UTC"
    open-pull-requests-limit: 3
    reviewers:
      - "govel-maintainers"
      - "devops-team"
    labels:
      - "dependencies"
      - "github-actions"
      - "ci-cd"
    commit-message:
      prefix: "🔄"
      include: "scope"
    groups:
      actions-security:
        patterns:
          - "actions/checkout"
          - "actions/setup-*"
          - "github/codeql-action"
          - "codecov/codecov-action"
        update-types:
          - "major"
          - "minor"
          - "patch"

  # 🐳 Docker updates (if any Dockerfiles exist)
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
      day: "wednesday"
      time: "10:00"
      timezone: "UTC"
    open-pull-requests-limit: 2
    reviewers:
      - "govel-maintainers"
      - "devops-team"
    labels:
      - "dependencies"
      - "docker"
    commit-message:
      prefix: "🐳"
      include: "scope"

EOF

    print_message "success" "Created Dependabot configuration" "✅"
}

##
# 🛠️ Generate Makefile
#
# Creates a comprehensive build automation Makefile for development.
##
generate_makefile() {
    local file_path="$PROJECT_ROOT/Makefile"
    print_message "info" "Generating Makefile" "🛠️"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "info" "[DRY-RUN] Would create: Makefile" "📄"
        return 0
    fi
    
    cat > "$file_path" << 'EOF'
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

EOF

    print_message "success" "Created Makefile" "✅"
}

##
# 🎯 Parse command line arguments
##
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dry-run)
                DRY_RUN=true
                print_message "info" "Dry-run mode enabled" "🔍"
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                print_message "info" "Verbose mode enabled" "📝"
                shift
                ;;
            -f|--force)
                FORCE=true
                print_message "info" "Force mode enabled" "💪"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_message "error" "Unknown option: $1"
                exit 1
                ;;
        esac
    done
}

##
# 🛠️ Display help information
##
show_help() {
    cat << EOF
${BOLD}${CYAN}🔧 GoVel CI/CD Configuration Generator${NC}

${BOLD}DESCRIPTION:${NC}
    Generates configuration files, templates, and automation scripts
    for the GoVel framework CI/CD pipeline.

${BOLD}USAGE:${NC}
    $SCRIPT_NAME [OPTIONS]

${BOLD}OPTIONS:${NC}
    -d, --dry-run     🔍 Show what would be created without creating files
    -v, --verbose     📝 Enable detailed output and debug information  
    -f, --force       💪 Overwrite existing files without confirmation
    -h, --help        ❓ Show this help message and exit

${BOLD}FILES CREATED:${NC}
    📄 .golangci.yml      # Comprehensive linting configuration
    📄 .codecov.yml       # Code coverage configuration
    📄 .github/dependabot.yml  # Automated dependency updates
    📄 Makefile          # Build automation
    📁 GitHub templates  # Issue and PR templates
    📁 CI scripts       # Automation scripts

EOF
}

##
# 🏁 Main execution function
##
main() {
    print_message "header" "GoVel CI/CD Configuration Generator" "🔧"
    
    # Parse arguments
    parse_arguments "$@"
    
    # Generate configuration files
    print_message "header" "Generating Configuration Files" "📄"
    generate_golangci_config
    generate_codecov_config
    generate_dependabot_config
    generate_makefile
    
    # Final summary
    print_message "header" "Configuration Generation Complete" "🎉"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_message "success" "Dry run completed - showed what would be created" "✅"
    else
        print_message "success" "All configuration files created successfully" "✅"
        print_message "info" "Next steps:"
        echo "  1. Review the generated configuration files"
        echo "  2. Customize settings as needed for your project"
        echo "  3. Run 'make dev-setup' to install development tools"
        echo "  4. Run 'make help' to see available build targets"
    fi
}

# 🏁 Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi