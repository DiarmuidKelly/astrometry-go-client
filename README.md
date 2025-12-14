# Astrometry Go Client

![Version](https://img.shields.io/github/v/release/DiarmuidKelly/Astrometry-Go-Client?label=version)
![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/DiarmuidKelly/Astrometry-Go-Client)
[![Go Report Card](https://goreportcard.com/badge/github.com/DiarmuidKelly/Astrometry-Go-Client)](https://goreportcard.com/report/github.com/DiarmuidKelly/Astrometry-Go-Client)

**Offline astrometric plate-solving for Go** - Solve astronomical images locally without internet access using the [dm90/astrometry](https://hub.docker.com/r/dm90/astrometry) Docker container. Complete privacy and control over your data with no dependency on external services.

## Features

- **100% Offline Operation** - No internet required after initial setup
- **Complete Privacy** - Your images never leave your machine
- **No Rate Limits** - Solve unlimited images without API quotas
- **Fast Local Processing** - No network latency or upload time
- Simple, type-safe API for plate-solving astronomical images
- Docker-based integration (no complex installation required)
- Support for scale hints, downsampling, and RA/Dec position hints
- WCS header parsing with structured results
- Context-aware with timeout support
- Comprehensive error handling
- CLI tool for quick command-line solving
- Full test coverage

## Why Offline?

Unlike cloud-based plate-solving services, this library runs entirely on your local machine:

- **Privacy First** - Your astrophotography images and location data stay private
- **No Internet Required** - Work in remote locations, observatories, or anywhere offline
- **Free & Unlimited** - No API keys, rate limits, or subscription fees
- **Faster Processing** - No upload/download time, just local computation
- **Full Control** - Customize solve parameters without service restrictions
- **Self-Hosted** - Perfect for automated pipelines, observatories, and air-gapped systems

## Prerequisites

**Note:** Internet is only required once during initial setup to download the Docker image and index files. After setup, the solver runs 100% offline.

### Required
- **Go 1.21+**
- **Docker** with the `dm90/astrometry` image pulled (one-time download)
- **Astrometry.net index files** downloaded to a local directory (one-time download)

### Quick Setup

```bash
# Clone the repository (for development)
git clone https://github.com/DiarmuidKelly/Astrometry-Go-Client.git
cd Astrometry-Go-Client

# Download all index files (~350MB)
./scripts/download-indexes.sh
```

### Index Files Guide

Astrometry.net requires index files that match your camera's field of view (FOV). **Using an index that's too wide for your FOV guarantees failure** - the star patterns physically won't fit in your image.

#### Choosing the Right Index Files

First, calculate your FOV:
```
FOV (degrees) = (sensor_width_mm / focal_length_mm) × 57.3
```

Then download 2-3 indexes that bracket your FOV:

| Index File | Field Width | Size | Use Case | Download |
|------------|-------------|------|----------|----------|
| index-4119 | 0.1° - 0.2° | 144 KB | Planetary imaging, very long focal length | [Download](http://data.astrometry.net/4100/index-4119.fits) |
| index-4118 | 0.2° - 0.28° | 187 KB | Long focal length telescopes | [Download](http://data.astrometry.net/4100/index-4118.fits) |
| index-4117 | 0.28° - 0.4° | 248 KB | | [Download](http://data.astrometry.net/4100/index-4117.fits) |
| index-4116 | 0.4° - 0.56° | 409 KB | | [Download](http://data.astrometry.net/4100/index-4116.fits) |
| index-4115 | 0.56° - 0.8° | 740 KB | Medium-long focal length | [Download](http://data.astrometry.net/4100/index-4115.fits) |
| index-4114 | 0.8° - 1.1° | 1.4 MB | | [Download](http://data.astrometry.net/4100/index-4114.fits) |
| index-4113 | 1.1° - 1.6° | 2.7 MB | | [Download](http://data.astrometry.net/4100/index-4113.fits) |
| index-4112 | 1.6° - 2.2° | 5.3 MB | **DSLR + telephoto (common)** | [Download](http://data.astrometry.net/4100/index-4112.fits) |
| index-4111 | 2.2° - 3.0° | 10 MB | **DSLR + normal lens (common)** | [Download](http://data.astrometry.net/4100/index-4111.fits) |
| index-4110 | 3.0° - 4.2° | 25 MB | **Wide angle DSLR (common)** | [Download](http://data.astrometry.net/4100/index-4110.fits) |
| index-4109 | 4.2° - 5.6° | 50 MB | Very wide angle | [Download](http://data.astrometry.net/4100/index-4109.fits) |
| index-4108 | 5.6° - 8.0° | 95 MB | Ultra-wide, fisheye | [Download](http://data.astrometry.net/4100/index-4108.fits) |
| index-4107 | 8.0° - 11.0° | 165 MB | All-sky cameras | [Download](http://data.astrometry.net/4100/index-4107.fits) |

**Total size of all indexes:** ~350 MB

#### Example Setups

**DSLR Astrophotography (APS-C sensor, 24mm wide):**
```bash
mkdir -p ~/astrometry-data && cd ~/astrometry-data
# 200mm lens: ~7° FOV
wget http://data.astrometry.net/4100/index-4108.fits
wget http://data.astrometry.net/4100/index-4109.fits

# 50mm lens: ~27° FOV
wget http://data.astrometry.net/4100/index-4110.fits
wget http://data.astrometry.net/4100/index-4111.fits
```

**Telescope (1000mm focal length, APS-C sensor):**
```bash
mkdir -p ~/astrometry-data && cd ~/astrometry-data
# ~1.4° FOV
wget http://data.astrometry.net/4100/index-4113.fits
wget http://data.astrometry.net/4100/index-4114.fits
```

**Quick Solver (50mm-300mm lenses - recommended default):**
```bash
mkdir -p ~/astrometry-data && cd ~/astrometry-data
# Download 1.1° - 4.2° range (~43 MB total)
# Covers DSLR + 50-300mm focal lengths
wget http://data.astrometry.net/4100/index-4110.fits  # 3.0° - 4.2° (50-70mm)
wget http://data.astrometry.net/4100/index-4111.fits  # 2.2° - 3.0° (70-110mm)
wget http://data.astrometry.net/4100/index-4112.fits  # 1.6° - 2.2° (110-150mm)
wget http://data.astrometry.net/4100/index-4113.fits  # 1.1° - 1.6° (150-220mm)
```

**Note:** The `scripts/download-indexes.sh` script downloads all indexes automatically.

#### Why Index Matching Matters

Index files contain pre-computed star patterns (called "quads") at specific angular scales. **Using the wrong index causes failure, not just slowness.**

**What happens with mismatched indexes:**

| Scenario | Result | Why |
|----------|--------|-----|
| **Narrow image (0.5°) + Wide index (6-8°)** | ❌ **Guaranteed failure** | Star patterns in the index span 6-8°. Your 0.5° image physically cannot contain patterns that large - they don't fit in the frame. |
| **Wide image (8°) + Narrow index (0.5°)** | ⚠️ **Likely failure** | Your wide image does contain small sub-regions with narrow patterns, but matching is unreliable, very slow, and usually fails. |

**The critical rule:**
- **Index scale larger than your FOV → guaranteed failure** (patterns don't fit in your image)
- **Index scale smaller than your FOV → likely failure** (unreliable matching of sub-regions)

✅ **Best practice:** Download 2-3 indexes that **bracket** your expected FOV (one above, one below). This ensures reliable, fast solving regardless of minor FOV variations.

## Docker Setup

This library requires the `dm90/astrometry` Docker container to perform plate-solving. You have several options for running the dependency:

### Docker Execution Modes

The library supports two Docker execution modes:

#### 1. Docker Run Mode (Default)

**How it works**: The library spawns a new Docker container for each solve operation, then removes it.

**Pros**:
- Simple setup - no container management required
- Clean isolation per solve

**Cons**:
- Slower for multiple solves (~1-2s container startup overhead per solve)
- More Docker overhead

**Setup**: Pull the image, then use the library - that's it!

```bash
docker pull dm90/astrometry:latest
```

**Client configuration**:
```go
config := &solver.ClientConfig{
    IndexPath: "/path/to/astrometry-data",  // Path to your index files
}
client, err := solver.NewClient(config)
```

#### 2. Docker Exec Mode (Recommended for Development)

**How it works**: Uses a long-running Docker container and executes solve commands via `docker exec`.

**Pros**:
- **Faster**: No container startup overhead (typically 1-2s faster per solve)
- Ideal for development, testing, or batch processing

**Cons**:
- Requires manual container management

**Setup Option A: Using Docker Compose (Recommended)**

```bash
# 1. Copy and configure environment
cp .env.example .env
# Edit .env and set ASTROMETRY_INDEX_PATH to your index files directory

# 2. Start the container
docker compose up -d

# 3. Verify it's running
docker ps | grep astrometry-solver
```

**Setup Option B: Manual Docker Run**

```bash
docker run -d \
  --name astrometry-solver \
  -v ~/astrometry-data:/usr/local/astrometry/data:ro \
  -v /tmp/astrometry-shared:/shared-data \
  dm90/astrometry:latest \
  tail -f /dev/null
```

**Client configuration**:
```go
config := &solver.ClientConfig{
    IndexPath:     "/path/to/astrometry-data",
    UseDockerExec: true,
    ContainerName: "astrometry-solver",
}
client, err := solver.NewClient(config)
```

### Performance Comparison

| Operation | Docker Run Mode | Docker Exec Mode | Difference |
|-----------|-----------------|------------------|------------|
| First solve | ~15s | ~13s | ~2s faster |
| Subsequent solves (each) | ~15s | ~13s | ~2s faster |
| Setup overhead | None | Start container once | One-time |

**Recommendation**: Use **docker exec mode** for development/testing. Either mode works well for production depending on your orchestration setup.

### Full Stack Setup

For a complete REST API server with web interface, see the [Astrometry API Server](https://github.com/DiarmuidKelly/Astrometry-API-Server) project. It includes:
- docker-compose orchestration for both the API server and solver
- RESTful HTTP API
- Swagger documentation
- Health monitoring

## Installation

### As a Library

```bash
go get github.com/DiarmuidKelly/Astrometry-Go-Client
```

### CLI Tool

```bash
go install github.com/DiarmuidKelly/Astrometry-Go-Client/cmd/astro-cli@latest
```

Or build from source:

```bash
git clone https://github.com/DiarmuidKelly/Astrometry-Go-Client.git
cd Astrometry-Go-Client
make install
```

## Quick Start

### Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/DiarmuidKelly/Astrometry-Go-Client/pkg/solver"
)

func main() {
    // Create client
    config := &solver.ClientConfig{
        IndexPath: "/path/to/astrometry-data",
    }
    client, err := solver.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // Configure solve options
    opts := solver.DefaultSolveOptions()
    opts.ScaleLow = 1.0   // 1 arcmin/width
    opts.ScaleHigh = 3.0  // 3 arcmin/width

    // Solve the image
    result, err := client.Solve(context.Background(), "image.jpg", opts)
    if err != nil {
        log.Fatal(err)
    }

    if result.Solved {
        fmt.Printf("RA: %.6f, Dec: %.6f\n", result.RA, result.Dec)
        fmt.Printf("Pixel Scale: %.2f arcsec/pixel\n", result.PixelScale)
    } else {
        fmt.Println("Image could not be solved")
    }
}
```

### CLI Usage

```bash
astro-cli \
  --image photo.jpg \
  --index-path ~/astrometry-data \
  --scale-low 1 \
  --scale-high 3 \
  --downsample 2
```

Output (JSON):

```json
{
  "solved": true,
  "ra": 120.123456,
  "dec": 45.987654,
  "pixel_scale": 1.23,
  "rotation": 15.5,
  "field_width": 2.5,
  "field_height": 1.8,
  "solve_time": 12.34
}
```

## API Reference

### Client Configuration

```go
type ClientConfig struct {
    DockerImage   string        // Default: "dm90/astrometry"
    IndexPath     string        // Required: path to index files
    TempDir       string        // Optional: temp directory for processing
    Timeout       time.Duration // Default: 5 minutes
    UseDockerExec bool          // Use docker exec mode (default: false)
    ContainerName string        // Container name for docker exec mode
}
```

### Solve Options

```go
type SolveOptions struct {
    ScaleLow         float64  // Lower bound of image scale
    ScaleHigh        float64  // Upper bound of image scale
    ScaleUnits       string   // "degwidth", "arcminwidth", "arcsecperpix"
    DownsampleFactor int      // Reduce resolution (default: 2)
    DepthLow         int      // Min quads to try (default: 10)
    DepthHigh        int      // Max quads to try (default: 20)
    NoPlots          bool     // Disable plot generation (default: true)
    RA               float64  // RA hint in degrees (optional)
    Dec              float64  // Dec hint in degrees (optional)
    Radius           float64  // Search radius in degrees (optional)
    Verbose          bool     // Enable verbose output
}
```

### Result Structure

```go
type Result struct {
    Solved      bool              // Whether the image was solved
    RA          float64           // Right ascension (J2000, degrees)
    Dec         float64           // Declination (J2000, degrees)
    PixelScale  float64           // arcsec/pixel
    Rotation    float64           // Field rotation (degrees)
    FieldWidth  float64           // Field of view width (degrees)
    FieldHeight float64           // Field of view height (degrees)
    WCSHeader   map[string]string // Raw WCS header fields
    OutputFiles []string          // Paths to generated files
    SolveTime   float64           // Solve duration (seconds)
}
```

### Methods

**`NewClient(config *ClientConfig) (*Client, error)`**

Creates a new astrometry client with the given configuration.

**`Solve(ctx context.Context, imagePath string, opts *SolveOptions) (*Result, error)`**

Solves a single image file and returns the plate solution.

**`SolveBytes(ctx context.Context, data []byte, format string, opts *SolveOptions) (*Result, error)`**

Solves image data from a byte slice (useful for in-memory images).

## Examples

See the [examples/](examples/) directory for more usage examples:

- `examples/basic/` - Basic plate-solving example
- `examples/batch/` - Batch processing multiple images (coming soon)
- `examples/with-hints/` - Using RA/Dec hints for faster solving (coming soon)

## Error Handling

The library provides structured error types:

```go
var (
    ErrNoSolution    = errors.New("no solution found")
    ErrTimeout       = errors.New("solve operation timed out")
    ErrDockerFailed  = errors.New("docker command failed")
    ErrInvalidInput  = errors.New("invalid input parameters")
    ErrWCSParseFailed = errors.New("failed to parse WCS output")
)
```

Example:

```go
result, err := client.Solve(ctx, imagePath, opts)
if err != nil {
    if errors.Is(err, solver.ErrTimeout) {
        log.Println("Solve timed out - try increasing timeout or downsample")
    } else if errors.Is(err, solver.ErrDockerFailed) {
        log.Println("Docker error - check Docker is running")
    }
    return err
}

if !result.Solved {
    log.Println("No solution - try adjusting scale parameters")
}
```

## Performance Tips

1. **Use scale bounds**: Providing `ScaleLow` and `ScaleHigh` dramatically speeds up solving
2. **Downsample**: Higher downsample factors (2-4) work well for most images
3. **RA/Dec hints**: If you know approximate coordinates, use them to reduce search space
4. **Index files**: Download only the indexes appropriate for your field of view

## Development

### Building

```bash
make build        # Build all binaries
make test         # Run tests
make lint         # Run linter
make clean        # Clean build artifacts
```

### Running Tests

```bash
go test ./...
```

Integration tests requiring Docker can be run with:

```bash
go test -tags=integration ./...
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for workflow details.

**Quick Start:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes with conventional commits
4. Push and create a PR with `[MAJOR]`, `[MINOR]`, or `[PATCH]` prefix

Auto-release workflow handles versioning, changelog, and releases on PR merge.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Links

- [Astrometry API Server](https://github.com/DiarmuidKelly/Astrometry-API-Server) - REST API server for this library
- [Changelog](CHANGELOG.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Issues](https://github.com/DiarmuidKelly/Astrometry-Go-Client/issues)
- [Astrometry.net](http://astrometry.net/)
- [dm90/astrometry Docker Image](https://hub.docker.com/r/dm90/astrometry)

## Acknowledgments

- [Astrometry.net](http://astrometry.net/) project for the plate-solving engine
- [dm90/astrometry](https://hub.docker.com/r/dm90/astrometry) for the containerized version
