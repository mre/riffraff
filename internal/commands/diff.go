package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/pmezard/go-difflib/difflib"
)

type Diff struct {
	jenkins *gojenkins.Jenkins
	jobName string
	build1  int64
	build2  int64
}

func NewDiff(jenkins *gojenkins.Jenkins, jobName string, build1, build2 int64) *Diff {
	return &Diff{jenkins, jobName, build1, build2}
}

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
