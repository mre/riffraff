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
	// Unknown represents the default or unknown status if an error occurs.
	Unknown = yellow("?")
	// Running is use to indicate jobs that are currently running.
	Running = green("↻")
	// Good is a generic status that represents a healthy job/node.
	Good = green("✓")
	// Bad should be used to represent a problem with a job/node.
	Bad = red("✗")
)

// Status is a type that is used to get the current build status of a Jenkins job.
// Including the regex that is used to match the desired job names.
type Status struct {
	jenkins     *gojenkins.Jenkins
	regex       string
	onlyFailing bool
}

func NewStatus(jenkins *gojenkins.Jenkins, regex string, onlyFailing bool) *Status {
	return &Status{jenkins, regex, onlyFailing}
}

// Exec is responsible for finding all of the matching jobs on the Jenkins server
// and prints the status of each job that matches the provided regular expression.
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
