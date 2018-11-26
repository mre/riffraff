package commands

import (
	"testing"

	"github.com/fatih/color"
)

func TestGetMarker(t *testing.T) {
	statusTests := []struct {
		status string
		marker string
		color  color.Attribute
	}{
		{status: StatusRunning, marker: MarkerRunning, color: color.FgGreen},
		{status: StatusSuccess, marker: MarkerSuccess, color: color.FgGreen},
		{status: StatusFailure, marker: MarkerFailure, color: color.FgRed},
		{status: "", marker: MarkerDefault, color: color.FgYellow},
		{status: "aone1231", marker: MarkerDefault, color: color.FgYellow},
		{status: "UNKNOWN", marker: MarkerDefault, color: color.FgYellow},
	}

	for _, tt := range statusTests {

		status := tt.status
		c := color.New(tt.color).SprintFunc()

		want := c(tt.marker)
		got := GetMarker(status)

		if got != want {
			t.Errorf("expected %v, want %v", got, want)
		}
	}
}
