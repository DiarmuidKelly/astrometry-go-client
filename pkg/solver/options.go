// Package solver provides a Go client for the Astrometry.net plate-solving service
// via the dm90/astrometry Docker container.
package solver

import "time"

const (
	// DefaultDockerImage is the default Docker image used for plate-solving
	DefaultDockerImage = "dm90/astrometry"
)

// ClientConfig holds configuration for the Astrometry client.
type ClientConfig struct {
	// DockerImage specifies the Docker image to use for solving.
	// Default: "dm90/astrometry"
	DockerImage string

	// IndexPath is the host path to the astrometry index files.
	// This directory will be mounted into the Docker container.
	// Required.
	IndexPath string

	// TempDir is the working directory for images and output files.
	// If empty, the system temporary directory will be used.
	TempDir string

	// Timeout is the maximum duration for the solve operation.
	// Default: 5 minutes
	Timeout time.Duration

	// UseDockerExec enables using docker exec on an existing container
	// instead of spawning new containers with docker run.
	// When true, ContainerName must be specified.
	// Default: false
	UseDockerExec bool

	// ContainerName is the name of the running container to exec commands in.
	// Only used when UseDockerExec is true.
	ContainerName string
}

// SolveOptions holds parameters for a plate-solving operation.
type SolveOptions struct {
	// ScaleLow is the lower bound of the image scale in the specified units.
	ScaleLow float64

	// ScaleHigh is the upper bound of the image scale in the specified units.
	ScaleHigh float64

	// ScaleUnits specifies the units for ScaleLow and ScaleHigh.
	// Valid values: "degwidth", "arcminwidth", "arcsecperpix"
	// Default: "arcminwidth"
	ScaleUnits string

	// DownsampleFactor reduces the image resolution by this factor.
	// Higher values speed up solving but reduce accuracy.
	// Default: 2
	DownsampleFactor int

	// DepthLow is the minimum number of quads to try.
	// DepthHigh is the maximum number of quads to try.
	// Default: [10, 20]
	DepthLow  int
	DepthHigh int

	// NoPlots disables generation of plot files (RedGreen, etc.).
	// Default: true (no plots)
	NoPlots bool

	// RA, Dec, and Radius provide a search hint for the solver.
	// RA and Dec are in degrees (J2000).
	// Radius is the search radius in degrees.
	// If RA is 0, no search hint is used.
	RA     float64
	Dec    float64
	Radius float64

	// OverwriteExisting allows overwriting existing output files.
	// Default: false
	OverwriteExisting bool

	// Verbose enables verbose output from solve-field.
	// Default: false
	Verbose bool

	// KeepTempFiles preserves temporary files for debugging.
	// When true, temp directory and all solve output files are not deleted.
	// Default: false
	KeepTempFiles bool
}

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		DockerImage: DefaultDockerImage,
		Timeout:     5 * time.Minute,
	}
}

// DefaultSolveOptions returns SolveOptions with sensible defaults.
func DefaultSolveOptions() *SolveOptions {
	return &SolveOptions{
		ScaleUnits:       "arcminwidth",
		DownsampleFactor: 2,
		DepthLow:         10,
		DepthHigh:        20,
		NoPlots:          true,
		Verbose:          false,
	}
}
