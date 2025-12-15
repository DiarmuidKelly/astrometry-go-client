//go:build integration
// +build integration

package solver

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// GroundTruth represents the expected solve results for a test image
type GroundTruth struct {
	Description string `json:"description"`
	Source      string `json:"source"`
	Camera      struct {
		LensMM                  float64 `json:"lens_mm"`
		Mount                   string  `json:"mount"`
		Sensor                  string  `json:"sensor"`
		CropFactor              float64 `json:"crop_factor"`
		EffectiveFocalLengthMM  float64 `json:"effective_focal_length_mm"`
		Notes                   string  `json:"notes"`
	} `json:"camera"`
	Solution struct {
		RA                       float64 `json:"ra"`
		Dec                      float64 `json:"dec"`
		PixelScaleArcsecPerPixel float64 `json:"pixel_scale_arcsec_per_pixel"`
		RotationDegrees          float64 `json:"rotation_degrees"`
		FieldWidthDegrees        float64 `json:"field_width_degrees"`
		FieldHeightDegrees       float64 `json:"field_height_degrees"`
		ImageWidthPixels         int     `json:"image_width_pixels"`
		ImageHeightPixels        int     `json:"image_height_pixels"`
	} `json:"solution"`
	Tolerance struct {
		PositionArcsec     float64 `json:"position_arcsec"`
		PixelScalePercent  float64 `json:"pixel_scale_percent"`
		RotationDegrees    float64 `json:"rotation_degrees"`
		FieldSizePercent   float64 `json:"field_size_percent"`
	} `json:"tolerance"`
	Notes []string `json:"notes"`
}

// loadGroundTruth loads ground truth data from testdata/ground_truth.json
func loadGroundTruth(t *testing.T, imageFilename string) *GroundTruth {
	groundTruthPath := filepath.Join("testdata", "ground_truth.json")

	data, err := os.ReadFile(groundTruthPath)
	if err != nil {
		t.Fatalf("Failed to read ground truth file: %v", err)
	}

	var groundTruths map[string]GroundTruth
	if err := json.Unmarshal(data, &groundTruths); err != nil {
		t.Fatalf("Failed to parse ground truth JSON: %v", err)
	}

	gt, ok := groundTruths[imageFilename]
	if !ok {
		t.Fatalf("No ground truth found for %s", imageFilename)
	}

	return &gt
}

// validateResult checks if a solve result matches ground truth within tolerance
func validateResult(t *testing.T, result *Result, gt *GroundTruth) {
	t.Helper()

	if !result.Solved {
		t.Fatal("Image was not solved")
	}

	// Validate RA (convert to arcseconds for comparison)
	raDiff := math.Abs(result.RA - gt.Solution.RA) * 3600.0
	if raDiff > gt.Tolerance.PositionArcsec {
		t.Errorf("RA difference %.2f arcsec exceeds tolerance %.2f arcsec (got %.6f, expected %.6f)",
			raDiff, gt.Tolerance.PositionArcsec, result.RA, gt.Solution.RA)
	}

	// Validate Dec (convert to arcseconds)
	decDiff := math.Abs(result.Dec - gt.Solution.Dec) * 3600.0
	if decDiff > gt.Tolerance.PositionArcsec {
		t.Errorf("Dec difference %.2f arcsec exceeds tolerance %.2f arcsec (got %.6f, expected %.6f)",
			decDiff, gt.Tolerance.PositionArcsec, result.Dec, gt.Solution.Dec)
	}

	// Validate pixel scale (percent difference)
	pixelScaleDiff := math.Abs((result.PixelScale - gt.Solution.PixelScaleArcsecPerPixel) / gt.Solution.PixelScaleArcsecPerPixel * 100.0)
	if pixelScaleDiff > gt.Tolerance.PixelScalePercent {
		t.Errorf("Pixel scale difference %.2f%% exceeds tolerance %.2f%% (got %.2f, expected %.2f)",
			pixelScaleDiff, gt.Tolerance.PixelScalePercent, result.PixelScale, gt.Solution.PixelScaleArcsecPerPixel)
	}

	// Validate rotation (degrees)
	rotationDiff := math.Abs(result.Rotation - gt.Solution.RotationDegrees)
	if rotationDiff > gt.Tolerance.RotationDegrees {
		t.Errorf("Rotation difference %.2f° exceeds tolerance %.2f° (got %.2f, expected %.2f)",
			rotationDiff, gt.Tolerance.RotationDegrees, result.Rotation, gt.Solution.RotationDegrees)
	}

	// Validate field width (percent difference)
	fieldWidthDiff := math.Abs((result.FieldWidth - gt.Solution.FieldWidthDegrees) / gt.Solution.FieldWidthDegrees * 100.0)
	if fieldWidthDiff > gt.Tolerance.FieldSizePercent {
		t.Errorf("Field width difference %.2f%% exceeds tolerance %.2f%% (got %.4f, expected %.4f)",
			fieldWidthDiff, gt.Tolerance.FieldSizePercent, result.FieldWidth, gt.Solution.FieldWidthDegrees)
	}

	// Validate field height (percent difference)
	fieldHeightDiff := math.Abs((result.FieldHeight - gt.Solution.FieldHeightDegrees) / gt.Solution.FieldHeightDegrees * 100.0)
	if fieldHeightDiff > gt.Tolerance.FieldSizePercent {
		t.Errorf("Field height difference %.2f%% exceeds tolerance %.2f%% (got %.4f, expected %.4f)",
			fieldHeightDiff, gt.Tolerance.FieldSizePercent, result.FieldHeight, gt.Solution.FieldHeightDegrees)
	}

	t.Logf("Validation passed:")
	t.Logf("  RA:          %.6f° (diff: %.2f arcsec)", result.RA, raDiff)
	t.Logf("  Dec:         %.6f° (diff: %.2f arcsec)", result.Dec, decDiff)
	t.Logf("  Pixel scale: %.2f arcsec/px (diff: %.2f%%)", result.PixelScale, pixelScaleDiff)
	t.Logf("  Rotation:    %.2f° (diff: %.2f°)", result.Rotation, rotationDiff)
	t.Logf("  Field size:  %.4f° x %.4f°", result.FieldWidth, result.FieldHeight)
}

// TestM42WithGroundTruth tests all 3 Docker images against real M42 image with ground truth
func TestM42WithGroundTruth(t *testing.T) {
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Get index path
	indexPath := os.Getenv("ASTROMETRY_INDEX_PATH")
	if indexPath == "" {
		indexPath = filepath.Join(os.Getenv("HOME"), "astrometry-data")
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Skipf("Index path does not exist: %s. Set ASTROMETRY_INDEX_PATH or download indexes.", indexPath)
	}

	// Load test image and ground truth
	testImageFilename := "IMG_2820.JPG"
	testImagePath := filepath.Join("../../images", testImageFilename)

	if _, err := os.Stat(testImagePath); os.IsNotExist(err) {
		t.Fatalf("Test image not found: %s", testImagePath)
	}

	gt := loadGroundTruth(t, testImageFilename)
	t.Logf("Testing with: %s", gt.Description)
	t.Logf("Ground truth source: %s", gt.Source)

	// Test all Docker images
	images := []struct {
		name  string
		image string
	}{
		{"Latest", "diarmuidk/astrometry-dockerised-solver:latest"},
		{"Legacy", "dm90/astrometry:latest"},
	}

	results := make(map[string]*Result)

	for _, tc := range images {
		t.Run(tc.name, func(t *testing.T) {
			// Check if image is available
			checkCmd := exec.Command("docker", "image", "inspect", tc.image)
			if err := checkCmd.Run(); err != nil {
				t.Skipf("Image %s not available locally. Pull with: docker pull %s", tc.image, tc.image)
			}

			// Create client
			config := &ClientConfig{
				IndexPath:   indexPath,
				DockerImage: tc.image,
				Timeout:     3 * time.Minute,
			}

			client, err := NewClient(config)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Configure solve options - use reasonable scale range
			opts := DefaultSolveOptions()
			opts.ScaleLow = 1.0
			opts.ScaleHigh = 180.0
			opts.ScaleUnits = "degwidth"
			opts.DownsampleFactor = 2

			// Solve
			ctx := context.Background()
			result, err := client.Solve(ctx, testImagePath, opts)
			if err != nil {
				t.Fatalf("Solve failed: %v", err)
			}

			// Validate against ground truth
			validateResult(t, result, gt)

			// Store result for cross-comparison
			results[tc.name] = result
		})
	}

	// Cross-compare: all 3 images should produce nearly identical results
	if len(results) >= 2 {
		t.Run("CrossCompare", func(t *testing.T) {
			var refName string
			var refResult *Result

			// Get reference result (first available)
			for name, result := range results {
				refName = name
				refResult = result
				break
			}

			t.Logf("Using %s as reference for cross-comparison", refName)

			for name, result := range results {
				if name == refName {
					continue
				}

				// RA/Dec should match within 5 arcsec
				raDiff := math.Abs(result.RA - refResult.RA) * 3600.0
				decDiff := math.Abs(result.Dec - refResult.Dec) * 3600.0

				if raDiff > 5.0 {
					t.Errorf("%s vs %s: RA differs by %.2f arcsec (%.6f vs %.6f)",
						name, refName, raDiff, result.RA, refResult.RA)
				}

				if decDiff > 5.0 {
					t.Errorf("%s vs %s: Dec differs by %.2f arcsec (%.6f vs %.6f)",
						name, refName, decDiff, result.Dec, refResult.Dec)
				}

				// Pixel scale should match within 2%
				pixelScaleDiff := math.Abs((result.PixelScale - refResult.PixelScale) / refResult.PixelScale * 100.0)
				if pixelScaleDiff > 2.0 {
					t.Errorf("%s vs %s: Pixel scale differs by %.2f%% (%.2f vs %.2f)",
						name, refName, pixelScaleDiff, result.PixelScale, refResult.PixelScale)
				}

				t.Logf("%s matches %s within tolerance", name, refName)
			}
		})
	}
}

// TestConvertedJPEGWithGroundTruth tests converted standard JPEG against same ground truth
func TestConvertedJPEGWithGroundTruth(t *testing.T) {
	if !isDockerAvailable(t) {
		t.Skip("Docker is not available")
	}

	// Get index path
	indexPath := os.Getenv("ASTROMETRY_INDEX_PATH")
	if indexPath == "" {
		indexPath = filepath.Join(os.Getenv("HOME"), "astrometry-data")
	}

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Skipf("Index path does not exist: %s. Set ASTROMETRY_INDEX_PATH or download indexes.", indexPath)
	}

	// Test converted JPEG file
	convertedFilename := "IMG_2820-converted.jpg"
	convertedPath := filepath.Join("../../images", convertedFilename)

	if _, err := os.Stat(convertedPath); os.IsNotExist(err) {
		t.Skipf("Converted test image not found: %s", convertedPath)
	}

	// Use same ground truth as MPO version
	gt := loadGroundTruth(t, "IMG_2820.JPG")
	t.Logf("Testing converted JPEG: %s", convertedFilename)
	t.Logf("Using ground truth from: IMG_2820.JPG (MPO)")

	// Test with :latest image
	config := &ClientConfig{
		IndexPath:   indexPath,
		DockerImage: "diarmuidk/astrometry-dockerised-solver:latest",
		Timeout:     3 * time.Minute,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Configure solve options - use reasonable scale range
	opts := DefaultSolveOptions()
	opts.ScaleLow = 1.0
	opts.ScaleHigh = 180.0
	opts.ScaleUnits = "degwidth"
	opts.DownsampleFactor = 2

	// Solve
	ctx := context.Background()
	result, err := client.Solve(ctx, convertedPath, opts)
	if err != nil {
		t.Fatalf("Solve failed: %v", err)
	}

	// Validate against same ground truth as MPO
	validateResult(t, result, gt)

	t.Logf("Converted JPEG produces results within tolerance of MPO ground truth")
}
