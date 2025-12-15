# Test Data for Integration Tests

This directory contains real astronomical images and ground truth data for integration testing.

## Files

### Test Images (located in `/images/`)

#### IMG_2820.JPG
- **Location**: `/images/IMG_2820.JPG`
- **Subject**: Orion Nebula (M42) region
- **Camera**: Canon EOS M50m2
- **Lens**: 200mm EF lens
- **Format**: MPO (Multi-Picture Object) - Canon's proprietary JPEG format
- **Size**: 6000x4000 pixels (6.1 MB)
- **Purpose**: Real astronomical image in MPO format for validating plate-solving with actual user workflow

#### IMG_2820-converted.jpg
- **Location**: `/images/IMG_2820-converted.jpg`
- **Subject**: Same as IMG_2820.JPG
- **Format**: Standard JPEG (converted from MPO at 100% quality)
- **Size**: 6000x4000 pixels (8.3 MB)
- **Purpose**: Tests standard JPEG format compatibility

#### IMG_2820_labelled.jpeg
- **Location**: `/images/IMG_2820_labelled.jpeg`
- **Purpose**: Annotated version showing identified stars and celestial objects from Astrometry.net nova

### Test Metadata (located in `testdata/`)

#### wcs.fits
- **Source**: Reference WCS solution file from astrometry.net web
- **Format**: FITS file containing World Coordinate System (WCS) solution
- **Purpose**: Reference data for WCS format validation

#### ground_truth.json
- **Purpose**: Structured ground truth data from solving IMG_2820.JPG
- **Contents**:
  - RA/Dec coordinates (J2000)
  - Pixel scale (arcsec/pixel)
  - Field of view (degrees)
  - Rotation angle
  - Tolerance values for test validation

## Index Files Required

The integration tests require **only index-4110.fits** (24 MB) for the IMG_2820.JPG test image:

- **index-4110.fits** - Covers 3.0° to 4.2° field of view
- **Image FOV**: ~6.6° (200mm lens on APS-C sensor)
- **Why this index**: The solver successfully solves using index-4110 which brackets the 3.96 arcsec/pixel scale

Run `make test-integration-setup` to download the index file automatically.

**Note**: Additional index files can be added when testing images with different field of views.

## Usage in Tests

Integration tests use this data to:

1. **Load** IMG_2820.JPG (MPO format) as primary test input
2. **Solve** using Docker images (diarmuidk/astrometry-dockerised-solver:0.97, dm90/astrometry:latest)
3. **Validate** results match ground truth within tolerance
4. **Verify** all images produce consistent results
5. **Test** both MPO and standard JPEG format support

## Regenerating Ground Truth

If you need to regenerate ground truth data (run from repository root):

```bash
# 1. Ensure you have the required index file
make test-integration-setup

# 2. Solve image with Docker solver
docker run --rm \
  -v "$(pwd)/images:/data" \
  -v "$(pwd)/astrometry-data:/usr/local/astrometry/data" \
  diarmuidk/astrometry-dockerised-solver:latest \
  solve-field --no-plots --overwrite \
  --scale-units degwidth --scale-low 1.0 --scale-high 180.0 \
  --downsample 2 /data/IMG_2820.JPG

# 3. Extract solution values from solve output
# Field center: (RA,Dec) in degrees
# Pixel scale in arcsec/pixel
# Field size in degrees
# Rotation angle (up is X degrees E of N)

# 4. Update testdata/ground_truth.json with new values
```

## Notes

- **MPO format** is fully supported - files are detected and processed correctly
- **Tolerances** are set to account for minor differences between solver implementations
- **Position tolerance** (10 arcsec) allows for WCS calculation differences
- **Pixel scale tolerance** (5%) accounts for rounding and numerical precision
- All Docker images should agree with each other within these tolerances
