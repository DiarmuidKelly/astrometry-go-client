# astrometry-go-client: Scope & Architecture

## What This Package IS

**A lightweight, reusable Go library for offline astrometric plate-solving**

- **Core purpose:** Wrap the astrometry.net Docker solver with a clean Go API
- **Target users:** Go developers who need plate-solving in their projects
- **Deployment:** Imported as a Go module (`go get github.com/DiarmuidKelly/astrometry-go-client`)

## Package Structure (Desired State)

**Unified client pattern with internal implementation:**

```
astrometry-go-client/
├── client.go                # Unified Client, NewClient()
├── solver.go                # Type re-exports (SolveOptions, Result)
├── config.go                # Public ClientConfig
├── errors.go                # Public error types
├── annotate.go              # OPTIONAL: Annotation support (Annotate() method)
├── integration_test.go      # Integration tests with Docker
├── testdata/                # Test fixtures (images, WCS files)
├── internal/
│   └── solver/              # Core plate-solving implementation
│       ├── solver.go        # Solver logic
│       ├── result.go        # WCS parsing
│       ├── options.go       # SolveOptions
│       └── solver_test.go   # Unit tests
├── fov/                     # FOV & sensor utilities (public subpackage)
│   ├── fov.go               # FOV calculations
│   ├── image.go             # EXIF extraction
│   ├── indexes.go           # Index recommendations
│   ├── constants.go         # Sensor database
│   └── *_test.go            # Tests
├── cmd/
│   └── astro-cli/           # CLI tool (supported)
│       └── main.go
└── examples/
    └── basic/               # Usage examples
        └── main.go
```

**Import paths:**
- Main API: `github.com/DiarmuidKelly/astrometry-go-client/client`
- FOV utilities: `github.com/DiarmuidKelly/astrometry-go-client/fov`

## Current Implementation Status

**✅ Completed (as of recent refactor):**

The package structure has been reorganized with a clean separation of concerns:

```
astrometry-go-client/
├── client.go                # Public API: Client wrapper
├── solver.go                # Type re-exports (SolveOptions, Result)
├── config.go                # Public ClientConfig
├── errors.go                # Public error types
├── internal/solver/         # Core solving implementation (internal)
│   ├── solver.go            # Solver logic, Solve(), SolveBytes()
│   ├── result.go            # WCS parsing
│   └── options.go           # SolveOptions
├── fov/                     # FOV utilities (public, fully implemented)
│   ├── fov.go               # FOV calculations
│   ├── image.go             # EXIF extraction
│   ├── indexes.go           # Index recommendations
│   └── constants.go         # Sensor database
├── cmd/astro-cli/           # CLI tool
└── examples/basic/          # Usage examples
```

**Import paths (current):**
- Main API: `github.com/DiarmuidKelly/astrometry-go-client/client`
- FOV utilities: `github.com/DiarmuidKelly/astrometry-go-client/fov`

**Implemented features:**
- ✅ Unified `client.Client` wrapping internal solver
- ✅ `Solve()` and `SolveBytes()` methods
- ✅ Docker run/exec modes
- ✅ Full WCS parsing with Result struct
- ✅ FOV calculation from EXIF
- ✅ Index file recommendations
- ✅ Sensor database with common cameras
- ✅ CLI tool
- ✅ Integration tests with ground truth validation

**🚧 Future additions (planned):**
- ⏳ Optional annotation support (`Annotate()` method)
- ⏳ Additional astrometry.net tool wrappers (image2xy, wcs-xy2rd, wcs-rd2xy, etc.)
- ⏳ Additional utility methods (see "Future methods" comment in client.go)

## What This Package DOES (Desired State)

### Core Functionality

1. **Solve single images** via `client.Solve(ctx, imagePath, options)`
2. **Parse WCS results** into structured Go types
3. **Support both Docker modes** (run/exec)
4. **Detect camera sensors** from EXIF (FOV calculation helpers)
5. **Optional annotation** via `client.Annotate()` - exposes Docker image tools:
   - Constellation overlays
   - NGC/Messier objects
   - Named bright stars
   - RA/Dec grid overlays
6. **CLI tool** for quick testing and bash scripting
7. **Provide sensible defaults** for common use cases

### Design Principles

- **Minimal dependencies** (only `goexif` for sensor detection)
- **No external services** (100% offline)
- **Stateless API** (each Solve() call is independent)
- **Context-aware** (timeouts, cancellation)
- **Well-tested** (unit + integration tests)
- **Clear error types** for caller-controlled retry logic

## What This Package DOES NOT DO

### Out of Scope

❌ **Built-in batch processing** - Call `Solve()` in your own loop (see examples)
❌ **Custom image processing** - Only expose what Docker image provides
❌ **HTTP API server** - See separate Astrometry-API-Server repo
❌ **File watching/orchestration** - User's responsibility
❌ **Image format conversion** - Accept what astrometry.net accepts
❌ **Database/persistence** - Caller handles storage
❌ **Progress tracking** - Caller implements if needed
❌ **Retry logic** - Provide clear error types; caller decides retry strategy
❌ **Rate limiting** - Not needed (offline, local)
❌ **Authentication** - Not applicable

## API Surface (Keep Minimal)

### Public Types

```go
// Core client
type Client struct { ... }
type ClientConfig struct {
    IndexPath string        // Required
    DockerImage string      // Optional (default provided)
    Timeout time.Duration   // Optional (default 5min)
    UseDockerExec bool      // Use docker exec vs docker run
    ContainerName string    // For docker exec mode
}

// Solving
type SolveOptions struct {
    ScaleLow, ScaleHigh float64
    ScaleUnits string
    DownsampleFactor int
    RA, Dec, Radius float64  // Optional hints
    NoPlots bool
    Verbose bool
}

type Result struct {
    Solved bool
    RA, Dec float64
    PixelScale, Rotation float64
    FieldWidth, FieldHeight float64
    WCSHeader map[string]string
    OutputFiles []string    // Generated files (.wcs, .corr, etc)
    SolveTime float64
}

// Optional annotation support
type AnnotationOptions struct {
    ShowConstellations bool
    ShowNGC bool
    ShowBrightStars bool
    GridSpacing int         // RA/Dec grid every N arcmin (0 = no grid)
}

// Error types for caller decisions
var (
    ErrNoSolution   = errors.New("no solution found")        // Don't retry
    ErrTimeout      = errors.New("solve timed out")          // Maybe retry with longer timeout
    ErrDockerFailed = errors.New("docker command failed")    // Check deployment
    ErrInvalidInput = errors.New("invalid input")            // Don't retry
)

// Core methods
func NewClient(config *ClientConfig) (*Client, error)
func (c *Client) Solve(ctx context.Context, imagePath string, opts *SolveOptions) (*Result, error)
func (c *Client) SolveBytes(ctx context.Context, data []byte, format string, opts *SolveOptions) (*Result, error)

// Optional annotation (exposes Docker image tools)
func (c *Client) Annotate(ctx context.Context, result *Result, imagePath string, opts *AnnotationOptions) (outputPath string, err error)
```

### Helper Package (Optional Extras)

```go
// fov/ sub-package - Camera sensor utilities
func DetectSensorFromEXIF(imagePath string) (*SensorInfo, error)
func CalculateFOV(sensorWidth, focalLength float64) float64
func RecommendIndexes(fov float64) []string

type SensorInfo struct {
    Width, Height float64  // mm
    CameraModel string
}
```

## CLI Tool Position

**`cmd/astro-cli` is a supported tool, not just a reference implementation**

### Purpose

- ✅ **Quick testing** - Fast manual verification without writing Go code
- ✅ **Bash scripting** - Enables batch processing via shell scripts
- ✅ **Documentation by example** - Shows how to use the library
- ✅ **Dogfooding** - Forces us to use our own API

### Design Principles

- Keep core functions available via CLI flags
- Simple flag parsing → library call → JSON output
- No CLI-specific logic (all logic in library)
- Flags mirror library API where possible

### Usage Guidance

- ✅ **Use for:** Quick testing, bash scripts, manual solves
- ⚠️ **Acceptable for:** Simple automation, cron jobs
- ❌ **Don't use for:** Complex production pipelines (import the library instead)

### Batch Processing with CLI

Users can write bash scripts that call `astro-cli` in a loop:

```bash
for img in *.jpg; do
  astro-cli --image "$img" --index-path ~/astrometry-data > "${img%.jpg}.json"
done
```

See `examples/batch/` for patterns.

## Example Implementations

### In This Repo

- **CLI tool** (`cmd/astro-cli`) - Supported tool for testing and scripting
- **Example scripts** (`examples/batch/`) - Bash/Python patterns for batch processing
- **Library examples** (`examples/basic/`) - Go code showing API usage

### Separate Repos (User-Built)

- **Astrometry-API-Server** - HTTP API wrapper (separate repo, maintained by us)
- **Custom pipelines** - Users import library in their Go code
- **Web UIs** - Users build on top of API server

## Maintenance Philosophy

### What Gets Added

✅ Bug fixes in WCS parsing
✅ Support for new Docker image versions
✅ Better error messages
✅ Performance improvements in core solving
✅ Additional camera sensors in FOV database
✅ Exposing existing Docker annotation tools (plot-constellations, etc)
✅ CLI features that mirror library capabilities

### What Gets Rejected

❌ "Add built-in batch processing with progress bars"
❌ "Add web UI"
❌ "Add database support"
❌ "Add custom image processing (beyond Docker tools)"
❌ "Add automatic retry logic"
❌ "Add file watching"

**Response:** "That's an application concern. Call `Solve()` in your own loop or import the library. See `examples/` and Astrometry-API-Server for patterns."

## Success Criteria

This package succeeds when:

1. Other Go projects can `import` and use it easily
2. API is stable and rarely changes
3. No feature creep - stays focused on solving
4. Well-documented with clear examples
5. Minimal dependencies and maintenance burden

## File Boundaries

### Keep in Go-Client

- Core solving logic (`solver.go`, `result.go`, `options.go`)
- WCS parsing
- Docker integration (run/exec modes)
- Optional annotation via Docker tools (`annotate.go`)
- FOV/sensor helpers (`fov/` package)
- Error types (`errors.go`)
- CLI tool (`cmd/astro-cli/`)
- Tests for all above
- Example scripts (`examples/`)
- Library documentation

### Push to Users/Other Repos

- HTTP servers (see Astrometry-API-Server)
- Complex batch orchestration (users write their own)
- File watching/monitoring
- Database integration
- Custom retry strategies
- Application-specific logic
- Web UIs

## Batch Processing Support

**Stance:** Supported via examples, not built-in.

### Patterns

1. **Bash script with CLI:** Call `astro-cli` in a loop
2. **Go with goroutine pool:** Call `Solve()` concurrently
3. **Python wrapper:** Subprocess calls to CLI
4. **API server:** HTTP endpoints for remote batch jobs

All patterns documented in `examples/batch/`.

**Why no built-in batch method:**

- Overhead of multiple `Solve()` calls is minimal
- Users have different needs (parallelism, error handling, progress)
- Keeps library focused and flexible
- Examples cover common patterns
