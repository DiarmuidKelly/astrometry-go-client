#!/bin/bash
# Download Astrometry.net index files

set -e

echo "Astrometry Index Files Download"
echo "================================"
echo ""
echo "This will download ALL index files (4107-4119)"
echo "Total size: ~350 MB"
echo "Coverage: 0.1° to 11.0° field width"
echo ""

# Create directory
mkdir -p astrometry-data
cd astrometry-data

echo "Downloading to: $(pwd)"
echo ""

# Download all index files
wget -c http://data.astrometry.net/4100/index-4119.fits
wget -c http://data.astrometry.net/4100/index-4118.fits
wget -c http://data.astrometry.net/4100/index-4117.fits
wget -c http://data.astrometry.net/4100/index-4116.fits
wget -c http://data.astrometry.net/4100/index-4115.fits
wget -c http://data.astrometry.net/4100/index-4114.fits
wget -c http://data.astrometry.net/4100/index-4113.fits
wget -c http://data.astrometry.net/4100/index-4112.fits
wget -c http://data.astrometry.net/4100/index-4111.fits
wget -c http://data.astrometry.net/4100/index-4110.fits
wget -c http://data.astrometry.net/4100/index-4109.fits
wget -c http://data.astrometry.net/4100/index-4108.fits
wget -c http://data.astrometry.net/4100/index-4107.fits

echo ""
echo "Download complete!"
echo "Index files saved to: $(pwd)"
echo ""
echo "Use this path in your configuration:"
echo "  ASTROMETRY_INDEX_PATH=$(pwd)"
