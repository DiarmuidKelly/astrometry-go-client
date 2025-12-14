package solver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	if config.DockerImage != "dm90/astrometry" {
		t.Errorf("expected DockerImage to be 'dm90/astrometry', got '%s'", config.DockerImage)
	}

	if config.Timeout.Minutes() != 5 {
		t.Errorf("expected Timeout to be 5 minutes, got %v", config.Timeout)
	}
}

func TestDefaultSolveOptions(t *testing.T) {
	opts := DefaultSolveOptions()

	if opts.ScaleUnits != "arcminwidth" {
		t.Errorf("expected ScaleUnits to be 'arcminwidth', got '%s'", opts.ScaleUnits)
	}

	if opts.DownsampleFactor != 2 {
		t.Errorf("expected DownsampleFactor to be 2, got %d", opts.DownsampleFactor)
	}

	if !opts.NoPlots {
		t.Error("expected NoPlots to be true")
	}
}

func TestNewClient_MissingIndexPath(t *testing.T) {
	config := &ClientConfig{
		DockerImage: "dm90/astrometry",
	}

	_, err := NewClient(config)
	if err == nil {
		t.Error("expected error when IndexPath is missing")
	}

	if !strings.Contains(err.Error(), "IndexPath is required") {
		t.Errorf("expected error about missing IndexPath, got: %v", err)
	}
}

func TestNewClient_NonexistentIndexPath(t *testing.T) {
	config := &ClientConfig{
		IndexPath: "/nonexistent/path/to/indexes",
	}

	_, err := NewClient(config)
	if err == nil {
		t.Error("expected error when IndexPath does not exist")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("expected error about nonexistent IndexPath, got: %v", err)
	}
}

func TestNewClient_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	config := &ClientConfig{
		IndexPath: tempDir,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	// Check defaults were set
	if client.config.DockerImage != "dm90/astrometry" {
		t.Errorf("expected DockerImage to be set to default")
	}

	if client.config.Timeout.Minutes() != 5 {
		t.Errorf("expected Timeout to be set to default")
	}
}

func TestParseWCSFile_Nonexistent(t *testing.T) {
	_, err := ParseWCSFile("/nonexistent/file.wcs")
	if err == nil {
		t.Error("expected error when WCS file does not exist")
	}
}

func TestParseWCSFile_Success(t *testing.T) {
	// Create a mock WCS file
	tempDir := t.TempDir()
	wcsPath := filepath.Join(tempDir, "test.wcs")

	// Create FITS format WCS content (80-character records)
	lines := []string{
		"SIMPLE  =                    T / file does conform to FITS standard             ",
		"BITPIX  =                    8 / number of bits per data pixel                  ",
		"NAXIS   =                    0 / number of data axes                            ",
		"EXTEND  =                    T / FITS dataset may contain extensions            ",
		"CRPIX1  =            2048.5000 / Coordinate reference pixel                     ",
		"CRPIX2  =            1534.5000 / Coordinate reference pixel                     ",
		"CRVAL1  =      120.12345678901 / RA of reference point (deg)                    ",
		"CRVAL2  =       45.98765432101 / Dec of reference point (deg)                   ",
		"CD1_1   =  -0.0003055555555556 / Transformation matrix                          ",
		"CD1_2   =   0.0000000000000000 / Transformation matrix                          ",
		"CD2_1   =   0.0000000000000000 / Transformation matrix                          ",
		"CD2_2   =   0.0003055555555556 / Transformation matrix                          ",
		"CROTA2  =              15.5000 / Rotation angle (deg)                           ",
		"CTYPE1  = 'RA---TAN'           / WCS projection type                            ",
		"CTYPE2  = 'DEC--TAN'           / WCS projection type                            ",
		"NAXIS1  =                 4096 / Image width                                    ",
		"NAXIS2  =                 3068 / Image height                                   ",
		"END                                                                             ",
	}

	var wcsContent []byte
	for _, line := range lines {
		if len(line) != 80 {
			t.Fatalf("FITS line must be exactly 80 characters, got %d for: %s", len(line), line)
		}
		wcsContent = append(wcsContent, []byte(line)...)
	}

	if err := os.WriteFile(wcsPath, wcsContent, 0644); err != nil {
		t.Fatalf("failed to create test WCS file: %v", err)
	}

	result, err := ParseWCSFile(wcsPath)
	if err != nil {
		t.Fatalf("unexpected error parsing WCS file: %v", err)
	}

	if !result.Solved {
		t.Error("expected Solved to be true")
	}

	// Check parsed values
	if result.RA != 120.12345678901 {
		t.Errorf("expected RA to be 120.12345678901, got %.11f", result.RA)
	}

	if result.Dec != 45.98765432101 {
		t.Errorf("expected Dec to be 45.98765432101, got %.11f", result.Dec)
	}

	// CD1_1 = -0.0003055555555556 deg/pixel = 1.1 arcsec/pixel (approx)
	expectedPixelScale := 0.0003055555555556 * 3600.0
	if result.PixelScale < expectedPixelScale*0.99 || result.PixelScale > expectedPixelScale*1.01 {
		t.Errorf("expected PixelScale to be approximately %.2f, got %.2f", expectedPixelScale, result.PixelScale)
	}

	if result.Rotation != 15.5 {
		t.Errorf("expected Rotation to be 15.5, got %.1f", result.Rotation)
	}

	// Check WCS header map
	if len(result.WCSHeader) == 0 {
		t.Error("expected WCSHeader to contain entries")
	}

	if result.WCSHeader["CTYPE1"] != "RA---TAN" {
		t.Errorf("expected CTYPE1 to be 'RA---TAN', got '%s'", result.WCSHeader["CTYPE1"])
	}
}

func TestBuildSolveArgs(t *testing.T) {
	tempDir := t.TempDir()
	config := &ClientConfig{
		IndexPath: tempDir,
	}

	client, _ := NewClient(config)

	opts := &SolveOptions{
		ScaleLow:         1.0,
		ScaleHigh:        3.0,
		ScaleUnits:       "arcminwidth",
		DownsampleFactor: 2,
		DepthLow:         10,
		DepthHigh:        20,
		NoPlots:          true,
		RA:               120.0,
		Dec:              45.0,
		Radius:           5.0,
	}

	args := client.buildSolveArgs("test.jpg", tempDir, opts)

	// Check that key arguments are present
	argsStr := strings.Join(args, " ")

	if !strings.Contains(argsStr, "--downsample 2") {
		t.Error("expected --downsample argument")
	}

	if !strings.Contains(argsStr, "--no-plots") {
		t.Error("expected --no-plots argument")
	}

	if !strings.Contains(argsStr, "--ra 120") {
		t.Error("expected --ra argument")
	}

	if !strings.Contains(argsStr, "--dec 45") {
		t.Error("expected --dec argument")
	}

	if !strings.Contains(argsStr, "/data/test.jpg") {
		t.Error("expected image path argument")
	}
}
