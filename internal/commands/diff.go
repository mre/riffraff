package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/pmezard/go-difflib/difflib"
)

// Diff is a type that can be used to retrieve the differences between the console
// output of two different builds.
type Diff struct {
	jenkins *gojenkins.Jenkins
	jobName string
	build1  int64
	build2  int64
}

// NewDiff is a convenience method for configuring a new Diff instance with the
// provided job name and build numbers.
func NewDiff(jenkins *gojenkins.Jenkins, jobName string, build1, build2 int64) *Diff {
	return &Diff{jenkins, jobName, build1, build2}
}

// Exec will get the console output of two jobs and print the differences between
// the two console output strings.
func (d Diff) Exec() error {
	build, err := d.jenkins.GetJob(d.jobName)
	if err != nil {
		return err
	}

	build1Logs, err := build.GetBuild(d.build1)
	if err != nil {
		return err
	}

	build2Logs, err := build.GetBuild(d.build2)
	if err != nil {
		return err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(build1Logs.GetConsoleOutput()),
		B:        difflib.SplitLines(build2Logs.GetConsoleOutput()),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Printf(text)

	return nil
}
