// Package fov provides field of view calculations and sensor size detection
// for astronomical imaging.
//
// This package helps determine the angular field of view for various camera
// and lens combinations, and recommends appropriate astrometry index files
// for plate-solving.
//
// # Calculating FOV
//
//	// Calculate FOV for a specific setup
//	fov := fov.CalculateFOV(200, fov.APSCCanon)
//	fmt.Printf("FOV: %.2f° x %.2f°\n", fov.WidthDegrees, fov.HeightDegrees)
//
// # Analyzing Images
//
//	// Extract camera details from EXIF and calculate FOV
//	info, err := fov.AnalyzeImage("photo.jpg")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Camera: %s %s\n", info.Make, info.Model)
//	fmt.Printf("Recommended scale: %.0f-%.0f arcminwidth\n",
//		info.ScaleLow, info.ScaleHigh)
//
// # Recommending Index Files
//
//	// Find which index files are needed for a lens range
//	rec := fov.RecommendIndexesForLens(50, 300, fov.APSCCanon, 1.2)
//	fmt.Println(rec.String())
package fov

import (
	"fmt"
	"math"
)

// SensorSize represents the physical dimensions of a camera sensor in millimeters.
type SensorSize struct {
	Width  float64 // Sensor width in mm
	Height float64 // Sensor height in mm
	Name   string  // Descriptive name
}

// Common sensor sizes
var (
	// Full Frame (35mm)
	FullFrame = SensorSize{Width: 36.0, Height: 24.0, Name: "Full Frame (35mm)"}

	// APS-C variants
	APSCCanon = SensorSize{Width: 22.3, Height: 14.9, Name: "APS-C Canon"}
	APSCNikon = SensorSize{Width: 23.6, Height: 15.7, Name: "APS-C Nikon/Sony"}
	APSCFuji  = SensorSize{Width: 23.5, Height: 15.6, Name: "APS-C Fujifilm"}

	// Micro Four Thirds
	MicroFourThirds = SensorSize{Width: 17.3, Height: 13.0, Name: "Micro Four Thirds"}

	// Smaller sensors
	OneInch = SensorSize{Width: 13.2, Height: 8.8, Name: "1\" sensor"}
)

// FieldOfView represents the calculated field of view for an imaging setup.
type FieldOfView struct {
	WidthDegrees  float64 // FOV width in degrees
	HeightDegrees float64 // FOV height in degrees
	WidthArcmin   float64 // FOV width in arcminutes
	HeightArcmin  float64 // FOV height in arcminutes
	DiagonalDeg   float64 // Diagonal FOV in degrees
}

// CalculateFOV calculates the field of view for a given focal length and sensor size.
//
// Formula: FOV (radians) = 2 * arctan(sensor_dimension / (2 * focal_length))
//
// Example:
//
//	fov := fov.CalculateFOV(50, fov.APSCNikon)
//	fmt.Printf("FOV: %.2f° x %.2f°\n", fov.WidthDegrees, fov.HeightDegrees)
func CalculateFOV(focalLengthMM float64, sensor SensorSize) FieldOfView {
	// Calculate FOV in radians
	fovWidthRad := 2 * math.Atan(sensor.Width/(2*focalLengthMM))
	fovHeightRad := 2 * math.Atan(sensor.Height/(2*focalLengthMM))

	// Calculate diagonal
	diagonalMM := math.Sqrt(sensor.Width*sensor.Width + sensor.Height*sensor.Height)
	fovDiagonalRad := 2 * math.Atan(diagonalMM/(2*focalLengthMM))

	// Convert to degrees
	widthDeg := fovWidthRad * 180 / math.Pi
	heightDeg := fovHeightRad * 180 / math.Pi
	diagonalDeg := fovDiagonalRad * 180 / math.Pi

	return FieldOfView{
		WidthDegrees:  widthDeg,
		HeightDegrees: heightDeg,
		WidthArcmin:   widthDeg * 60,
		HeightArcmin:  heightDeg * 60,
		DiagonalDeg:   diagonalDeg,
	}
}

// CalculateFOVRange calculates the FOV range for a zoom lens.
func CalculateFOVRange(minFocalLength, maxFocalLength float64, sensor SensorSize) (minFOV, maxFOV FieldOfView) {
	// At max focal length, FOV is smallest
	minFOV = CalculateFOV(maxFocalLength, sensor)
	// At min focal length, FOV is largest
	maxFOV = CalculateFOV(minFocalLength, sensor)
	return
}

// String returns a human-readable representation of the FOV.
func (fov FieldOfView) String() string {
	return fmt.Sprintf("%.2f° x %.2f° (%.1f' x %.1f')",
		fov.WidthDegrees, fov.HeightDegrees,
		fov.WidthArcmin, fov.HeightArcmin)
}
