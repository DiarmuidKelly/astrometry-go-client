package client

// This file re-exports solver types for public API

import (
	"github.com/DiarmuidKelly/astrometry-go-client/internal/solver"
)

// SolveOptions holds parameters for a plate-solving operation.
type SolveOptions = solver.SolveOptions

// Result holds the plate-solving results.
type Result = solver.Result

// DefaultSolveOptions returns SolveOptions with sensible defaults.
func DefaultSolveOptions() *SolveOptions {
	return solver.DefaultSolveOptions()
}
