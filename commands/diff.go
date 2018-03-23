package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/pmezard/go-difflib/difflib"
)

func DiffExec(jenkins *gojenkins.Jenkins, jobName string, build1, build2 int64) error {
	build, err := jenkins.GetJob(jobName)
	if err != nil {
		return err
	}

	build1Logs, err := build.GetBuild(build1)
	if err != nil {
		return err
	}

	build2Logs, err := build.GetBuild(build2)
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
