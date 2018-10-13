package commands

import (
	"github.com/fatih/color"
)

const (
	// MarkerRunning status
	MarkerRunning = "↻"
	// MarkerSuccess status
	MarkerSuccess = "✓"
	// MarkerFailure status
	MarkerFailure = "✗"
	// MarkerDefault status
	MarkerDefault = "?"
	// StatusRunning text
	StatusRunning = "RUNNING"
	// StatusSuccess text
	StatusSuccess = "SUCCESS"
	// StatusFailure text
	StatusFailure = "FAILURE"
)

// GetMarker for status
func GetMarker(status string) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch status {
	case "RUNNING":
		return green(MarkerRunning)
	case "SUCCESS":
		return green(MarkerSuccess)
	case "FAILURE":
		return red(MarkerFailure)
	default:
		return yellow(MarkerDefault)
	}
}
