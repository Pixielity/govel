# 🚀 GoVel CI/CD - Quick Reference Guide

## 🎯 Most Used Commands

### Root-Level Operations
```bash
# Quick start
make help                    # Show all available commands
make info                    # Project information
make all                     # Run complete CI pipeline

# Development workflow
make test                    # Run all tests with coverage
make lint                    # Lint all packages
make security               # Security scans
make build                  # Build all packages

# Quality checks
make format                 # Format all code
make vet                    # Run go vet
make coverage-report        # Coverage summary
```

### Package-Level Operations
```bash
cd packages/[package-name]

# Quick start
make help                   # Show package commands
make all                    # Complete package pipeline

# Development
make build                  # Build package
make test                   # Test with coverage
make lint                   # Package linting
make security              # Package security scan
```

---

## ⚡ Quick Workflows

### 🔄 Daily Development
```bash
# 1. Start working on feature
git checkout -b feature/my-feature

# 2. Make changes, then quick check
make test lint

# 3. Full pipeline before commit
make all

# 4. Commit and push
git add . && git commit -m "feat: description"
git push origin feature/my-feature
```

### 📦 Package Development
```bash
# Navigate to your package
cd packages/my-package

# Quick iteration cycle
make build test           # Build and test
make lint fmt            # Quality checks
make all                 # Full pipeline

# Return to root for commit
cd ../.. && git add . && git commit -m "feat(my-package): description"
```

### 🚀 Pre-Release Workflow
```bash
# Complete pre-release pipeline
make pre-release

# Or step by step
make deps-verify
make security
make test
make coverage-report
```

---

## 🛠️ Setup Commands

### Initial Setup
```bash
# Run CI/CD setup (first time)
./setup-cicd.sh

# Add per-package CI/CD
./setup-per-package-cicd.sh

# Setup development environment
make dev-setup
```

### New Package Setup
```bash
# 1. Create package structure
mkdir -p packages/new-package/src
cd packages/new-package/src
go mod init govel/new-package

# 2. Return to root and setup CI/CD
cd ../../..
./setup-per-package-cicd.sh

# 3. Verify setup
cd packages/new-package
make help && make build
```

---

## 🔍 Troubleshooting Commands

### Build Issues
```bash
# Dependency problems
make deps-download          # Root level
cd packages/[package] && make deps  # Package level

# Module issues
cd packages/[package]/src
go mod tidy
go mod verify
```

### Test Issues
```bash
# Race conditions
go test -race ./...

# Verbose output
make test-verbose

# Platform-specific
GOOS=linux go test ./...
```

### Coverage Issues
```bash
# Generate HTML report
make coverage-html
open coverage/coverage.html

# Text summary
make coverage-report

# Per-package coverage
cd packages/[package] && make coverage
```

### Linting Issues
```bash
# Auto-fix where possible
make lint-fix

# Check specific issues
golangci-lint run --verbose

# Format code
make format
```

---

## 📊 Monitoring & Reports

### Coverage Tracking
```bash
make coverage-merge         # Merge all coverage reports
make coverage-html          # Generate HTML report
make coverage-report        # Show summary
```

### Security Monitoring
```bash
make security               # All security scans
make security-gosec         # Static analysis
make security-govulncheck   # Vulnerability check
```

### Project Health
```bash
make info                   # Project overview
make report-summary         # Comprehensive report
make check-all             # All quality checks
```

---

## 🎛️ Configuration Files

### Key Files to Know
```
.github/workflows/ci.yml           # Main CI pipeline
.github/workflows/security.yml     # Security scanning
.golangci.yml                      # Root linting rules
.codecov.yml                       # Coverage configuration
Makefile                           # Root automation
packages/[name]/.golangci.yml      # Package linting
packages/[name]/Makefile           # Package automation
```

### Common Customizations
```bash
# Adjust coverage threshold (in .codecov.yml)
target: 80%

# Add linting rules (in .golangci.yml)
linters:
  enable:
    - newlinter

# Modify build targets (in Makefile)
custom-target: ## Description
    @echo "Custom command"
```

---

## 🚦 CI/CD Status Checks

### GitHub Actions
- ✅ Main CI Pipeline - Comprehensive testing
- ✅ Security Scanning - Vulnerability checks  
- ✅ Per-Package CI - Individual package testing
- ✅ Dependency Updates - Automated maintenance

### Quality Gates
- 📊 **80% Code Coverage** - Minimum threshold
- 🔒 **Security Scanning** - No high vulnerabilities
- 🧪 **All Tests Pass** - Across Go 1.21-1.23
- 📝 **Linting Clean** - 50+ quality rules
- 🎨 **Format Check** - Consistent code style

---

## 🎯 Performance Tips

### Faster Local Development
```bash
make quick                  # Quick build and test
make test-quick            # Tests without coverage
cd packages/[name] && make test-short  # Package quick test
```

### Parallel Execution
```bash
make -j4 build             # Parallel builds
go test -parallel 4 ./...  # Parallel tests
```

### Caching
```bash
# Use Go build cache
export GOCACHE=/tmp/go-cache
go build -x ./...

# Clean caches when needed
make clean
```

---

## 🆘 Emergency Commands

### Fix Broken Build
```bash
# Reset everything
make clean
make deps-download
make build

# Per-package reset
cd packages/[package]
make clean deps build
```

### Fix Dependencies
```bash
# Root level
make deps-tidy deps-verify

# Per-package
cd packages/[package]/src
go mod tidy && go mod verify
```

### Fix Coverage Reports
```bash
go install github.com/wadey/gocovmerge@latest
make coverage-merge
```

### Reset Development Environment
```bash
make dev-reset
make dev-setup
```

---

## 📞 Need Help?

### Documentation
- 📖 [Full Documentation](./CICD_DOCUMENTATION.md)
- 🏗️ [Development Rules](./DEVELOPMENT_RULES.md)
- 📋 [Application Blueprint](./APPLICATION_PACKAGE_BLUEPRINT.md)

### Commands for Help
```bash
make help                   # Root-level help
cd packages/[name] && make help  # Package-level help
./setup-cicd.sh --help     # Setup script help
```

### Common File Locations
```bash
ls .github/workflows/       # GitHub Actions
ls packages/*/Makefile     # Package automation
find . -name ".golangci.yml"  # Linting configs
```

---

*Keep this guide handy for quick reference during development!*

**Version**: 1.0.0  
**Last Updated**: 2025-09-13