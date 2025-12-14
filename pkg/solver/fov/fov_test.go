package fov

import (
	"math"
	"testing"
)

func TestCalculateFOV(t *testing.T) {
	tests := []struct {
		name          string
		focalLength   float64
		sensor        SensorSize
		expectedWidth float64 // degrees
		tolerance     float64 // degrees
	}{
		{
			name:          "50mm on Full Frame",
			focalLength:   50,
			sensor:        FullFrame,
			expectedWidth: 39.6, // Known value
			tolerance:     0.5,
		},
		{
			name:          "50mm on APS-C Nikon",
			focalLength:   50,
			sensor:        APSCNikon,
			expectedWidth: 26.6, // Calculated from formula
			tolerance:     0.1,
		},
		{
			name:          "200mm on APS-C Nikon",
			focalLength:   200,
			sensor:        APSCNikon,
			expectedWidth: 6.75,
			tolerance:     0.1,
		},
		{
			name:          "300mm on APS-C Nikon",
			focalLength:   300,
			sensor:        APSCNikon,
			expectedWidth: 4.5,
			tolerance:     0.1,
		},
		{
			name:          "100mm on APS-C Nikon",
			focalLength:   100,
			sensor:        APSCNikon,
			expectedWidth: 13.46,
			tolerance:     0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fov := CalculateFOV(tt.focalLength, tt.sensor)

			if math.Abs(fov.WidthDegrees-tt.expectedWidth) > tt.tolerance {
				t.Errorf("FOV width = %.2f°, want %.2f° (±%.2f°)",
					fov.WidthDegrees, tt.expectedWidth, tt.tolerance)
			}

			// Verify arcmin conversion
			expectedArcmin := tt.expectedWidth * 60
			if math.Abs(fov.WidthArcmin-expectedArcmin) > tt.tolerance*60 {
				t.Errorf("FOV arcmin = %.1f', want %.1f'",
					fov.WidthArcmin, expectedArcmin)
			}

			// Verify aspect ratio is approximately maintained
			// Note: Due to trigonometric projection, exact aspect ratio isn't preserved
			// at wide angles, but should be close for normal focal lengths
			aspectRatio := tt.sensor.Width / tt.sensor.Height
			calculatedAspect := fov.WidthDegrees / fov.HeightDegrees
			if math.Abs(calculatedAspect-aspectRatio) > 0.05 {
				t.Errorf("Aspect ratio = %.3f, want %.3f (difference too large)",
					calculatedAspect, aspectRatio)
			}
		})
	}
}

func TestCalculateFOVRange(t *testing.T) {
	// Test 50-300mm zoom on APS-C
	minFOV, maxFOV := CalculateFOVRange(50, 300, APSCNikon)

	// At 300mm (max focal length), FOV should be smallest
	if minFOV.WidthDegrees > 5 || minFOV.WidthDegrees < 4 {
		t.Errorf("Min FOV (at 300mm) = %.2f°, expected ~4.5°", minFOV.WidthDegrees)
	}

	// At 50mm (min focal length), FOV should be largest
	if maxFOV.WidthDegrees > 28 || maxFOV.WidthDegrees < 25 {
		t.Errorf("Max FOV (at 50mm) = %.2f°, expected ~26.6°", maxFOV.WidthDegrees)
	}

	// Max should be larger than min
	if maxFOV.WidthDegrees <= minFOV.WidthDegrees {
		t.Error("Max FOV should be larger than min FOV")
	}
}

func TestFieldOfViewString(t *testing.T) {
	fov := CalculateFOV(50, APSCNikon)
	str := fov.String()

	// Should contain degrees and arcminutes
	if str == "" {
		t.Error("String() returned empty string")
	}

	// Should contain reasonable values (not checking exact format)
	t.Logf("FOV String representation: %s", str)
}

func TestSensorPresets(t *testing.T) {
	sensors := []SensorSize{
		FullFrame,
		APSCCanon,
		APSCNikon,
		APSCFuji,
		MicroFourThirds,
		OneInch,
	}

	for _, sensor := range sensors {
		t.Run(sensor.Name, func(t *testing.T) {
			// Verify sensor has reasonable dimensions
			if sensor.Width <= 0 || sensor.Height <= 0 {
				t.Errorf("Invalid sensor dimensions: %.1f x %.1f mm",
					sensor.Width, sensor.Height)
			}

			// Verify sensor has a name
			if sensor.Name == "" {
				t.Error("Sensor name is empty")
			}

			// Verify width > height (landscape orientation)
			if sensor.Width < sensor.Height {
				t.Errorf("Sensor width (%.1f) should be greater than height (%.1f)",
					sensor.Width, sensor.Height)
			}

			// Test that FOV calculation works
			fov := CalculateFOV(50, sensor)
			if fov.WidthDegrees <= 0 || fov.HeightDegrees <= 0 {
				t.Error("FOV calculation produced invalid results")
			}
		})
	}
}

// Benchmark FOV calculation
func BenchmarkCalculateFOV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateFOV(50, APSCNikon)
	}
}
