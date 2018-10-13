package commands

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMarker(t *testing.T) {
	Convey("Given a status", t, func() {
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
			description := fmt.Sprintf("expected %s for %s", tt.marker, tt.status)

			Convey(description, func() {
				So(got, ShouldEqual, want)
			})
		}
	})
}
