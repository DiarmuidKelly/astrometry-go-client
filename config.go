package client

import "time"

const (
	// DefaultDockerImage is the default Docker image used for plate-solving
	// Compatible images: "dm90/astrometry", "diarmuidk/astrometry-dockerised-solver", "ghcr.io/diarmuidkelly/astrometry-dockerised-solver"
	DefaultDockerImage = "diarmuidk/astrometry-dockerised-solver"
)

// ClientConfig holds configuration for the Astrometry client.
type ClientConfig struct {
	// DockerImage specifies the Docker image to use for solving.
	// Compatible with dm90/astrometry, diarmuidk/astrometry-dockerised-solver, or ghcr.io/diarmuidkelly/astrometry-dockerised-solver
	// Default: "diarmuidk/astrometry-dockerised-solver"
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

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		DockerImage: DefaultDockerImage,
		Timeout:     5 * time.Minute,
	}
}
