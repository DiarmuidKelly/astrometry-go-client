package fov

import (
	"os"
	"strings"
	"testing"
)

func TestDetectSensor(t *testing.T) {
	tests := []struct {
		name           string
		make           string
		model          string
		expectedSensor SensorSize
		expectedFrom   string
	}{
		// Canon Full Frame
		{
			name:           "Canon EOS R5",
			make:           "Canon",
			model:          "Canon EOS R5",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Canon EOS 5D Mark IV",
			make:           "Canon",
			model:          "Canon EOS 5D Mark IV",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Canon EOS 6D",
			make:           "Canon",
			model:          "Canon EOS 6D",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		// Canon APS-C
		{
			name:           "Canon EOS M50",
			make:           "Canon",
			model:          "Canon EOS M50",
			expectedSensor: APSCCanon,
			expectedFrom:   "exif",
		},
		{
			name:           "Canon EOS M50m2",
			make:           "Canon",
			model:          "Canon EOS M50m2",
			expectedSensor: APSCCanon,
			expectedFrom:   "exif",
		},
		{
			name:           "Canon EOS 7D",
			make:           "Canon",
			model:          "Canon EOS 7D Mark II",
			expectedSensor: APSCCanon,
			expectedFrom:   "exif",
		},
		{
			name:           "Canon EOS Rebel T7i",
			make:           "Canon",
			model:          "Canon EOS Rebel T7i",
			expectedSensor: APSCCanon,
			expectedFrom:   "exif",
		},
		// Nikon Full Frame
		{
			name:           "Nikon Z9",
			make:           "Nikon",
			model:          "Nikon Z 9",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Nikon D850",
			make:           "Nikon",
			model:          "Nikon D850",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Nikon Z6",
			make:           "NIKON CORPORATION",
			model:          "NIKON Z 6",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		// Nikon APS-C
		{
			name:           "Nikon Z50",
			make:           "Nikon",
			model:          "Nikon Z 50",
			expectedSensor: APSCNikon,
			expectedFrom:   "exif",
		},
		{
			name:           "Nikon D7500",
			make:           "Nikon",
			model:          "Nikon D7500",
			expectedSensor: APSCNikon,
			expectedFrom:   "exif",
		},
		{
			name:           "Nikon D500",
			make:           "Nikon",
			model:          "Nikon D500",
			expectedSensor: APSCNikon,
			expectedFrom:   "exif",
		},
		// Sony Full Frame
		{
			name:           "Sony A7 III",
			make:           "Sony",
			model:          "ILCE-7M3",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Sony A7R IV",
			make:           "SONY",
			model:          "ILCE-7RM4",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		{
			name:           "Sony A9 II",
			make:           "Sony",
			model:          "ILCE-9M2",
			expectedSensor: FullFrame,
			expectedFrom:   "exif",
		},
		// Sony APS-C
		{
			name:           "Sony A6400",
			make:           "Sony",
			model:          "ILCE-6400",
			expectedSensor: APSCNikon,
			expectedFrom:   "exif",
		},
		{
			name:           "Sony ZV-E10",
			make:           "Sony",
			model:          "ZV-E10",
			expectedSensor: APSCNikon,
			expectedFrom:   "exif",
		},
		// Olympus Micro Four Thirds
		{
			name:           "Olympus OM-D E-M1",
			make:           "OLYMPUS CORPORATION",
			model:          "E-M1",
			expectedSensor: MicroFourThirds,
			expectedFrom:   "exif",
		},
		{
			name:           "OM System OM-1",
			make:           "OM SYSTEM",
			model:          "OM-1",
			expectedSensor: MicroFourThirds,
			expectedFrom:   "exif",
		},
		// Panasonic Micro Four Thirds
		{
			name:           "Panasonic GH5",
			make:           "Panasonic",
			model:          "DC-GH5",
			expectedSensor: MicroFourThirds,
			expectedFrom:   "exif",
		},
		{
			name:           "Panasonic G9",
			make:           "Panasonic",
			model:          "DC-G9",
			expectedSensor: MicroFourThirds,
			expectedFrom:   "exif",
		},
		// Unknown camera - default to APS-C
		{
			name:           "Unknown Camera",
			make:           "Unknown Manufacturer",
			model:          "Unknown Model",
			expectedSensor: APSCNikon,
			expectedFrom:   "default",
		},
		{
			name:           "Empty strings",
			make:           "",
			model:          "",
			expectedSensor: APSCNikon,
			expectedFrom:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sensor, from := detectSensor(tt.make, tt.model)

			if sensor.Width != tt.expectedSensor.Width || sensor.Height != tt.expectedSensor.Height {
				t.Errorf("detectSensor(%q, %q) sensor = %+v, want %+v",
					tt.make, tt.model, sensor, tt.expectedSensor)
			}

			if from != tt.expectedFrom {
				t.Errorf("detectSensor(%q, %q) from = %q, want %q",
					tt.make, tt.model, from, tt.expectedFrom)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"HELLO", "HELLO"},
		{"HeLLo", "HELLO"},
		{"Canon EOS M50", "CANON EOS M50"},
		{"123abc", "123ABC"},
		{"", ""},
		{"MixedCase123", "MIXEDCASE123"},
		{"special!@#chars", "SPECIAL!@#CHARS"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toUpper(tt.input)
			if result != tt.expected {
				t.Errorf("toUpper(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{"substring exists", "Canon EOS M50", "EOS", true},
		{"substring at start", "Canon EOS M50", "Canon", true},
		{"substring at end", "Canon EOS M50", "M50", true},
		{"substring not found", "Canon EOS M50", "Nikon", false},
		{"empty substring", "Canon EOS M50", "", true},
		{"empty string", "", "test", false},
		{"both empty", "", "", true},
		{"substring longer than string", "EOS", "Canon EOS", false},
		{"exact match", "Canon", "Canon", true},
		{"case sensitive", "Canon", "canon", false},
		{"partial match", "ILCE-7M3", "A7", false},
		{"partial match success", "ILCE-7M3", "7M3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.str, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, want %v",
					tt.str, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestImageInfoString(t *testing.T) {
	tests := []struct {
		name     string
		info     ImageInfo
		contains []string
	}{
		{
			name: "Complete info",
			info: ImageInfo{
				Make:         "Canon",
				Model:        "EOS M50",
				FocalLength:  50.0,
				Sensor:       APSCCanon,
				FOV:          FieldOfView{WidthDegrees: 15.0, HeightDegrees: 10.0, WidthArcmin: 900, HeightArcmin: 600},
				ScaleLow:     750,
				ScaleHigh:    1080,
				HasEXIF:      true,
				DetectedFrom: "exif",
			},
			contains: []string{"Canon EOS M50", "50mm", "APS-C Canon", "exif", "15.00°", "750-1080 arcminwidth"},
		},
		{
			name: "No EXIF data",
			info: ImageInfo{
				HasEXIF: false,
			},
			contains: []string{"No EXIF"},
		},
		{
			name: "EXIF but no FOV",
			info: ImageInfo{
				Make:         "Nikon",
				Model:        "D850",
				Sensor:       FullFrame,
				HasEXIF:      true,
				DetectedFrom: "exif",
			},
			contains: []string{"Nikon D850", "Full Frame", "exif"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.String()

			for _, expectedSubstr := range tt.contains {
				if !strings.Contains(result, expectedSubstr) {
					t.Errorf("ImageInfo.String() missing expected substring %q\nGot: %s",
						expectedSubstr, result)
				}
			}
		})
	}
}

func TestAnalyzeImage_NonexistentFile(t *testing.T) {
	_, err := AnalyzeImage("/nonexistent/file.jpg")
	if err == nil {
		t.Error("AnalyzeImage() expected error for nonexistent file, got nil")
	}
}

func TestAnalyzeImage_InvalidEXIF(t *testing.T) {
	// Create a temporary file without EXIF data
	tmpFile := t.TempDir() + "/no-exif.txt"
	if err := os.WriteFile(tmpFile, []byte("not an image"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info, err := AnalyzeImage(tmpFile)
	if err == nil {
		t.Error("AnalyzeImage() expected error for non-image file, got nil")
	}
	if info == nil {
		t.Error("AnalyzeImage() should return partial ImageInfo even on error")
	}
	if info != nil && info.HasEXIF {
		t.Error("AnalyzeImage() HasEXIF should be false for file without EXIF")
	}
}
