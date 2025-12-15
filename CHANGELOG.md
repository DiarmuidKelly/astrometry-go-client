# Changelog

## [1.3.1] - 2025-12-15

### Changes

- Release created from PR merge


## [1.3.0] - 2025-12-15

### Changes

- Release created from PR merge

## [1.2.0] - 2025-12-14

### Changes

- Release created from PR merge

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2025-12-15

### Changed

- **Reorganized package structure with clean separation of concerns**
  - Introduced `client` package as the main public API
  - Moved core solving implementation to `internal/solver/` (internal use only)
  - Created type re-exports in root `solver.go` for backwards compatibility
  - Separated public configuration into `config.go`
  - Public import paths: `github.com/DiarmuidKelly/astrometry-go-client/client` (main API), `github.com/DiarmuidKelly/astrometry-go-client/fov` (utilities)

- **Introduced unified `client.Client` wrapper**
  - Wraps internal solver implementation with clean public interface
  - Provides foundation for future expansion (e.g., annotation support, additional astrometry tools)
  - Methods: `NewClient()`, `Solve()`, `SolveBytes()`

- **Enhanced package architecture documentation**
  - Added "Current Implementation Status" section documenting completed refactor
  - Clearly separates current state from planned features (desired state)
  - Documents import paths, implemented features, and future additions

### Technical Details

- Internal package prevents direct access to implementation details
- Type re-exports maintain simple API surface while enabling future extensibility
- Structure prepares codebase for additional astrometry.net tools (image2xy, wcs utilities, etc.)

## [1.2.1] - 2025-12-15

### Fixed

- **WCS parser now correctly calculates field center coordinates** - Previously read reference pixel coordinates (CRVAL) instead of transforming to image center using CD matrix
- **Rotation angle calculation** - Added proper CD matrix-based rotation calculation (was returning 0°)
- **Pixel scale calculation** - Now uses full CD matrix instead of single element

### Changed

- **Reduced index file requirements** - Integration tests now download only index-4110.fits (24 MB) instead of 4 files (340 MB)
- **Updated test documentation** - Corrected Docker image references and solve parameters

### Technical Details

- Implemented WCS coordinate transformation: field_center = CRVAL + CD_matrix × (image_center - CRPIX)
- Rotation formula: `180° - atan2(CD1_2, CD1_1) × 180/π`
- All integration tests pass within tolerance (RA/Dec: <6 arcsec, rotation: 0.03°, pixel scale: 0.01%)

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

[Unreleased]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v1.2.1...HEAD
[1.2.1]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v0.3.1...v1.0.0
[0.3.1]: https://github.com/DiarmuidKelly/astrometry-go-client/compare/v0.1.0...v0.3.1
[0.1.0]: https://github.com/DiarmuidKelly/astrometry-go-client/releases/tag/v0.1.0
