# 🚀 GoVel Framework - CI/CD Setup Guide

## Welcome to GoVel Development! 👋

This guide will help you set up the complete CI/CD infrastructure for the GoVel Framework, whether you're a new developer joining the project or setting up a fresh environment.

---

## 📋 Prerequisites

### Required Software

Before starting, ensure you have the following installed:

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Git** - Version control
- **Make** - Build automation (usually pre-installed on macOS/Linux)
- **GitHub CLI** (optional) - For easier GitHub integration

### Verify Prerequisites
```bash
# Check Go version
go version                 # Should show 1.21 or higher

# Check Git
git --version

# Check Make
make --version

# Check GitHub CLI (optional)
gh --version
```

### Environment Setup
```bash
# Set Go environment (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:$(go env GOPATH)/bin
export GOCACHE="$HOME/.cache/go-build"
export GOMODCACHE="$HOME/go/pkg/mod"
```

---

## 🏁 Quick Start (5 Minutes)

### 1. Clone and Navigate
```bash
# Clone the repository
git clone https://github.com/your-org/govel.git
cd govel

# Verify project structure
ls -la
```

### 2. Run Setup Scripts
```bash
# Setup main CI/CD infrastructure
./setup-cicd.sh

# Setup per-package CI/CD
./setup-per-package-cicd.sh

# Setup development environment
make dev-setup
```

### 3. Verify Setup
```bash
# Test the setup
make help                  # Should show all available commands
make info                  # Should show project information
make build                 # Should build all packages successfully
```

### 4. Run First Test
```bash
# Run complete CI pipeline
make all

# This will:
# ✅ Build all packages
# ✅ Run all tests
# ✅ Perform linting
# ✅ Run security scans
# ✅ Generate coverage reports
```

**🎉 If all commands complete successfully, you're ready to develop!**

---

## 📊 Detailed Setup Steps

### Step 1: Project Structure Understanding

After cloning, you'll see this structure:
```
govel/
├── .github/workflows/     # GitHub Actions CI/CD
├── packages/              # Individual Go packages
├── docs/                  # Documentation
├── scripts/               # Automation scripts
├── Makefile               # Root build automation
├── .golangci.yml          # Linting configuration
├── .codecov.yml           # Coverage configuration
└── setup-*.sh             # Setup scripts
```

### Step 2: Main CI/CD Setup

Run the main setup script:
```bash
./setup-cicd.sh --verbose
```

This creates:
- **GitHub Actions workflows** in `.github/workflows/`
- **Linting configuration** (`.golangci.yml`)
- **Coverage configuration** (`.codecov.yml`)
- **Build automation** (`Makefile`)
- **Dependency automation** (`.github/dependabot.yml`)

### Step 3: Per-Package CI/CD Setup

Run the per-package setup script:
```bash
./setup-per-package-cicd.sh --verbose
```

This creates for each package:
- **Individual CI workflow** (`packages/[name]/.github/workflows/ci.yml`)
- **Package linting config** (`packages/[name]/.golangci.yml`)
- **Package Makefile** (`packages/[name]/Makefile`)

### Step 4: Development Tools Installation

Install all development tools:
```bash
make dev-setup
```

This installs:
- **golangci-lint** - Code linting
- **gosec** - Security scanning
- **govulncheck** - Vulnerability checking
- **gocovmerge** - Coverage report merging

### Step 5: Verification and Testing

Verify everything works:
```bash
# Project information
make info

# Build all packages
make build

# Run all tests
make test

# Run linting
make lint

# Run security scans
make security

# Generate coverage report
make coverage-report

# Run complete pipeline
make all
```

---

## 🔧 Configuration Options

### Customizing Coverage Thresholds

Edit `.codecov.yml`:
```yaml
coverage:
  status:
    project:
      default:
        target: 80%        # Change this value
        threshold: 1%
```

### Customizing Linting Rules

Edit `.golangci.yml`:
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    # Add more linters here

linters-settings:
  gocyclo:
    min-complexity: 15     # Adjust complexity threshold
```

### Adding Custom Make Targets

Edit `Makefile`:
```makefile
my-custom-task: ## 📋 My custom task description
	@echo "Running my custom task"
	# Add your commands here
```

---

## 🏗️ Working with Packages

### Understanding Package Structure

Each package follows this structure:
```
packages/[package-name]/
├── .github/workflows/ci.yml    # Package CI pipeline
├── .golangci.yml              # Package linting rules
├── Makefile                   # Package build automation
└── src/                       # Package source code
    ├── go.mod                 # Go module definition
    ├── go.sum                 # Dependencies
    └── ...                    # Package implementation
```

### Working on an Existing Package

```bash
# Navigate to package
cd packages/encryption

# See available commands
make help

# Build the package
make build

# Run package tests
make test

# Run package linting
make lint

# Run complete package pipeline
make all
```

### Creating a New Package

```bash
# 1. Create package structure
mkdir -p packages/my-new-package/src
cd packages/my-new-package/src

# 2. Initialize Go module
go mod init govel/my-new-package

# 3. Create basic Go file
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello from my-new-package!")
}
EOF

# 4. Return to root and setup CI/CD
cd ../../..
./setup-per-package-cicd.sh

# 5. Verify new package setup
cd packages/my-new-package
make help
make build
make test
```

---

## 🚀 GitHub Integration

### Required Repository Settings

#### 1. GitHub Actions Permissions
1. Go to **Settings → Actions → General**
2. Set **Actions permissions** to "Allow all actions and reusable workflows"
3. Set **Workflow permissions** to "Read and write permissions"
4. Check "Allow GitHub Actions to create and approve pull requests"

#### 2. Required Secrets

Add these secrets in **Settings → Secrets and variables → Actions**:

| Secret Name | Description | Required |
|-------------|-------------|----------|
| `CODECOV_TOKEN` | Coverage reporting integration | Yes |

To get Codecov token:
1. Sign up at [codecov.io](https://codecov.io)
2. Connect your GitHub repository
3. Copy the repository token
4. Add it as `CODECOV_TOKEN` secret

#### 3. Branch Protection Rules (Recommended)

For `main` branch in **Settings → Branches**:
- ✅ Require pull request reviews
- ✅ Require status checks to pass
- ✅ Require branches to be up to date
- ✅ Include administrators

---

## 🧪 Testing Your Setup

### Local Testing Workflow

```bash
# 1. Create a test branch
git checkout -b test/ci-setup

# 2. Make a small change
echo "# Test Change" >> README.md

# 3. Run local CI pipeline
make all

# 4. If successful, commit and push
git add .
git commit -m "test: verify CI/CD setup"
git push origin test/ci-setup

# 5. Create PR and watch GitHub Actions run
```

### What to Expect in GitHub Actions

When you create a PR, you should see:
- ✅ **Main CI Pipeline** - Tests all packages
- ✅ **Security Scanning** - Vulnerability checks
- ✅ **Per-Package CI** - Individual package testing (if package files changed)

### Monitoring CI Results

- **Actions tab**: View running and completed workflows
- **PR checks**: See status checks at bottom of PR
- **Security tab**: View security scan results
- **Codecov comments**: Coverage reports on PRs

---

## 🔍 Troubleshooting Setup

### Common Issues and Solutions

#### ❌ Setup Script Permissions
```bash
# Problem: Permission denied when running setup scripts
# Solution: Make scripts executable
chmod +x setup-cicd.sh setup-per-package-cicd.sh
```

#### ❌ Go Module Issues
```bash
# Problem: Module download failures
# Solution: Clear module cache and retry
go clean -modcache
make deps-download
```

#### ❌ Development Tools Installation
```bash
# Problem: Tools not installing with make dev-setup
# Solution: Install manually
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/wadey/gocovmerge@latest
```

#### ❌ Build Failures
```bash
# Problem: Build failures due to missing dependencies
# Solution: Install package dependencies
cd packages/[failing-package]/src
go mod tidy
go mod download

# Or use Make targets
cd packages/[failing-package]
make deps
```

#### ❌ GitHub Actions Not Running
```bash
# Problem: Workflows not triggering
# Solution: Check repository settings
# 1. Ensure Actions are enabled
# 2. Check workflow file syntax
# 3. Verify branch protection rules
```

### Getting Help

If you encounter issues:

1. **Check the logs**:
   ```bash
   make test 2>&1 | tee build.log
   ```

2. **Verify environment**:
   ```bash
   make info
   go env
   ```

3. **Reset and retry**:
   ```bash
   make clean
   make dev-setup
   make all
   ```

4. **Check documentation**:
   - [Full Documentation](./CICD_DOCUMENTATION.md)
   - [Quick Reference](./CICD_QUICK_REFERENCE.md)

---

## 📚 Next Steps

### After Setup

1. **Read the documentation**:
   - [CI/CD Documentation](./CICD_DOCUMENTATION.md)
   - [Development Rules](./DEVELOPMENT_RULES.md)
   - [Quick Reference Guide](./CICD_QUICK_REFERENCE.md)

2. **Explore the codebase**:
   ```bash
   make info                    # Project overview
   ls packages/                 # Available packages
   cd packages/encryption && make help  # Package commands
   ```

3. **Make your first contribution**:
   - Find an issue to work on
   - Create a feature branch
   - Make changes and test locally
   - Create a pull request

### Recommended Learning Path

1. **Week 1**: Get familiar with the build system
   ```bash
   make help                    # Learn available commands
   make build test              # Basic build and test
   cd packages/types && make all  # Package-level development
   ```

2. **Week 2**: Understand CI/CD workflows
   - Create test PRs
   - Watch GitHub Actions runs
   - Experiment with different change types

3. **Week 3**: Advanced usage
   - Custom linting rules
   - Coverage optimization
   - Performance improvements

---

## ✅ Verification Checklist

Before starting development, ensure:

- [ ] **Repository cloned** and accessible
- [ ] **Setup scripts executed** without errors
- [ ] **Development tools installed** (`make dev-setup` successful)
- [ ] **Project builds** (`make build` successful)
- [ ] **Tests pass** (`make test` successful)
- [ ] **Linting passes** (`make lint` successful)
- [ ] **Security scans clean** (`make security` successful)
- [ ] **GitHub Actions enabled** in repository settings
- [ ] **Codecov token added** to repository secrets
- [ ] **First PR created** and CI passes

---

## 🎯 Summary

You now have a complete CI/CD infrastructure set up for the GoVel Framework:

### ✅ What You've Accomplished

- **🏗️ Complete CI/CD Pipeline**: Root-level and per-package automation
- **🔒 Security Integration**: Automated vulnerability scanning
- **📊 Coverage Tracking**: Comprehensive test coverage monitoring  
- **🧪 Multi-Platform Testing**: Go 1.21-1.23 across Linux/macOS/Windows
- **🛠️ Development Tools**: Linting, formatting, and quality checks
- **⚡ Build Automation**: 25+ Make targets for common tasks
- **🤖 Dependency Management**: Automated updates via Dependabot

### 🚀 Ready for Development

You can now:
- **Create feature branches** and develop with confidence
- **Run quality checks** locally before committing
- **Monitor CI/CD results** in GitHub Actions
- **Contribute to any package** in the ecosystem
- **Add new packages** with automated CI/CD setup

---

**Welcome to the GoVel Framework development team! 🎉**

*For questions or issues, refer to the [full documentation](./CICD_DOCUMENTATION.md) or create an issue in the repository.*

**Version**: 1.0.0  
**Last Updated**: 2025-09-13  
**Setup Time**: ~5-10 minutes