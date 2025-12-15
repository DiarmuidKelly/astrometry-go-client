// Package client provides a unified Go client for astrometry.net operations.
//
// This package wraps multiple astrometry.net tools (solve-field, image2xy, wcs utilities, etc.)
// in a convenient Go API, running them via Docker containers.
//
// # Basic Usage
//
//	config := &client.ClientConfig{
//		IndexPath: "/path/to/index/files",
//	}
//
//	c, err := client.NewClient(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	opts := &solver.SolveOptions{
//		ScaleLow: 1.0,
//		ScaleHigh: 3.0,
//		ScaleUnits: "arcminwidth",
//	}
//
//	result, err := c.Solve(context.Background(), "image.jpg", opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if result.Solved {
//		fmt.Printf("RA: %.6f, Dec: %.6f\n", result.RA, result.Dec)
//	}
//
// # Docker Execution Modes
//
// The client supports two Docker execution modes:
//
// 1. Docker Run Mode (default): Spawns a new container for each operation.
// Simpler but slower for multiple operations.
//
// 2. Docker Exec Mode: Uses an existing long-running container via docker exec.
// Faster for multiple operations, requires a running container.
//
//	config := &client.ClientConfig{
//		IndexPath:     "/path/to/indexes",
//		UseDockerExec: true,
//		ContainerName: "astrometry-solver",
//	}
package client

import (
	"context"
	"fmt"
	"os"

	"github.com/DiarmuidKelly/astrometry-go-client/internal/solver"
)

// Client is the unified interface for all astrometry.net operations.
type Client struct {
	config       *ClientConfig
	solverClient *solver.Client
}

// NewClient creates a new astrometry Client with the given configuration.
func NewClient(config *ClientConfig) (*Client, error) {
	if config == nil {
		config = DefaultClientConfig()
	}

	// Validate required fields
	if config.IndexPath == "" {
		return nil, fmt.Errorf("%w: IndexPath is required", ErrInvalidInput)
	}

	// Check that index path exists
	if _, err := os.Stat(config.IndexPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: IndexPath does not exist: %s", ErrInvalidInput, config.IndexPath)
	}

	// Set defaults
	if config.DockerImage == "" {
		config.DockerImage = DefaultDockerImage
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultClientConfig().Timeout
	}
	if config.TempDir == "" {
		config.TempDir = os.TempDir()
	}

	// Convert to internal solver config
	solverCfg := &solver.ClientConfig{
		DockerImage:   config.DockerImage,
		IndexPath:     config.IndexPath,
		TempDir:       config.TempDir,
		Timeout:       config.Timeout,
		UseDockerExec: config.UseDockerExec,
		ContainerName: config.ContainerName,
	}

	// Create solver client
	solverClient, err := solver.NewClient(solverCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:       config,
		solverClient: solverClient,
	}, nil
}

// Solve performs plate-solving on the given image file.
//
// This wraps the astrometry.net solve-field command, which identifies
// celestial coordinates and orientation by matching star patterns.
func (c *Client) Solve(ctx context.Context, imagePath string, opts *SolveOptions) (*Result, error) {
	return c.solverClient.Solve(ctx, imagePath, opts)
}

// SolveBytes performs plate-solving on image data provided as bytes.
//
// The data is written to a temporary file, solved, and cleaned up.
func (c *Client) SolveBytes(ctx context.Context, data []byte, format string, opts *SolveOptions) (*Result, error) {
	return c.solverClient.SolveBytes(ctx, data, format, opts)
}

// Future methods to be added:
// - ExtractSources(ctx, imagePath) - wraps image2xy
// - FitWCS(ctx, xyList) - wraps fit-wcs
// - XYToRaDec(ctx, wcsFile, x, y) - wraps wcs-xy2rd
// - RaDecToXY(ctx, wcsFile, ra, dec) - wraps wcs-rd2xy
// - AnalyzeFOV(imagePath) - analyze image and calculate FOV
