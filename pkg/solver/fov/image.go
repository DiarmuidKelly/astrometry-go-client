package fov

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

const (
	// detectionSourceEXIF indicates sensor was detected from EXIF camera model
	detectionSourceEXIF = "exif"
	// detectionSourceDefault indicates default sensor was used
	detectionSourceDefault = "default"
)

// ImageInfo contains camera and lens information extracted from an image.
type ImageInfo struct {
	Make         string  // Camera manufacturer (e.g., "Canon")
	Model        string  // Camera model (e.g., "Canon EOS M50m2")
	FocalLength  float64 // Focal length in mm
	Sensor       SensorSize
	FOV          FieldOfView
	ScaleLow     float64 // Recommended lower scale bound (arcminwidth)
	ScaleHigh    float64 // Recommended upper scale bound (arcminwidth)
	HasEXIF      bool    // Whether EXIF data was found
	DetectedFrom string  // How sensor was detected ("exif" or "default")
}

// AnalyzeImage extracts camera information from an image file and calculates FOV.
//
// This function reads EXIF data to determine the camera model and focal length,
// then attempts to match the camera to a known sensor size. If successful,
// it calculates the field of view and recommends scale parameters for solving.
//
// Example:
//
//	info, err := fov.AnalyzeImage("photo.jpg")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Camera: %s\n", info.Model)
//	fmt.Printf("FOV: %s\n", info.FOV.String())
//	fmt.Printf("Use scale: %.0f-%.0f arcminwidth\n", info.ScaleLow, info.ScaleHigh)
func AnalyzeImage(imagePath string) (*ImageInfo, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	// Try to decode EXIF data
	x, err := exif.Decode(file)
	if err != nil {
		return &ImageInfo{
			HasEXIF: false,
		}, fmt.Errorf("failed to decode EXIF data: %w", err)
	}

	info := &ImageInfo{
		HasEXIF: true,
	}

	// Extract camera make
	if makeTag, makeErr := x.Get(exif.Make); makeErr == nil {
		if makeStr, strErr := makeTag.StringVal(); strErr == nil {
			info.Make = makeStr
		}
	}

	// Extract camera model
	if model, modelErr := x.Get(exif.Model); modelErr == nil {
		if modelStr, strErr := model.StringVal(); strErr == nil {
			info.Model = modelStr
		}
	}

	// Extract focal length
	if focalTag, focalErr := x.Get(exif.FocalLength); focalErr == nil {
		if num, denom, ratErr := focalTag.Rat2(0); ratErr == nil && denom != 0 {
			info.FocalLength = float64(num) / float64(denom)
		}
	}

	// Detect sensor size from camera model
	info.Sensor, info.DetectedFrom = detectSensor(info.Make, info.Model)

	// Calculate FOV if we have focal length and sensor size
	if info.FocalLength > 0 && info.Sensor.Width > 0 {
		info.FOV = CalculateFOV(info.FocalLength, info.Sensor)

		// Calculate recommended scale bounds with 20% margin
		margin := 1.2
		info.ScaleLow = info.FOV.WidthArcmin / margin
		info.ScaleHigh = info.FOV.WidthArcmin * margin
	}

	return info, nil
}

// detectSensor attempts to identify the sensor size based on camera make and model.
// Camera mappings are defined in constants.go and should be reviewed for accuracy.
func detectSensor(cameraMake, model string) (SensorSize, string) {
	// Normalize strings for comparison
	makeUpper := toUpper(cameraMake)
	modelUpper := toUpper(model)

	// Canon cameras
	if contains(makeUpper, "CANON") {
		for _, mapping := range canonMappings {
			if contains(modelUpper, mapping.Pattern) {
				return mapping.Sensor, detectionSourceEXIF
			}
		}
	}

	// Nikon cameras
	if contains(makeUpper, "NIKON") {
		for _, mapping := range nikonMappings {
			if contains(modelUpper, mapping.Pattern) {
				return mapping.Sensor, detectionSourceEXIF
			}
		}
	}

	// Sony cameras
	if contains(makeUpper, "SONY") {
		for _, mapping := range sonyMappings {
			if contains(modelUpper, mapping.Pattern) {
				return mapping.Sensor, detectionSourceEXIF
			}
		}
	}

	// Olympus/OM System
	if contains(makeUpper, "OLYMPUS") || contains(makeUpper, "OM SYSTEM") {
		for _, mapping := range olympusMappings {
			if contains(modelUpper, mapping.Pattern) {
				return mapping.Sensor, detectionSourceEXIF
			}
		}
		// Default to Micro Four Thirds for Olympus/OM System
		return MicroFourThirds, detectionSourceEXIF
	}

	// Panasonic cameras
	if contains(makeUpper, "PANASONIC") {
		for _, mapping := range panasonicMappings {
			if contains(modelUpper, mapping.Pattern) {
				return mapping.Sensor, detectionSourceEXIF
			}
		}
		// Default to Micro Four Thirds for Panasonic Lumix
		return MicroFourThirds, detectionSourceEXIF
	}

	// Default to APS-C Nikon as most common sensor size
	return APSCNikon, detectionSourceDefault
}

// String returns a human-readable summary of the image info.
func (i *ImageInfo) String() string {
	if !i.HasEXIF {
		return "No EXIF data found"
	}

	result := fmt.Sprintf("Camera: %s %s\n", i.Make, i.Model)
	if i.FocalLength > 0 {
		result += fmt.Sprintf("Focal Length: %.0fmm\n", i.FocalLength)
	}
	result += fmt.Sprintf("Sensor: %s (%s)\n", i.Sensor.Name, i.DetectedFrom)
	if i.FOV.WidthDegrees > 0 {
		result += fmt.Sprintf("FOV: %s\n", i.FOV.String())
		result += fmt.Sprintf("Recommended scale: %.0f-%.0f arcminwidth", i.ScaleLow, i.ScaleHigh)
	}
	return result
}

// Helper functions
func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c = c - 'a' + 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
