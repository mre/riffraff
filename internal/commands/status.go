package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
	"github.com/mre/riffraff/internal/job"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

var (
	Unknown = yellow("?")
	Running = green("↻")
	Good    = green("✓")
	Bad     = red("✗")
)

type Status struct {
	jenkins     *gojenkins.Jenkins
	regex       string
	onlyFailing bool
}

func NewStatus(jenkins *gojenkins.Jenkins, regex string, onlyFailing bool) *Status {
	return &Status{jenkins, regex, onlyFailing}
}

func (s Status) Exec() error {
	jobs, err := job.FindMatchingJobs(s.jenkins, s.regex)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := range jobs {
		wg.Add(1)
		go func(job gojenkins.InnerJob) {
			defer wg.Done()
			s.print(job)
		}(jobs[i])
	}
	wg.Wait()
	return nil
}

func (s Status) print(job gojenkins.InnerJob) error {
	// Buffer full output to avoid race conditions between jobs

	build, err := s.jenkins.GetJob(job.Name)
	if err != nil {
		return err
	}

	lastBuild, err := build.GetLastBuild()
	marker := Unknown
	if err == nil {
		result := lastBuild.GetResult()
		if lastBuild.IsRunning() {
			marker = Running
		} else if result == "SUCCESS" {
			marker = Good
		} else if result == "FAILURE" {
			marker = Bad
		}
	}

	if s.onlyFailing && marker != Bad {
		return nil
	}

	fmt.Printf("%v %v (%v)\n", marker, job.Name, job.Url)
	return nil
}
