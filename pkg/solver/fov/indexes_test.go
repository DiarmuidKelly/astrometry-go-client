package fov

import (
	"strings"
	"testing"
)

func TestRecommendIndexes(t *testing.T) {
	tests := []struct {
		name             string
		fovDegrees       float64
		margin           float64
		expectedMinCount int
		expectedMaxCount int
	}{
		{
			name:             "Narrow FOV (200mm lens ~7°)",
			fovDegrees:       7.0,
			margin:           1.5,
			expectedMinCount: 2,
			expectedMaxCount: 5,
		},
		{
			name:             "Medium-wide FOV (100mm lens ~13°)",
			fovDegrees:       13.0,
			margin:           1.3,
			expectedMinCount: 1, // At edge of coverage (13*1.3=16.9° exceeds max of 11°)
			expectedMaxCount: 3,
		},
		{
			name:             "Very narrow FOV (500mm lens ~2.7°)",
			fovDegrees:       2.7,
			margin:           1.5,
			expectedMinCount: 2,
			expectedMaxCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := RecommendIndexes(tt.fovDegrees, tt.margin)

			// Check we got a reasonable number of recommendations
			count := len(rec.Indexes)
			if count < tt.expectedMinCount || count > tt.expectedMaxCount {
				t.Errorf("Got %d recommendations, expected between %d and %d",
					count, tt.expectedMinCount, tt.expectedMaxCount)
			}

			// Verify total size is positive
			if rec.TotalSizeMB <= 0 {
				t.Error("Total size should be positive")
			}

			// Verify download script is not empty
			if rec.DownloadScript == "" {
				t.Error("Download script is empty")
			}

			// Verify download script contains wget commands
			if !strings.Contains(rec.DownloadScript, "wget") {
				t.Error("Download script should contain wget commands")
			}

			t.Logf("Recommendations for %.1f°: %d indexes, %.1f MB total",
				tt.fovDegrees, count, rec.TotalSizeMB)
		})
	}
}

func TestRecommendIndexesForFOV(t *testing.T) {
	// Use 200mm lens which has FOV within index coverage
	fov := CalculateFOV(200, APSCNikon)
	rec := RecommendIndexesForFOV(fov, 1.5)

	// Verify TargetFOV is set
	if rec.TargetFOV.WidthDegrees == 0 {
		t.Error("TargetFOV should be set")
	}

	// Verify we got recommendations (200mm ~6.75° should match indexes)
	if len(rec.Indexes) == 0 {
		t.Error("Should get at least one index recommendation for 200mm lens")
	}
}

func TestRecommendIndexesForLens(t *testing.T) {
	// Test 50-300mm zoom on APS-C Nikon
	rec := RecommendIndexesForLens(50, 300, APSCNikon, 1.3)

	// Should get several indexes to cover the wide range
	if len(rec.Indexes) < 3 {
		t.Errorf("50-300mm zoom should need at least 3 indexes, got %d", len(rec.Indexes))
	}

	// Verify total size is reasonable
	if rec.TotalSizeMB <= 0 || rec.TotalSizeMB > 500 {
		t.Errorf("Total size %.1f MB seems unreasonable", rec.TotalSizeMB)
	}

	// Verify download script mentions the lens
	if !strings.Contains(rec.DownloadScript, "50-300mm") {
		t.Error("Download script should mention lens focal range")
	}

	// Verify download script mentions sensor
	if !strings.Contains(rec.DownloadScript, APSCNikon.Name) {
		t.Error("Download script should mention sensor name")
	}

	t.Logf("50-300mm zoom recommendations:\n%s", rec.String())
}

func TestRecommendIndexesForLens_PrimeLens(t *testing.T) {
	// Test 200mm prime on APS-C (FOV ~6.75°, well within index coverage)
	rec := RecommendIndexesForLens(200, 200, APSCNikon, 1.5)

	// Should get 2-4 indexes for a prime lens
	if len(rec.Indexes) < 2 || len(rec.Indexes) > 5 {
		t.Errorf("200mm prime should need 2-5 indexes, got %d", len(rec.Indexes))
	}

	// Total size should be moderate
	if rec.TotalSizeMB > 400 {
		t.Errorf("Total size %.1f MB seems too large for prime lens", rec.TotalSizeMB)
	}

	t.Logf("200mm prime recommendations: %d indexes, %.1f MB",
		len(rec.Indexes), rec.TotalSizeMB)
}

func TestIndexRecommendationString(t *testing.T) {
	rec := RecommendIndexes(10.0, 1.5)
	str := rec.String()

	if str == "" {
		t.Error("String() should not be empty")
	}

	// Should contain size information
	if !strings.Contains(str, "MB") {
		t.Error("String should contain size information")
	}

	t.Logf("Recommendation string:\n%s", str)
}

func TestAllIndexFilesMetadata(t *testing.T) {
	// Verify all index files have valid metadata
	for _, idx := range AllIndexFiles {
		t.Run(idx.Name, func(t *testing.T) {
			// Check name
			if idx.Name == "" {
				t.Error("Index name is empty")
			}

			// Check FOV range
			if idx.MinFOV <= 0 || idx.MaxFOV <= 0 {
				t.Errorf("Invalid FOV range: %.2f - %.2f", idx.MinFOV, idx.MaxFOV)
			}
			if idx.MinFOV >= idx.MaxFOV {
				t.Errorf("MinFOV (%.2f) should be less than MaxFOV (%.2f)",
					idx.MinFOV, idx.MaxFOV)
			}

			// Check size
			if idx.SizeMB <= 0 {
				t.Errorf("Invalid size: %.2f MB", idx.SizeMB)
			}

			// Check download URL
			if !strings.HasPrefix(idx.DownloadURL, "http://data.astrometry.net/") {
				t.Errorf("Invalid download URL: %s", idx.DownloadURL)
			}
			if !strings.HasSuffix(idx.DownloadURL, ".fits") {
				t.Errorf("Download URL should end with .fits: %s", idx.DownloadURL)
			}
		})
	}

	// Verify indexes are ordered from widest to narrowest
	for i := 0; i < len(AllIndexFiles)-1; i++ {
		if AllIndexFiles[i].MinFOV < AllIndexFiles[i+1].MinFOV {
			t.Errorf("Indexes not properly ordered: %s (%.2f) before %s (%.2f)",
				AllIndexFiles[i].Name, AllIndexFiles[i].MinFOV,
				AllIndexFiles[i+1].Name, AllIndexFiles[i+1].MinFOV)
		}
	}
}

// Test real-world scenarios
func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		minFL       float64
		maxFL       float64
		sensor      SensorSize
		description string
	}{
		{
			name:        "DSLR kit lens",
			minFL:       18,
			maxFL:       55,
			sensor:      APSCCanon,
			description: "Canon APS-C with 18-55mm kit lens",
		},
		{
			name:        "Telephoto zoom",
			minFL:       70,
			maxFL:       200,
			sensor:      FullFrame,
			description: "Full frame with 70-200mm telephoto",
		},
		{
			name:        "Ultra-wide prime",
			minFL:       14,
			maxFL:       14,
			sensor:      FullFrame,
			description: "Full frame with 14mm ultra-wide",
		},
		{
			name:        "Superzoom",
			minFL:       24,
			maxFL:       240,
			sensor:      OneInch,
			description: "1\" sensor with 24-240mm superzoom",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			rec := RecommendIndexesForLens(scenario.minFL, scenario.maxFL, scenario.sensor, 1.3)

			// Note: Very wide FOVs (>11°) may not have matching indexes
			// Only fail if FOV is within available index range
			minFOV, _ := CalculateFOVRange(scenario.minFL, scenario.maxFL, scenario.sensor)
			if minFOV.WidthDegrees <= 11.0 && len(rec.Indexes) == 0 {
				t.Error("Should get at least one recommendation for FOV within index coverage")
			}

			t.Logf("%s:\n  Indexes: %d\n  Total size: %.1f MB\n  Coverage: %s",
				scenario.description,
				len(rec.Indexes),
				rec.TotalSizeMB,
				rec.TargetFOV.String())
		})
	}
}
