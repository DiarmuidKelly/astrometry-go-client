// Package client provides a unified Go client for astrometry.net operations
package client

import "errors"

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
