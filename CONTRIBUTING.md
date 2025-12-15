# Contributing Guide

Thank you for contributing! This project uses automated PR-based releases with semantic versioning.

## Quick Start

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/astrometry-go-client.git
cd astrometry-go-client

# 2. Create feature branch
git checkout -b feat/my-feature

# 3. Make changes with conventional commits
git commit -m "feat: add batch processing support"

# 4. Run tests
make test
make lint

# 5. Push and create PR
git push origin feat/my-feature
```

## PR Title Format

Your PR title determines the version bump:

| PR Title              | Version Change | Use Case         |
| --------------------- | -------------- | ---------------- |
| `[MAJOR] description` | 0.1.0 → 1.0.0  | Breaking changes |
| `[MINOR] description` | 0.1.0 → 0.2.0  | New features     |
| `[PATCH] description` | 0.1.0 → 0.1.1  | Bug fixes        |
| `[SKIP] description`  | No release     | Docs/chore only  |

**Or use conventional commits:**

- `feat: description` → MINOR
- `fix: description` → PATCH
- `docs: description` → SKIP
- `chore: description` → SKIP

## Commit Message Format

Use **Conventional Commits** for automatic changelog generation:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, etc)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Tests
- `chore`: Maintenance

### Examples

```bash
feat: add support for FITS image format
fix: correct WCS header parsing for rotation
docs: update API reference with timeout examples
feat!: change Solve() signature to accept context  # Breaking change
```

## What Happens When PR Merges?

1. **PR title validated** automatically
2. **Labels auto-applied** based on type
3. **On merge:**
   - ✅ Changelog generated from commits
   - ✅ Version bumped in VERSION file
   - ✅ Git tag created
   - ✅ GitHub release published with binaries
   - ✅ Comment added to PR with release link

## Testing Checklist

Before submitting PR:

- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Added tests for new functionality
- [ ] Updated documentation (if applicable)
- [ ] Commits follow conventional format
- [ ] PR title follows required format
- [ ] Code builds successfully (`make build`)

## Code Style

- Follow standard Go conventions
- Run `gofmt` on all code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Write table-driven tests where appropriate

## PR Review Process

Your PR will be merged if:

- ✅ Code follows Go best practices
- ✅ Commits follow conventional format
- ✅ All tests pass
- ✅ Documentation is updated
- ✅ PR title is valid format
- ✅ Code owner approves (@DiarmuidKelly)

## Version Bumping Rules

- **PATCH** (0.0.X): Bug fixes, docs, minor tweaks
- **MINOR** (0.X.0): New features, backward-compatible
- **MAJOR** (X.0.0): Breaking changes, incompatible API

**Note:** You don't need to manually update version numbers. When your PR is merged, the CI/CD workflow automatically updates the `VERSION` file based on your PR title.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker (for integration tests)
- golangci-lint (for linting)

### Install Dependencies

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Download go dependencies
go mod download
```

### Running Tests

```bash
# Unit tests
make test

# Run with coverage
make test-coverage

# Integration tests (requires Docker)
make test-integration

# All tests
make test-all
```

### Running Linter

```bash
# Run linter
make lint

# Auto-fix linting issues where possible
make lint-fix
```

### Building

```bash
# Build all binaries
make build

# Build CLI only
make build-cli

# Install to GOPATH
make install
```

## Adding New Features

1. **Create an issue** first to discuss the feature
2. **Write tests** before implementing (TDD approach)
3. **Update documentation** (README.md, godoc comments)
4. **Add examples** if introducing new API surface
5. **Update CHANGELOG.md** header with unreleased changes (optional - CI will do this)

## Reporting Bugs

When reporting bugs, include:

- Go version (`go version`)
- Operating system
- Docker version
- Minimal code to reproduce
- Expected vs actual behavior
- Error messages and logs

## Questions?

- Check [README.md](README.md) for project overview
- See [.github/BRANCH_PROTECTION.md](.github/BRANCH_PROTECTION.md) for branch rules
- Open an [issue](https://github.com/DiarmuidKelly/astrometry-go-client/issues) for questions

Thank you for contributing!
