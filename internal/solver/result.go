package solver

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Result holds the plate-solving results.
type Result struct {
	// Solved indicates whether the image was successfully plate-solved.
	Solved bool

	// RA is the right ascension of the image center in degrees (J2000).
	RA float64

	// Dec is the declination of the image center in degrees (J2000).
	Dec float64

	// PixelScale is the image scale in arcseconds per pixel.
	PixelScale float64

	// Rotation is the field rotation in degrees.
	Rotation float64

	// FieldWidth is the field of view width in degrees.
	FieldWidth float64

	// FieldHeight is the field of view height in degrees.
	FieldHeight float64

	// WCSHeader contains the raw parsed WCS header fields.
	WCSHeader map[string]string

	// OutputFiles contains paths to generated output files (.wcs, .corr, etc.).
	OutputFiles []string

	// SolveTime is the duration of the solve operation.
	SolveTime float64 // seconds

	// RawOutput contains the raw stdout/stderr from solve-field.
	// Only populated when Verbose option is enabled.
	RawOutput string
}

var (
	// ErrNoSolution indicates that astrometry.net could not solve the image.
	ErrNoSolution = errors.New("no solution found")

	// ErrTimeout indicates that the solve operation exceeded the timeout.
	ErrTimeout = errors.New("solve operation timed out")

	// ErrDockerFailed indicates that the Docker command failed.
	ErrDockerFailed = errors.New("docker command failed")

	// ErrInvalidInput indicates invalid input parameters.
	ErrInvalidInput = errors.New("invalid input parameters")

	// ErrWCSParseFailed indicates failure to parse WCS output.
	ErrWCSParseFailed = errors.New("failed to parse WCS output")
)

// ParseWCSFile parses a FITS WCS header file and returns a Result.
// The WCS file uses FITS header format with fixed 80-character records.
func ParseWCSFile(wcsPath string) (*Result, error) {
	file, err := os.Open(wcsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open WCS file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck // Read-only file, close error not critical
	}()

	result := &Result{
		Solved:    true,
		WCSHeader: make(map[string]string),
	}

	// Read FITS file in 80-character records
	buf := make([]byte, 80)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		if n != 80 {
			break
		}

		line := string(buf)

		// Skip empty lines and END marker
		if len(line) == 0 || strings.HasPrefix(line, "END") {
			continue
		}

		// Parse FITS header line: "KEY     = VALUE / COMMENT"
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valuePart := strings.TrimSpace(parts[1])

		// Remove comment (everything after '/')
		if idx := strings.Index(valuePart, "/"); idx != -1 {
			valuePart = valuePart[:idx]
		}
		valuePart = strings.TrimSpace(valuePart)

		// Remove quotes from string values
		valuePart = strings.Trim(valuePart, "'")
		valuePart = strings.TrimSpace(valuePart)

		// Store in map
		result.WCSHeader[key] = valuePart
	}

	// Parse WCS transformation parameters
	var crval1, crval2, crpix1, crpix2 float64
	var cd11, cd12, cd21, cd22 float64
	var imageW, imageH float64
	var hasWCS bool

	// Extract all WCS parameters
	if val, ok := result.WCSHeader["CRVAL1"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			crval1 = v
			hasWCS = true
		}
	}
	if val, ok := result.WCSHeader["CRVAL2"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			crval2 = v
		}
	}
	if val, ok := result.WCSHeader["CRPIX1"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			crpix1 = v
		}
	}
	if val, ok := result.WCSHeader["CRPIX2"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			crpix2 = v
		}
	}
	if val, ok := result.WCSHeader["CD1_1"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			cd11 = v
		}
	}
	if val, ok := result.WCSHeader["CD1_2"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			cd12 = v
		}
	}
	if val, ok := result.WCSHeader["CD2_1"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			cd21 = v
		}
	}
	if val, ok := result.WCSHeader["CD2_2"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			cd22 = v
		}
	}

	// Get image dimensions (prefer IMAGEW/IMAGEH, fallback to NAXIS1/NAXIS2)
	if val, ok := result.WCSHeader["IMAGEW"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			imageW = v
		}
	} else if val, ok := result.WCSHeader["NAXIS1"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			imageW = v
		}
	}
	if val, ok := result.WCSHeader["IMAGEH"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			imageH = v
		}
	} else if val, ok := result.WCSHeader["NAXIS2"]; ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			imageH = v
		}
	}

	// Calculate field center coordinates using WCS transformation
	// The reference pixel (CRPIX) has coordinates CRVAL at that pixel
	// To get field center, transform from reference pixel to image center
	if hasWCS && imageW > 0 && imageH > 0 {
		// Image center in pixel coordinates
		centerX := imageW / 2.0
		centerY := imageH / 2.0

		// Offset from reference pixel to center
		dx := centerX - crpix1
		dy := centerY - crpix2

		// Apply CD matrix transformation to get sky coordinate offsets
		dRA := cd11*dx + cd12*dy
		dDec := cd21*dx + cd22*dy

		// Calculate field center coordinates
		result.RA = crval1 + dRA
		result.Dec = crval2 + dDec

		// Calculate pixel scale from CD matrix
		// Pixel scale = sqrt(CD1_1^2 + CD2_1^2) in degrees/pixel
		pixelScaleDeg := math.Sqrt(cd11*cd11 + cd21*cd21)
		result.PixelScale = pixelScaleDeg * 3600.0 // Convert to arcsec/pixel

		// Calculate rotation from CD matrix
		// Position angle: how many degrees E of N the "up" direction points
		// Formula: 180 - atan2(CD1_2, CD1_1) * 180/π
		rotationRad := math.Atan2(cd12, cd11)
		result.Rotation = 180.0 - (rotationRad * 180.0 / math.Pi)

		// Normalize to 0-360 range
		for result.Rotation < 0 {
			result.Rotation += 360.0
		}
		for result.Rotation >= 360 {
			result.Rotation -= 360.0
		}
	} else {
		// Fallback: use CRVAL as field center if we can't calculate
		result.RA = crval1
		result.Dec = crval2

		// Fallback pixel scale from CD1_1
		if cd11 != 0 {
			result.PixelScale = abs(cd11) * 3600.0
		}

		// Extract rotation angle if available
		if val, ok := result.WCSHeader["CROTA2"]; ok {
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				result.Rotation = v
			}
		}
	}

	// Calculate field dimensions from image size and pixel scale if available
	// Try IMAGEW/IMAGEH first (astrometry.net specific), then fall back to NAXIS1/NAXIS2
	if imagew, ok := result.WCSHeader["IMAGEW"]; ok {
		if w, err := strconv.ParseFloat(imagew, 64); err == nil && result.PixelScale > 0 {
			result.FieldWidth = (w * result.PixelScale) / 3600.0 // degrees
		}
	}
	if imageh, ok := result.WCSHeader["IMAGEH"]; ok {
		if h, err := strconv.ParseFloat(imageh, 64); err == nil && result.PixelScale > 0 {
			result.FieldHeight = (h * result.PixelScale) / 3600.0 // degrees
		}
	}

	// Fall back to NAXIS if IMAGEW/IMAGEH not found
	if naxis1, ok := result.WCSHeader["NAXIS1"]; ok {
		if w, err := strconv.ParseFloat(naxis1, 64); err == nil && result.PixelScale > 0 {
			result.FieldWidth = (w * result.PixelScale) / 3600.0 // degrees
		}
	}
	if naxis2, ok := result.WCSHeader["NAXIS2"]; ok {
		if h, err := strconv.ParseFloat(naxis2, 64); err == nil && result.PixelScale > 0 {
			result.FieldHeight = (h * result.PixelScale) / 3600.0 // degrees
		}
	}

	// Validate that we got essential fields
	if result.RA == 0 && result.Dec == 0 && result.PixelScale == 0 {
		return nil, fmt.Errorf("%w: no valid WCS fields found", ErrWCSParseFailed)
	}

	return result, nil
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
