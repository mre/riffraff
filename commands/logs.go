package commands

import (
	"fmt"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
)

type Logs struct {
	jenkins *gojenkins.Jenkins
	jobName string
	salt    bool
}

func NewLogs(jenkins *gojenkins.Jenkins, jobName string, salt bool) *Logs {
	return &Logs{jenkins, jobName, salt}
}

func (l Logs) Exec() error {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	build, err := l.jenkins.GetJob(l.jobName)
	if err != nil {
		return err
	}

	lastBuild, err := build.GetLastBuild()
	var result string
	if err != nil {
		result = fmt.Sprintf("UNKNOWN (%v)", err)
	} else {
		result = lastBuild.GetResult()
	}

	var marker string
	switch result {
	case "SUCCESS":
		marker = green("✓")
	case "FAILURE":
		marker = red("✗")
	default:
		marker = yellow("?")
	}

	fmt.Printf("%v %v (%v)\n", marker, l.jobName, lastBuild.GetUrl())

	fmt.Printf("Jenkins result code: %v\n", result)
	consoleOutput := lastBuild.GetConsoleOutput()
	if l.salt {
		for _, stateOutput := range getFailedSaltStates(consoleOutput) {
			fmt.Println(stateOutput)
		}
	} else {
		fmt.Printf(consoleOutput)
	}
	fmt.Printf("%v/consoleText\n", lastBuild.GetUrl())
	return nil
}

func getFailedSaltStates(output string) []string {
	saltStates := strings.Split(output, "----------")
	var failedStates []string
	for _, state := range saltStates {
		if strings.Contains(state, "Result: False") {
			failedStates = append(failedStates, state)
		}
	}
	return failedStates
}
