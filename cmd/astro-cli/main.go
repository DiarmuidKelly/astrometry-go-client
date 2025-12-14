// Package main provides a command-line interface for astrometric plate-solving.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/DiarmuidKelly/Astrometry-Go-Client/pkg/solver"
)

const version = "0.1.0"

func main() {
	// Define flags
	imagePath := flag.String("image", "", "Path to the image file to solve (required)")
	indexPath := flag.String("index-path", "", "Path to astrometry index files (required)")
	scaleLow := flag.Float64("scale-low", 0, "Lower bound of image scale")
	scaleHigh := flag.Float64("scale-high", 0, "Upper bound of image scale")
	scaleUnits := flag.String("scale-units", "arcminwidth", "Units for scale (degwidth, arcminwidth, arcsecperpix)")
	downsample := flag.Int("downsample", 2, "Downsample factor")
	ra := flag.Float64("ra", 0, "RA hint in degrees (optional)")
	dec := flag.Float64("dec", 0, "Dec hint in degrees (optional)")
	radius := flag.Float64("radius", 0, "Search radius in degrees (optional)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("astro-cli version %s\n", version)
		os.Exit(0)
	}

	// Validate required flags
	if *imagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: --image is required")
		flag.Usage()
		os.Exit(1)
	}

	if *indexPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --index-path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create client config
	config := &solver.ClientConfig{
		IndexPath: *indexPath,
	}

	client, err := solver.NewClient(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	// Create solve options
	opts := solver.DefaultSolveOptions()
	opts.ScaleLow = *scaleLow
	opts.ScaleHigh = *scaleHigh
	opts.ScaleUnits = *scaleUnits
	opts.DownsampleFactor = *downsample
	opts.RA = *ra
	opts.Dec = *dec
	opts.Radius = *radius
	opts.Verbose = *verbose

	// Solve the image
	ctx := context.Background()
	result, err := client.Solve(ctx, *imagePath, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving image: %v\n", err)
		os.Exit(1)
	}

	// Output result as JSON
	output := struct {
		Solved      bool              `json:"solved"`
		RA          float64           `json:"ra,omitempty"`
		Dec         float64           `json:"dec,omitempty"`
		PixelScale  float64           `json:"pixel_scale,omitempty"`
		Rotation    float64           `json:"rotation,omitempty"`
		FieldWidth  float64           `json:"field_width,omitempty"`
		FieldHeight float64           `json:"field_height,omitempty"`
		SolveTime   float64           `json:"solve_time,omitempty"`
		OutputFiles []string          `json:"output_files,omitempty"`
		WCSHeader   map[string]string `json:"wcs_header,omitempty"`
	}{
		Solved:      result.Solved,
		RA:          result.RA,
		Dec:         result.Dec,
		PixelScale:  result.PixelScale,
		Rotation:    result.Rotation,
		FieldWidth:  result.FieldWidth,
		FieldHeight: result.FieldHeight,
		SolveTime:   result.SolveTime,
		OutputFiles: result.OutputFiles,
		WCSHeader:   result.WCSHeader,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	if !result.Solved {
		os.Exit(1)
	}
}
