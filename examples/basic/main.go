package main

import (
	"context"
	"fmt"
	"log"

	solver "github.com/DiarmuidKelly/astrometry-go-client"
)

func main() {
	// Configure the client
	config := &solver.ClientConfig{
		IndexPath: "/path/to/astrometry/indexes", // Update with your index path
	}

	// Create a new client
	client, err := solver.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Configure solve options
	opts := solver.DefaultSolveOptions()
	opts.ScaleLow = 1.0  // 1 arcmin/width
	opts.ScaleHigh = 3.0 // 3 arcmin/width
	opts.DownsampleFactor = 2

	// Optional: provide RA/Dec hint for faster solving
	// opts.RA = 120.5
	// opts.Dec = 45.2
	// opts.Radius = 5.0

	// Solve the image
	ctx := context.Background()
	result, err := client.Solve(ctx, "path/to/your/image.jpg", opts)
	if err != nil {
		log.Fatalf("Failed to solve image: %v", err)
	}

	// Check if the image was solved
	if !result.Solved {
		log.Fatal("Image could not be solved")
	}

	// Print the results
	fmt.Printf("Image successfully solved!\n")
	fmt.Printf("RA (J2000):       %.6f degrees\n", result.RA)
	fmt.Printf("Dec (J2000):      %.6f degrees\n", result.Dec)
	fmt.Printf("Pixel Scale:      %.2f arcsec/pixel\n", result.PixelScale)
	fmt.Printf("Rotation:         %.2f degrees\n", result.Rotation)
	fmt.Printf("Field Width:      %.4f degrees\n", result.FieldWidth)
	fmt.Printf("Field Height:     %.4f degrees\n", result.FieldHeight)
	fmt.Printf("Solve Time:       %.2f seconds\n", result.SolveTime)
	fmt.Printf("\nOutput files:\n")
	for _, file := range result.OutputFiles {
		fmt.Printf("  - %s\n", file)
	}
}
