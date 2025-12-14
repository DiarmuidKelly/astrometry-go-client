# Index Files for 50-300mm Lens on APS-C

## Your Camera Setup

- **Lens**: 50-300mm zoom
- **Sensor**: APS-C (23.6mm x 15.7mm)
- **Field of View Range**:
  - At 50mm (wide): **26.56° x 17.85°**
  - At 300mm (tele): **4.50° x 3.00°**

## The Problem

**Your current index files (4110-4113) only cover 1.1° - 4.2° FOV.**

This means:
- ✅ Images at **>220mm** focal length will solve (FOV < 4.2°)
- ❌ Images at **<220mm** focal length won't solve (FOV > 4.2°)

## Required Index Files

You need to download **wider** index files to cover your zoom range:

```bash
cd /home/diarmuid/code/Astrometry-Go-Client/astrometry-data

# Download indexes 4107-4110 (335 MB total)
wget http://data.astrometry.net/4100/index-4107.fits  # 8.00° - 11.00° (165 MB)
wget http://data.astrometry.net/4100/index-4108.fits  # 5.60° - 8.00° (95 MB)
wget http://data.astrometry.net/4100/index-4109.fits  # 4.20° - 5.60° (50 MB)
wget http://data.astrometry.net/4100/index-4110.fits  # 3.00° - 4.20° (25 MB)

# Keep your existing narrow indexes for telephoto end
# index-4111.fits  # 2.2° - 3.0° (10 MB) - already have
# index-4112.fits  # 1.6° - 2.2° (5.3 MB) - already have
# index-4113.fits  # 1.1° - 1.6° (2.7 MB) - already have
```

## Coverage After Download

With all index files (4107-4113), you'll be able to solve images at:

| Focal Length | FOV Width | Status |
|--------------|-----------|--------|
| 50-135mm | 26.6° - 10.6° | ⚠️ **Partially outside coverage** (>11° won't solve) |
| 135-165mm | 10.6° - 8.6° | ✅ **Full coverage** (index-4107: 8-11°) |
| 165-200mm | 8.6° - 7.1° | ✅ **Full coverage** (index-4108: 5.6-8°) |
| 200-265mm | 7.1° - 5.4° | ✅ **Full coverage** (index-4108/4109) |
| 265-300mm | 5.4° - 4.5° | ✅ **Full coverage** (index-4109/4110) |

### Limitation at Wide End

**Images shot at 50-135mm (FOV > 11°) cannot be solved** because there are no standard 4100-series index files that cover FOVs wider than 11°.

To solve very wide images (50-135mm), you would need:
- 4200-series indexes (covers wider FOVs)
- Or shoot at >135mm focal length

## Using the FOV Calculator

You can use the new FOV helper to calculate exact focal lengths:

```go
package main

import (
    "fmt"
    "github.com/DiarmuidKelly/Astrometry-Go-Client/pkg/solver/fov"
)

func main() {
    // Check what FOV your 100mm lens produces
    myFOV := fov.CalculateFOV(100, fov.APSCNikon)
    fmt.Printf("100mm produces: %s\n", myFOV.String())

    // Get index recommendations
    rec := fov.RecommendIndexesForLens(50, 300, fov.APSCNikon, 1.3)
    fmt.Println(rec.DownloadScript)
}
```

## Practical Recommendations

### For Your Test Images

If your test images (IMG_2819.JPG, etc.) were shot at 50-135mm, they **won't solve** even with the new indexes because the FOV exceeds 11°.

**Solution**:
1. Download indexes 4107-4110 (above command)
2. Test with images shot at **>135mm focal length**
3. Or use images shot with a longer lens (200mm+)

### Complete Coverage Setup

To cover your entire 50-300mm range as much as possible:

```bash
cd /home/diarmuid/code/Astrometry-Go-Client/astrometry-data

# Download ALL recommended indexes (377 MB total)
wget http://data.astrometry.net/4100/index-4107.fits  # 165 MB
wget http://data.astrometry.net/4100/index-4108.fits  # 95 MB
wget http://data.astrometry.net/4100/index-4109.fits  # 50 MB
# index-4110, 4111, 4112, 4113 already downloaded
```

This gives you coverage from **3.0° to 11.0°** (roughly 135mm to 300mm+ on your setup).

## Summary

- ✅ **Download indexes 4107-4110** (335 MB)
- ✅ **Keep existing indexes 4110-4113**
- ⚠️ **Wide end (50-135mm) won't work** - no indexes available
- ✅ **Tele end (135-300mm) will work** with new indexes
- 💡 **Test with images shot at >135mm** for guaranteed success
