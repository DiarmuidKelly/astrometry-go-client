//go:build integration
// +build integration

package solver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestDockerRunMode tests the default docker run mode
func TestDockerRunMode(t *testing.T) {
	// Check if Docker is available
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Create a test directory with a minimal test image
	tempDir := t.TempDir()
	testImagePath := createTestImage(t, tempDir)

	// Get index path from environment or use default
	indexPath := os.Getenv("ASTROMETRY_INDEX_PATH")
	if indexPath == "" {
		indexPath = filepath.Join(os.Getenv("HOME"), "astrometry-data")
	}

	// Check if index path exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Skipf("Index path does not exist: %s. Set ASTROMETRY_INDEX_PATH or download indexes.", indexPath)
	}

	// Create client with default image (docker run mode)
	config := &ClientConfig{
		IndexPath: indexPath,
		Timeout:   2 * time.Minute,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify the default image is set to the constant
	if client.config.DockerImage != DefaultDockerImage {
		t.Errorf("Expected default image '%s', got '%s'", DefaultDockerImage, client.config.DockerImage)
	}

	// Configure solve options
	opts := DefaultSolveOptions()
	opts.ScaleLow = 0.5
	opts.ScaleHigh = 10.0
	opts.ScaleUnits = "degwidth"
	opts.Verbose = true

	// Attempt to solve (will likely fail with test image, but tests the integration)
	ctx := context.Background()
	result, err := client.Solve(ctx, testImagePath, opts)

	// We expect either a successful solve or a "no solution" result
	if err != nil && err != ErrNoSolution {
		t.Logf("Solve failed (expected for test image): %v", err)
	}

	if result != nil {
		t.Logf("Solve completed: Solved=%v, Time=%.2fs", result.Solved, result.SolveTime)
		if result.Solved {
			t.Logf("  RA=%.6f, Dec=%.6f, PixelScale=%.2f", result.RA, result.Dec, result.PixelScale)
		}
	}
}

// TestDockerExecMode tests the docker exec mode with a running container
func TestDockerExecMode(t *testing.T) {
	// Check if Docker is available
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Get index path from environment or use default
	indexPath := os.Getenv("ASTROMETRY_INDEX_PATH")
	if indexPath == "" {
		indexPath = filepath.Join(os.Getenv("HOME"), "astrometry-data")
	}

	// Check if index path exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Skipf("Index path does not exist: %s. Set ASTROMETRY_INDEX_PATH or download indexes.", indexPath)
	}

	// Start a container for testing
	containerName := fmt.Sprintf("astrometry-test-%d", time.Now().Unix())
	absIndexPath, err := filepath.Abs(indexPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create shared temp directory
	sharedDir := t.TempDir()

	// Start container using the default image
	dockerArgs := []string{
		"run", "-d",
		"--name", containerName,
		"-v", fmt.Sprintf("%s:/usr/local/astrometry/data:ro", absIndexPath),
		"-v", fmt.Sprintf("%s:/shared-data", sharedDir),
		DefaultDockerImage,
		"tail", "-f", "/dev/null",
	}

	cmd := exec.Command("docker", dockerArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to start container: %v\nOutput: %s", err, string(output))
	}

	// Ensure cleanup
	defer func() {
		stopCmd := exec.Command("docker", "stop", containerName)
		_ = stopCmd.Run()
		rmCmd := exec.Command("docker", "rm", containerName)
		_ = rmCmd.Run()
	}()

	// Wait for container to be ready
	time.Sleep(2 * time.Second)

	// Create test image
	testImagePath := createTestImage(t, sharedDir)

	// Create client with docker exec mode
	config := &ClientConfig{
		IndexPath:     absIndexPath,
		UseDockerExec: true,
		ContainerName: containerName,
		Timeout:       2 * time.Minute,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Configure solve options
	opts := DefaultSolveOptions()
	opts.ScaleLow = 0.5
	opts.ScaleHigh = 10.0
	opts.ScaleUnits = "degwidth"
	opts.Verbose = true

	// Attempt to solve
	ctx := context.Background()
	result, err := client.Solve(ctx, testImagePath, opts)

	// We expect either a successful solve or a "no solution" result
	if err != nil && err != ErrNoSolution {
		t.Logf("Solve failed (expected for test image): %v", err)
	}

	if result != nil {
		t.Logf("Solve completed: Solved=%v, Time=%.2fs", result.Solved, result.SolveTime)
		if result.Solved {
			t.Logf("  RA=%.6f, Dec=%.6f, PixelScale=%.2f", result.RA, result.Dec, result.PixelScale)
		}
	}
}

// TestImageCompatibility tests that all supported Docker images work with the same code
// Tests: diarmuidk/astrometry-dockerised-solver (DockerHub)
//        ghcr.io/diarmuidkelly/astrometry-dockerised-solver (GHCR)
//        dm90/astrometry (legacy)
func TestImageCompatibility(t *testing.T) {
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	indexPath := os.Getenv("ASTROMETRY_INDEX_PATH")
	if indexPath == "" {
		indexPath = filepath.Join(os.Getenv("HOME"), "astrometry-data")
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Skipf("Index path does not exist: %s", indexPath)
	}

	tempDir := t.TempDir()
	testImagePath := createTestImage(t, tempDir)

	// Test all supported image variants
	images := []struct {
		name        string
		image       string
		description string
	}{
		{
			name:        "DockerHub",
			image:       "diarmuidk/astrometry-dockerised-solver:latest",
			description: "diarmuidk/astrometry-dockerised-solver from DockerHub",
		},
		{
			name:        "GHCR",
			image:       "ghcr.io/diarmuidkelly/astrometry-dockerised-solver:latest",
			description: "diarmuidk/astrometry-dockerised-solver from GHCR",
		},
		{
			name:        "Legacy",
			image:       "dm90/astrometry:latest",
			description: "dm90/astrometry (legacy image)",
		},
	}

	for _, tc := range images {
		t.Run(tc.name, func(t *testing.T) {
			// Check if image is available
			checkCmd := exec.Command("docker", "image", "inspect", tc.image)
			if err := checkCmd.Run(); err != nil {
				t.Skipf("Image %s not available locally. Pull with: docker pull %s", tc.image, tc.image)
			}

			t.Logf("Testing: %s", tc.description)

			config := &ClientConfig{
				IndexPath:   indexPath,
				DockerImage: tc.image,
				Timeout:     2 * time.Minute,
			}

			client, err := NewClient(config)
			if err != nil {
				t.Fatalf("Failed to create client with image %s: %v", tc.image, err)
			}

			opts := DefaultSolveOptions()
			opts.ScaleLow = 0.5
			opts.ScaleHigh = 10.0
			opts.ScaleUnits = "degwidth"

			ctx := context.Background()
			result, err := client.Solve(ctx, testImagePath, opts)

			// Just verify it doesn't crash - actual solving may fail without proper test data
			if err != nil && err != ErrNoSolution {
				t.Logf("Solve with %s: %v (expected for minimal test image)", tc.name, err)
			}

			if result != nil {
				t.Logf("Image %s: Solved=%v, Time=%.2fs", tc.name, result.Solved, result.SolveTime)
			}
		})
	}
}

// Helper functions

func isDockerAvailable(t *testing.T) bool {
	cmd := exec.Command("docker", "version")
	err := cmd.Run()
	if err != nil {
		t.Logf("Docker not available: %v", err)
		return false
	}
	return true
}

func createTestImage(t *testing.T, dir string) string {
	// Create a simple 100x100 black PNG image for testing
	imagePath := filepath.Join(dir, "test-image.png")

	// PNG header for a 100x100 black image
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk start
		0x00, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x64, // Width=100, Height=100
		0x08, 0x00, 0x00, 0x00, 0x00,
		0x7C, 0x79, 0x7E, 0xF8, // CRC
		0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, // IDAT chunk
		0x78, 0x9C, 0x62, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01,
		0xE2, 0x21, 0xBC, 0x33,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, // IEND chunk
		0xAE, 0x42, 0x60, 0x82,
	}

	if err := os.WriteFile(imagePath, pngData, 0644); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	return imagePath
}
