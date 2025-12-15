package fov

// Camera model to sensor size mappings
//
// IMPORTANT: These mappings are based on verified Wikipedia lists (accessed 2025-12-14):
// - APS-C cameras: https://en.wikipedia.org/wiki/Category:APS-C_digital_cameras
// - Full-frame DSLR: https://en.wikipedia.org/wiki/Category:Full-frame_DSLR_cameras
// - Full-frame mirrorless: https://en.wikipedia.org/wiki/Category:Full-frame_mirrorless_interchangeable_lens_cameras
// - Micro Four Thirds: https://en.wikipedia.org/wiki/Category:4/3-type_digital_cameras
//
// Detection works by checking if the camera model (converted to uppercase)
// contains any of the strings in these lists.

// CameraMapping represents a camera model pattern and its sensor size
type CameraMapping struct {
	Pattern string     // Substring to match in camera model (will be uppercased)
	Sensor  SensorSize // Sensor size for this camera
}

var (
	// canonMappings contains Canon camera model to sensor mappings
	// Source: Wikipedia Full-frame mirrorless & DSLR categories
	canonMappings = []CameraMapping{
		// Full Frame Mirrorless
		{Pattern: "EOS C50", Sensor: FullFrame},
		{Pattern: "EOS R1", Sensor: FullFrame},
		{Pattern: "EOS R3", Sensor: FullFrame},
		{Pattern: "EOS R5 MARK II", Sensor: FullFrame},
		{Pattern: "EOS R5", Sensor: FullFrame},
		{Pattern: "EOS R6 MARK II", Sensor: FullFrame},
		{Pattern: "EOS R6", Sensor: FullFrame},
		{Pattern: "EOS R8", Sensor: FullFrame},
		{Pattern: "EOS RA", Sensor: FullFrame},
		{Pattern: "EOS RP", Sensor: FullFrame},

		// Full Frame DSLR
		{Pattern: "EOS-1D C", Sensor: FullFrame},
		{Pattern: "EOS-1D X MARK III", Sensor: FullFrame},
		{Pattern: "EOS-1D X MARK II", Sensor: FullFrame},
		{Pattern: "EOS-1D X", Sensor: FullFrame},
		{Pattern: "EOS-1DS MARK III", Sensor: FullFrame},
		{Pattern: "EOS-1DS MARK II", Sensor: FullFrame},
		{Pattern: "EOS-1DS", Sensor: FullFrame},
		{Pattern: "EOS 5D MARK IV", Sensor: FullFrame},
		{Pattern: "EOS 5D MARK III", Sensor: FullFrame},
		{Pattern: "EOS 5D MARK II", Sensor: FullFrame},
		{Pattern: "EOS 5DS", Sensor: FullFrame},
		{Pattern: "EOS 5D", Sensor: FullFrame},
		{Pattern: "EOS 6D MARK II", Sensor: FullFrame},
		{Pattern: "EOS 6D", Sensor: FullFrame},

		// APS-C models (verified by user)
		{Pattern: "EOS M", Sensor: APSCCanon},
		{Pattern: "EOS R7", Sensor: APSCCanon},
		{Pattern: "EOS R10", Sensor: APSCCanon},
		{Pattern: "EOS M5", Sensor: APSCCanon},
		{Pattern: "EOS 7D", Sensor: APSCCanon},
		{Pattern: "EOS 77D", Sensor: APSCCanon},
		{Pattern: "EOS 80D", Sensor: APSCCanon},
		{Pattern: "EOS 90D", Sensor: APSCCanon},
		{Pattern: "REBEL", Sensor: APSCCanon},
		{Pattern: "KISS", Sensor: APSCCanon},
	}

	// nikonMappings contains Nikon camera model to sensor mappings
	// Source: Wikipedia Full-frame mirrorless & DSLR categories
	// NOTE: APS-C patterns must be checked FIRST to avoid substring matches with full-frame models
	nikonMappings = []CameraMapping{
		// APS-C models (checked first to prevent false matches)
		{Pattern: "Z FC", Sensor: APSCNikon},
		{Pattern: "Z 50", Sensor: APSCNikon},
		{Pattern: "Z50", Sensor: APSCNikon},
		{Pattern: "D7500", Sensor: APSCNikon},
		{Pattern: "D7200", Sensor: APSCNikon},
		{Pattern: "D7100", Sensor: APSCNikon},
		{Pattern: "D5600", Sensor: APSCNikon},
		{Pattern: "D5500", Sensor: APSCNikon},
		{Pattern: "D5300", Sensor: APSCNikon},
		{Pattern: "D3500", Sensor: APSCNikon},
		{Pattern: "D3400", Sensor: APSCNikon},
		{Pattern: "D3300", Sensor: APSCNikon},
		{Pattern: "D500", Sensor: APSCNikon},

		// Full Frame Mirrorless (with space variants)
		{Pattern: "Z 5 II", Sensor: FullFrame},
		{Pattern: "Z5II", Sensor: FullFrame},
		{Pattern: "Z 5", Sensor: FullFrame},
		{Pattern: "Z5", Sensor: FullFrame},
		{Pattern: "Z 6 III", Sensor: FullFrame},
		{Pattern: "Z6III", Sensor: FullFrame},
		{Pattern: "Z 6 II", Sensor: FullFrame},
		{Pattern: "Z6II", Sensor: FullFrame},
		{Pattern: "Z 6", Sensor: FullFrame},
		{Pattern: "Z6", Sensor: FullFrame},
		{Pattern: "Z 7 II", Sensor: FullFrame},
		{Pattern: "Z7II", Sensor: FullFrame},
		{Pattern: "Z 7", Sensor: FullFrame},
		{Pattern: "Z7", Sensor: FullFrame},
		{Pattern: "Z 8", Sensor: FullFrame},
		{Pattern: "Z8", Sensor: FullFrame},
		{Pattern: "Z 9", Sensor: FullFrame},
		{Pattern: "Z9", Sensor: FullFrame},
		{Pattern: "Z F", Sensor: FullFrame},
		{Pattern: "ZF", Sensor: FullFrame},
		{Pattern: "Z R", Sensor: FullFrame},
		{Pattern: "ZR", Sensor: FullFrame},

		// Full Frame DSLR
		{Pattern: "D3S", Sensor: FullFrame},
		{Pattern: "D3X", Sensor: FullFrame},
		{Pattern: "D3 ", Sensor: FullFrame}, // Space to avoid matching D3500, D3400, D3300
		{Pattern: "D4S", Sensor: FullFrame},
		{Pattern: "D4 ", Sensor: FullFrame}, // Space to avoid matching other models
		{Pattern: "D5 ", Sensor: FullFrame}, // Space to avoid matching D5600, D5500, D5300, D500
		{Pattern: "D6 ", Sensor: FullFrame},
		{Pattern: "D600", Sensor: FullFrame},
		{Pattern: "D610", Sensor: FullFrame},
		{Pattern: "D700", Sensor: FullFrame},
		{Pattern: "D750 ", Sensor: FullFrame}, // Space to avoid matching D7500
		{Pattern: "D780", Sensor: FullFrame},
		{Pattern: "D800", Sensor: FullFrame},
		{Pattern: "D810A", Sensor: FullFrame},
		{Pattern: "D810", Sensor: FullFrame},
		{Pattern: "D850", Sensor: FullFrame},
		{Pattern: "DF ", Sensor: FullFrame},
	}

	// sonyMappings contains Sony camera model to sensor mappings
	// Source: Wikipedia Full-frame mirrorless & DSLR categories
	sonyMappings = []CameraMapping{
		// Full Frame Mirrorless
		{Pattern: "ILCE-1", Sensor: FullFrame},    // α1 internal code
		{Pattern: "A1", Sensor: FullFrame},        // α1
		{Pattern: "ILCE-7M5", Sensor: FullFrame},  // α7 V internal code
		{Pattern: "ILCE-7M4", Sensor: FullFrame},  // α7 IV internal code
		{Pattern: "ILCE-7M3", Sensor: FullFrame},  // α7 III internal code
		{Pattern: "ILCE-7M2", Sensor: FullFrame},  // α7 II internal code
		{Pattern: "ILCE-7", Sensor: FullFrame},    // α7 internal code
		{Pattern: "ILCE-7RM5", Sensor: FullFrame}, // α7R V internal code
		{Pattern: "ILCE-7RM4", Sensor: FullFrame}, // α7R IV internal code
		{Pattern: "ILCE-7RM3", Sensor: FullFrame}, // α7R III internal code
		{Pattern: "ILCE-7RM2", Sensor: FullFrame}, // α7R II internal code
		{Pattern: "ILCE-7SM3", Sensor: FullFrame}, // α7S III internal code
		{Pattern: "ILCE-7SM2", Sensor: FullFrame}, // α7S II internal code
		{Pattern: "ILCE-9M3", Sensor: FullFrame},  // α9 III internal code
		{Pattern: "ILCE-9", Sensor: FullFrame},    // α9 internal code
		{Pattern: "FX3", Sensor: FullFrame},
		{Pattern: "FX6", Sensor: FullFrame},

		// Full Frame DSLR
		{Pattern: "ALPHA 99", Sensor: FullFrame},
		{Pattern: "ALPHA 850", Sensor: FullFrame},
		{Pattern: "ALPHA 900", Sensor: FullFrame},
		{Pattern: "A99", Sensor: FullFrame},

		// APS-C models (user verified)
		{Pattern: "ILCE-6", Sensor: APSCNikon}, // A6xxx series internal codes
		{Pattern: "A6", Sensor: APSCNikon},     // A6000, A6100, A6400, A6600
		{Pattern: "ZV-E10", Sensor: APSCNikon},
	}

	// olympusMappings contains Olympus/OM System camera model to sensor mappings
	// Source: Wikipedia 4/3-type digital cameras category
	// All Olympus/OM System cameras use Micro Four Thirds sensor
	olympusMappings = []CameraMapping{
		// OM-D series
		{Pattern: "E-M1X", Sensor: MicroFourThirds},
		{Pattern: "E-M1 MARK III", Sensor: MicroFourThirds},
		{Pattern: "E-M1 MARK II", Sensor: MicroFourThirds},
		{Pattern: "E-M1", Sensor: MicroFourThirds},
		{Pattern: "E-M5 MARK III", Sensor: MicroFourThirds},
		{Pattern: "E-M5 MARK II", Sensor: MicroFourThirds},
		{Pattern: "E-M5", Sensor: MicroFourThirds},
		{Pattern: "E-M10 MARK IV", Sensor: MicroFourThirds},
		{Pattern: "E-M10 MARK III", Sensor: MicroFourThirds},
		{Pattern: "E-M10 MARK II", Sensor: MicroFourThirds},
		{Pattern: "E-M10", Sensor: MicroFourThirds},

		// OM System
		{Pattern: "OM-1 MARK II", Sensor: MicroFourThirds},
		{Pattern: "OM-1", Sensor: MicroFourThirds},
		{Pattern: "OM-3", Sensor: MicroFourThirds},
		{Pattern: "OM-5", Sensor: MicroFourThirds},

		// E-series (legacy)
		{Pattern: "E-1", Sensor: MicroFourThirds},
		{Pattern: "E-3", Sensor: MicroFourThirds},
		{Pattern: "E-5", Sensor: MicroFourThirds},
		{Pattern: "E-30", Sensor: MicroFourThirds},
		{Pattern: "E-300", Sensor: MicroFourThirds},
		{Pattern: "E-330", Sensor: MicroFourThirds},
		{Pattern: "E-400", Sensor: MicroFourThirds},
		{Pattern: "E-410", Sensor: MicroFourThirds},
		{Pattern: "E-420", Sensor: MicroFourThirds},
		{Pattern: "E-450", Sensor: MicroFourThirds},
		{Pattern: "E-500", Sensor: MicroFourThirds},
		{Pattern: "E-510", Sensor: MicroFourThirds},
		{Pattern: "E-520", Sensor: MicroFourThirds},
		{Pattern: "E-620", Sensor: MicroFourThirds},

		// PEN series
		{Pattern: "PEN-F", Sensor: MicroFourThirds},
		{Pattern: "E-P1", Sensor: MicroFourThirds},
		{Pattern: "E-P2", Sensor: MicroFourThirds},
		{Pattern: "E-P3", Sensor: MicroFourThirds},
		{Pattern: "E-P5", Sensor: MicroFourThirds},
		{Pattern: "E-P7", Sensor: MicroFourThirds},
		{Pattern: "E-PL1", Sensor: MicroFourThirds},
		{Pattern: "E-PL2", Sensor: MicroFourThirds},
		{Pattern: "E-PL3", Sensor: MicroFourThirds},
		{Pattern: "E-PL5", Sensor: MicroFourThirds},
		{Pattern: "E-PL6", Sensor: MicroFourThirds},
		{Pattern: "E-PL7", Sensor: MicroFourThirds},
		{Pattern: "E-PL9", Sensor: MicroFourThirds},
		{Pattern: "E-PM1", Sensor: MicroFourThirds},
		{Pattern: "E-PM2", Sensor: MicroFourThirds},
	}

	// panasonicMappings contains Panasonic camera model to sensor mappings
	// Source: Wikipedia 4/3-type digital cameras category
	// All listed Panasonic Lumix G/GH/GF/GM/GX cameras use Micro Four Thirds sensor
	panasonicMappings = []CameraMapping{
		// DC-G series (newer models)
		{Pattern: "DC-G9 II", Sensor: MicroFourThirds},
		{Pattern: "DC-G9", Sensor: MicroFourThirds},

		// DC-GH series
		{Pattern: "DC-GH6", Sensor: MicroFourThirds},
		{Pattern: "DC-GH5S", Sensor: MicroFourThirds},
		{Pattern: "DC-GH5M2", Sensor: MicroFourThirds},
		{Pattern: "DC-GH5", Sensor: MicroFourThirds},

		// DC-GX series
		{Pattern: "DC-GX850", Sensor: MicroFourThirds},
		{Pattern: "DC-GX800", Sensor: MicroFourThirds},

		// DMC-G series
		{Pattern: "DMC-G85", Sensor: MicroFourThirds},
		{Pattern: "DMC-G80", Sensor: MicroFourThirds},
		{Pattern: "DMC-G10", Sensor: MicroFourThirds},
		{Pattern: "DMC-G7", Sensor: MicroFourThirds},
		{Pattern: "DMC-G6", Sensor: MicroFourThirds},
		{Pattern: "DMC-G5", Sensor: MicroFourThirds},
		{Pattern: "DMC-G3", Sensor: MicroFourThirds},
		{Pattern: "DMC-G2", Sensor: MicroFourThirds},
		{Pattern: "DMC-G1", Sensor: MicroFourThirds},

		// DMC-GF series
		{Pattern: "DMC-GF7", Sensor: MicroFourThirds},
		{Pattern: "DMC-GF6", Sensor: MicroFourThirds},
		{Pattern: "DMC-GF5", Sensor: MicroFourThirds},
		{Pattern: "DMC-GF3", Sensor: MicroFourThirds},
		{Pattern: "DMC-GF2", Sensor: MicroFourThirds},
		{Pattern: "DMC-GF1", Sensor: MicroFourThirds},

		// DMC-GH series
		{Pattern: "DMC-GH4", Sensor: MicroFourThirds},
		{Pattern: "DMC-GH3", Sensor: MicroFourThirds},
		{Pattern: "DMC-GH2", Sensor: MicroFourThirds},
		{Pattern: "DMC-GH1", Sensor: MicroFourThirds},

		// DMC-GM series
		{Pattern: "DMC-GM5", Sensor: MicroFourThirds},
		{Pattern: "DMC-GM1", Sensor: MicroFourThirds},

		// DMC-GX series
		{Pattern: "DMC-GX8", Sensor: MicroFourThirds},
		{Pattern: "DMC-GX7", Sensor: MicroFourThirds},
		{Pattern: "DMC-GX1", Sensor: MicroFourThirds},

		// DMC-L series
		{Pattern: "DMC-L10", Sensor: MicroFourThirds},
		{Pattern: "DMC-L1", Sensor: MicroFourThirds},
	}
)
