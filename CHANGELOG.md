# Changelog

## [1.2.0] - 2025-12-14

### Changes

- Release created from PR merge


All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-12-14

### Changed
- **Migrated golangci-lint from v1.64.8 to v2.7.2**
  - Updated configuration to v2 format using official migration tool
  - Config structure changes: `linters-settings` → `linters.settings`
  - Exclusion rules: `issues.exclude-rules` → `linters.exclusions.rules`
  - Exclusion paths: `issues.exclude-dirs` → `linters.exclusions.paths`
  - GitHub Action updated from v4 to v7
- Improved linter configuration with better exclusions for legitimate complexity

### Fixed
- Removed deprecated config options (`run.timeout`, `unused.check-exported`)
- Updated timeout handling (now managed by GitHub Actions)

## [1.0.0] - 2025-12-14

### Added
- Docker setup documentation with execution mode comparison
- docker-compose.yml for development convenience
- .env.example for configuration template
- Wikipedia-verified camera sensor detection (185+ models: Canon, Nikon, Sony, Olympus, Panasonic)
- FOV package with automatic sensor detection from EXIF data
- Test coverage for FOV package (85.4%)

### Fixed
- Fixed .gitignore preventing cmd/astro-cli source from being committed
- Addressed linter warnings (built-in shadowing, error handling, constants)

## [0.1.0] - 2025-12-14

### Added
- Initial release of Astrometry Go Client
- Core solver package with plate-solving functionality
- Docker-based integration with dm90/astrometry container
- Support for configurable solve options (scale, downsample, depth, RA/Dec hints)
- WCS header parsing and result extraction
- CLI tool (astro-cli) for command-line plate solving
- Comprehensive test suite
- Example code for library usage
- Full documentation and contribution guidelines

### Features
- `NewClient()` - Create astrometry client with configurable options
- `Solve()` - Solve image files with full WCS output
- `SolveBytes()` - Solve images from byte arrays
- Structured error types for common failure cases
- Timeout support for long-running solves
- Output file collection (.wcs, .corr, .solved, etc.)

[Unreleased]: https://github.com/DiarmuidKelly/Astrometry-Go-Client/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/DiarmuidKelly/Astrometry-Go-Client/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/DiarmuidKelly/Astrometry-Go-Client/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/DiarmuidKelly/Astrometry-Go-Client/releases/tag/v0.1.0
