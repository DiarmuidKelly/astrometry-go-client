package solver

import (
	"errors"
	"fmt"
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

		// Extract key fields
		switch key {
		case "CRVAL1": // RA of reference point
			if val, err := strconv.ParseFloat(valuePart, 64); err == nil {
				result.RA = val
			}
		case "CRVAL2": // Dec of reference point
			if val, err := strconv.ParseFloat(valuePart, 64); err == nil {
				result.Dec = val
			}
		case "CD1_1", "CDELT1": // Pixel scale information
			// CD1_1 is deg/pixel, convert to arcsec/pixel
			if val, err := strconv.ParseFloat(valuePart, 64); err == nil {
				result.PixelScale = abs(val) * 3600.0
			}
		case "CROTA2": // Rotation angle
			if val, err := strconv.ParseFloat(valuePart, 64); err == nil {
				result.Rotation = val
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
