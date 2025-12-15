package fov

import (
	"fmt"
	"sort"
)

// IndexFile represents an astrometry.net index file with its coverage and metadata.
type IndexFile struct {
	Name        string  // File name (e.g., "index-4110.fits")
	MinFOV      float64 // Minimum field width in degrees
	MaxFOV      float64 // Maximum field width in degrees
	SizeMB      float64 // File size in megabytes
	DownloadURL string  // URL to download the index
}

// AllIndexFiles contains metadata for all available 4100-series index files.
// These are ordered from widest to narrowest FOV.
var AllIndexFiles = []IndexFile{
	{Name: "index-4107", MinFOV: 8.0, MaxFOV: 11.0, SizeMB: 165, DownloadURL: "http://data.astrometry.net/4100/index-4107.fits"},
	{Name: "index-4108", MinFOV: 5.6, MaxFOV: 8.0, SizeMB: 95, DownloadURL: "http://data.astrometry.net/4100/index-4108.fits"},
	{Name: "index-4109", MinFOV: 4.2, MaxFOV: 5.6, SizeMB: 50, DownloadURL: "http://data.astrometry.net/4100/index-4109.fits"},
	{Name: "index-4110", MinFOV: 3.0, MaxFOV: 4.2, SizeMB: 25, DownloadURL: "http://data.astrometry.net/4100/index-4110.fits"},
	{Name: "index-4111", MinFOV: 2.2, MaxFOV: 3.0, SizeMB: 10, DownloadURL: "http://data.astrometry.net/4100/index-4111.fits"},
	{Name: "index-4112", MinFOV: 1.6, MaxFOV: 2.2, SizeMB: 5.3, DownloadURL: "http://data.astrometry.net/4100/index-4112.fits"},
	{Name: "index-4113", MinFOV: 1.1, MaxFOV: 1.6, SizeMB: 2.7, DownloadURL: "http://data.astrometry.net/4100/index-4113.fits"},
	{Name: "index-4114", MinFOV: 0.8, MaxFOV: 1.1, SizeMB: 1.4, DownloadURL: "http://data.astrometry.net/4100/index-4114.fits"},
	{Name: "index-4115", MinFOV: 0.56, MaxFOV: 0.8, SizeMB: 0.74, DownloadURL: "http://data.astrometry.net/4100/index-4115.fits"},
	{Name: "index-4116", MinFOV: 0.4, MaxFOV: 0.56, SizeMB: 0.409, DownloadURL: "http://data.astrometry.net/4100/index-4116.fits"},
	{Name: "index-4117", MinFOV: 0.28, MaxFOV: 0.4, SizeMB: 0.248, DownloadURL: "http://data.astrometry.net/4100/index-4117.fits"},
	{Name: "index-4118", MinFOV: 0.2, MaxFOV: 0.28, SizeMB: 0.187, DownloadURL: "http://data.astrometry.net/4100/index-4118.fits"},
	{Name: "index-4119", MinFOV: 0.1, MaxFOV: 0.2, SizeMB: 0.144, DownloadURL: "http://data.astrometry.net/4100/index-4119.fits"},
}

// IndexRecommendation contains recommended index files for a given FOV.
type IndexRecommendation struct {
	TargetFOV      FieldOfView
	Indexes        []IndexFile
	TotalSizeMB    float64
	DownloadScript string
}

// RecommendIndexes recommends index files for a given field of view.
// It returns 2-3 indexes that bracket the FOV for reliable solving.
//
// Parameters:
//   - fovDegrees: The field width in degrees
//   - margin: Additional margin as a multiplier (e.g., 1.5 means +50%)
//
// Example:
//
//	fov := fov.CalculateFOV(50, fov.APSCNikon)
//	rec := fov.RecommendIndexes(fov.WidthDegrees, 1.5)
//	fmt.Println(rec.DownloadScript)
func RecommendIndexes(fovDegrees, margin float64) IndexRecommendation {
	minFOV := fovDegrees / margin
	maxFOV := fovDegrees * margin

	var recommended []IndexFile
	for _, idx := range AllIndexFiles {
		// Include if there's overlap with our FOV range
		if idx.MaxFOV >= minFOV && idx.MinFOV <= maxFOV {
			recommended = append(recommended, idx)
		}
	}

	// Sort by MinFOV (narrowest to widest)
	sort.Slice(recommended, func(i, j int) bool {
		return recommended[i].MinFOV < recommended[j].MinFOV
	})

	// Calculate total size
	var totalSize float64
	for _, idx := range recommended {
		totalSize += idx.SizeMB
	}

	// Generate download script
	script := "#!/bin/bash\n# Download recommended astrometry index files\n\n"
	script += fmt.Sprintf("# Target FOV: %.2f degrees\n", fovDegrees)
	script += fmt.Sprintf("# Total download size: %.1f MB\n\n", totalSize)
	script += "mkdir -p astrometry-data && cd astrometry-data\n\n"
	for _, idx := range recommended {
		script += fmt.Sprintf("wget %s  # %.2f° - %.2f° (%.1f MB)\n",
			idx.DownloadURL, idx.MinFOV, idx.MaxFOV, idx.SizeMB)
	}

	return IndexRecommendation{
		Indexes:        recommended,
		TotalSizeMB:    totalSize,
		DownloadScript: script,
	}
}

// RecommendIndexesForFOV recommends indexes for a FieldOfView struct.
func RecommendIndexesForFOV(fov FieldOfView, margin float64) IndexRecommendation {
	rec := RecommendIndexes(fov.WidthDegrees, margin)
	rec.TargetFOV = fov
	return rec
}

// RecommendIndexesForLens recommends indexes for a specific lens and sensor combination.
//
// For zoom lenses, use minFocalLength and maxFocalLength.
// For prime lenses, use the same value for both parameters.
//
// Example:
//
//	// 50-300mm zoom on APS-C
//	rec := fov.RecommendIndexesForLens(50, 300, fov.APSCNikon, 1.3)
//	fmt.Println(rec.DownloadScript)
func RecommendIndexesForLens(minFocalLength, maxFocalLength float64, sensor SensorSize, margin float64) IndexRecommendation {
	minFOV, maxFOV := CalculateFOVRange(minFocalLength, maxFocalLength, sensor)

	// We need to cover from the narrowest FOV (maxFocalLength) to widest FOV (minFocalLength)
	// Apply margin to both ends
	fovMin := minFOV.WidthDegrees / margin
	fovMax := maxFOV.WidthDegrees * margin

	var recommended []IndexFile
	for _, idx := range AllIndexFiles {
		// Include if there's overlap with our FOV range
		if idx.MaxFOV >= fovMin && idx.MinFOV <= fovMax {
			recommended = append(recommended, idx)
		}
	}

	// Sort by MinFOV (narrowest to widest)
	sort.Slice(recommended, func(i, j int) bool {
		return recommended[i].MinFOV < recommended[j].MinFOV
	})

	// Calculate total size
	var totalSize float64
	for _, idx := range recommended {
		totalSize += idx.SizeMB
	}

	// Generate download script
	script := "#!/bin/bash\n# Download recommended astrometry index files\n\n"
	script += fmt.Sprintf("# Lens: %.0f-%.0fmm on %s\n", minFocalLength, maxFocalLength, sensor.Name)
	script += fmt.Sprintf("# FOV range: %s (wide) to %s (tele)\n", maxFOV.String(), minFOV.String())
	script += fmt.Sprintf("# Total download size: %.1f MB\n\n", totalSize)
	script += "mkdir -p astrometry-data && cd astrometry-data\n\n"
	for _, idx := range recommended {
		script += fmt.Sprintf("wget %s  # %.2f° - %.2f° (%.1f MB)\n",
			idx.DownloadURL, idx.MinFOV, idx.MaxFOV, idx.SizeMB)
	}

	return IndexRecommendation{
		TargetFOV:      maxFOV, // Use widest FOV as representative
		Indexes:        recommended,
		TotalSizeMB:    totalSize,
		DownloadScript: script,
	}
}

// String returns a human-readable summary of the recommendation.
func (r IndexRecommendation) String() string {
	result := fmt.Sprintf("Recommended indexes for FOV %s:\n", r.TargetFOV.String())
	result += fmt.Sprintf("Total download size: %.1f MB\n\n", r.TotalSizeMB)
	for _, idx := range r.Indexes {
		result += fmt.Sprintf("  %s: %.2f° - %.2f° (%.1f MB)\n",
			idx.Name, idx.MinFOV, idx.MaxFOV, idx.SizeMB)
	}
	return result
}
