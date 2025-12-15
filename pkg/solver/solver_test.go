package solver

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	if config.DockerImage != "diarmuidk/astrometry-dockerised-solver" {
		t.Errorf("expected DockerImage to be 'diarmuidk/astrometry-dockerised-solver', got '%s'", config.DockerImage)
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
		DockerImage: "diarmuidk/astrometry-dockerised-solver",
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
	if client.config.DockerImage != "diarmuidk/astrometry-dockerised-solver" {
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
	// Values based on real IMG_2820.JPG solve (testdata/wcs.fits)
	lines := []string{
		"SIMPLE  =                    T / file does conform to FITS standard             ",
		"BITPIX  =                    8 / number of bits per data pixel                  ",
		"NAXIS   =                    0 / number of data axes                            ",
		"EXTEND  =                    T / FITS dataset may contain extensions            ",
		"CRPIX1  =          3000.000000 / X reference pixel (image center)               ",
		"CRPIX2  =          2000.000000 / Y reference pixel (image center)               ",
		"CRVAL1  =        83.4230000000 / RA of reference point (deg)                    ",
		"CRVAL2  =        -5.8930000000 / Dec of reference point (deg)                   ",
		"CD1_1   =  -0.0010995000000000 / Transformation matrix                          ",
		"CD1_2   =   0.0004600000000000 / Transformation matrix (~22deg rotation)        ",
		"CD2_1   =  -0.0004500000000000 / Transformation matrix                          ",
		"CD2_2   =  -0.0011000000000000 / Transformation matrix                          ",
		"CTYPE1  = 'RA---TAN'           / WCS projection type                            ",
		"CTYPE2  = 'DEC--TAN'           / WCS projection type                            ",
		"IMAGEW  =                 6000 / Image width                                    ",
		"IMAGEH  =                 4000 / Image height                                   ",
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
	// Reference pixel is at image center (3000, 2000) = (IMAGEW/2, IMAGEH/2)
	// So CRVAL should equal field center (no transformation needed)
	expectedRA := 83.423
	expectedDec := -5.893

	if math.Abs(result.RA-expectedRA) > 0.001 {
		t.Errorf("expected RA to be approximately %.3f, got %.6f", expectedRA, result.RA)
	}

	if math.Abs(result.Dec-expectedDec) > 0.001 {
		t.Errorf("expected Dec to be approximately %.3f, got %.6f", expectedDec, result.Dec)
	}

	// Pixel scale from CD matrix: sqrt(CD1_1^2 + CD2_1^2) * 3600
	// sqrt(0.0010995^2 + 0.00045^2) * 3600 ≈ 4.3 arcsec/pixel
	expectedPixelScale := 4.3
	if math.Abs(result.PixelScale-expectedPixelScale) > 0.2 {
		t.Errorf("expected PixelScale to be approximately %.1f arcsec/px, got %.2f", expectedPixelScale, result.PixelScale)
	}

	// Rotation from CD matrix: atan2(CD1_2, CD1_1) with CD1_2=0.00046, CD1_1=-0.0010995
	// atan2(0.00046, -0.0010995) ≈ 2.76 rad ≈ 158°, so 180-158 ≈ 22°
	expectedRotation := 22.0
	if math.Abs(result.Rotation-expectedRotation) > 2.0 {
		t.Errorf("expected Rotation to be approximately %.0f°, got %.1f°", expectedRotation, result.Rotation)
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
