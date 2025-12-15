// Package solver provides a Go client for the Astrometry.net plate-solving service
// via Docker containers (diarmuidk/astrometry-dockerised-solver or dm90/astrometry).
//
// Plate-solving identifies the celestial coordinates and orientation of astronomical
// images by matching star patterns against index files. This package wraps the
// astrometry.net solver in a convenient Go API.
//
// # Basic Usage
//
//	config := &solver.ClientConfig{
//		IndexPath: "/path/to/index/files",
//	}
//
//	client, err := solver.NewClient(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	opts := solver.DefaultSolveOptions()
//	opts.ScaleLow = 300
//	opts.ScaleHigh = 500
//	opts.ScaleUnits = "arcminwidth"
//
//	result, err := client.Solve(context.Background(), "image.jpg", opts)
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
// 1. Docker Run Mode (default): Spawns a new container for each solve operation.
// Simpler but slower for multiple solves.
//
// 2. Docker Exec Mode: Uses an existing long-running container via docker exec.
// Faster for multiple solves, requires a running container.
//
//	config := &solver.ClientConfig{
//		IndexPath:     "/path/to/indexes",
//		UseDockerExec: true,
//		ContainerName: "astrometry-solver",
//	}
package solver

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Client is the main interface for astrometry.net plate solving.
type Client struct {
	config *ClientConfig
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
		config.Timeout = 5 * time.Minute
	}
	if config.TempDir == "" {
		config.TempDir = os.TempDir()
	}

	return &Client{config: config}, nil
}

// Solve performs plate-solving on the given image file.
func (c *Client) Solve(ctx context.Context, imagePath string, opts *SolveOptions) (*Result, error) {
	if opts == nil {
		opts = DefaultSolveOptions()
	}

	// Validate image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: image file does not exist: %s", ErrInvalidInput, imagePath)
	}

	// Get absolute paths
	absImagePath, err := filepath.Abs(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute image path: %w", err)
	}

	absIndexPath, err := filepath.Abs(c.config.IndexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute index path: %w", err)
	}

	// Create temp directory for this solve operation
	tempDir, err := os.MkdirTemp(c.config.TempDir, "astrometry-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	if !opts.KeepTempFiles {
		defer func() {
			if removeErr := os.RemoveAll(tempDir); removeErr != nil {
				log.Printf("warning: failed to remove temp directory: %v", removeErr)
			}
		}()
	} else {
		log.Printf("KeepTempFiles enabled: temp directory preserved at %s", tempDir)
	}

	// Copy image to temp directory (solve-field writes output alongside input)
	imageFilename := filepath.Base(absImagePath)
	tempImagePath := filepath.Join(tempDir, imageFilename)
	err = copyFile(absImagePath, tempImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to copy image to temp directory: %w", err)
	}

	// Build solve-field command arguments
	args := c.buildSolveArgs(imageFilename, tempDir, opts)

	// Build Docker command based on mode
	var dockerArgs []string
	if c.config.UseDockerExec {
		// Docker exec mode: use existing container
		dockerArgs = []string{"exec", c.config.ContainerName}
		dockerArgs = append(dockerArgs, args...)
	} else {
		// Docker run mode: spawn new container
		dockerArgs = []string{
			"run", "--rm",
			"-v", fmt.Sprintf("%s:/data", tempDir),
			"-v", fmt.Sprintf("%s:/usr/local/astrometry/data", absIndexPath),
			c.config.DockerImage,
		}
		dockerArgs = append(dockerArgs, args...)
	}

	// Create context with timeout
	solveCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	// Execute Docker command
	startTime := time.Now()
	cmd := exec.CommandContext(solveCtx, "docker", dockerArgs...)
	output, err := cmd.CombinedOutput()

	if solveCtx.Err() == context.DeadlineExceeded {
		return nil, ErrTimeout
	}

	if err != nil {
		// Check if it's a "no solution" case by examining output
		if strings.Contains(string(output), "Did not solve") ||
			strings.Contains(string(output), "Failed to solve") {
			return &Result{Solved: false}, nil
		}
		return nil, fmt.Errorf("%w: %v\nOutput: %s", ErrDockerFailed, err, string(output))
	}

	solveTime := time.Since(startTime).Seconds()

	// Parse WCS file
	wcsPath := filepath.Join(tempDir, strings.TrimSuffix(imageFilename, filepath.Ext(imageFilename))+".wcs")
	result, err := ParseWCSFile(wcsPath)
	if err != nil {
		// If WCS file doesn't exist, image wasn't solved
		if os.IsNotExist(err) {
			return &Result{Solved: false}, nil
		}
		return nil, err
	}

	result.SolveTime = solveTime

	// Collect output files
	result.OutputFiles = c.collectOutputFiles(tempDir, imageFilename)

	return result, nil
}

// SolveBytes performs plate-solving on image data provided as bytes.
// The data is written to a temporary file, solved, and cleaned up.
func (c *Client) SolveBytes(ctx context.Context, data []byte, format string, opts *SolveOptions) (*Result, error) {
	// Create temp file with appropriate extension
	tempFile, err := os.CreateTemp(c.config.TempDir, fmt.Sprintf("image-*.%s", format))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tempFile.Name()) //nolint:errcheck // Cleanup operation, error not critical
	}()

	// Write data
	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close() //nolint:errcheck // Best effort cleanup on error path
		return nil, fmt.Errorf("failed to write image data: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// Solve using the temp file
	return c.Solve(ctx, tempFile.Name(), opts)
}

// buildSolveArgs constructs the solve-field command arguments.
func (c *Client) buildSolveArgs(imageFilename, tempDir string, opts *SolveOptions) []string {
	args := []string{"solve-field"}

	// Scale bounds
	if opts.ScaleLow > 0 && opts.ScaleHigh > 0 {
		args = append(args, "-L", fmt.Sprintf("%.6f", opts.ScaleLow))
		args = append(args, "-H", fmt.Sprintf("%.6f", opts.ScaleHigh))
		args = append(args, "-u", opts.ScaleUnits)
	}

	// Downsample
	if opts.DownsampleFactor > 0 {
		args = append(args, "--downsample", fmt.Sprintf("%d", opts.DownsampleFactor))
	}

	// Depth
	if opts.DepthLow > 0 && opts.DepthHigh > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d-%d", opts.DepthLow, opts.DepthHigh))
	}

	// No plots
	if opts.NoPlots {
		args = append(args, "--no-plots")
	}

	// RA/Dec hint
	if opts.RA != 0 || opts.Dec != 0 {
		args = append(args, "--ra", fmt.Sprintf("%.6f", opts.RA))
		args = append(args, "--dec", fmt.Sprintf("%.6f", opts.Dec))
		if opts.Radius > 0 {
			args = append(args, "--radius", fmt.Sprintf("%.6f", opts.Radius))
		}
	}

	// Overwrite
	if opts.OverwriteExisting {
		args = append(args, "--overwrite")
	}

	// Verbose
	if !opts.Verbose {
		args = append(args, "--no-verify")
	}

	// Determine paths based on execution mode
	var workDir, imagePath string
	if c.config.UseDockerExec {
		// In exec mode, use the actual shared volume path
		workDir = tempDir
		imagePath = filepath.Join(tempDir, imageFilename)
	} else {
		// In run mode, paths are relative to /data mount
		workDir = "/data"
		imagePath = fmt.Sprintf("/data/%s", imageFilename)
	}

	// Output directory
	args = append(args, "--dir", workDir)

	// Image path
	args = append(args, imagePath)

	return args
}

// collectOutputFiles finds all output files generated by solve-field.
func (c *Client) collectOutputFiles(tempDir, imageFilename string) []string {
	baseName := strings.TrimSuffix(imageFilename, filepath.Ext(imageFilename))
	extensions := []string{".wcs", ".corr", ".solved", ".match", ".rdls", ".axy", "-indx.xyls"}

	var files []string
	for _, ext := range extensions {
		path := filepath.Join(tempDir, baseName+ext)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}
	return files
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
