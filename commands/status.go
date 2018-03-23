package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
	"github.com/mre/riffraff/job"
)

func StatusExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := job.FindMatchingJobs(jenkins, regex)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(jobs))
	defer waitGroup.Wait()
	for _, job := range jobs {
		go printStatus(&waitGroup, jenkins, job)
	}
	return nil
}

func printStatus(waitGroup *sync.WaitGroup, jenkins *gojenkins.Jenkins, job gojenkins.InnerJob) error {
	defer waitGroup.Done()
	// Buffer full output to avoid race conditions between jobs
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	build, err := jenkins.GetJob(job.Name)
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

	fmt.Printf("%v %v (%v)\n", marker, job.Name, job.Url)
	return nil
}
