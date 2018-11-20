package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/mre/riffraff/internal/job"
)

// Build can be used to kick of new builds of jobs with a name matching
// the provided regular expression.
type Build struct {
	jenkins *gojenkins.Jenkins
	regex   string
}

// NewBuild is a convenience method for initializing a new Build instance.
func NewBuild(jenkins *gojenkins.Jenkins, regex string) *Build {
	return &Build{jenkins, regex}
}

// Exec will send a request to the Jenkins server to start new builds of all jobs
// with names that match the configured regular expression.
func (b Build) Exec() error {
	jobs, err := job.FindMatchingJobs(b.jenkins, b.regex)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job gojenkins.InnerJob) {
			defer wg.Done()

			id, err := b.jenkins.BuildJob(job.Name)
			if err != nil {
				fmt.Printf("Triggering build for %v failed: %v\n", job.Name, err)
				return
			}
			build, err := b.jenkins.GetBuild(job.Name, id)
			if err != nil {
				fmt.Printf("Getting build for %v [%v] failed: %v\n", job.Name, id, err)
				return
			}
			fmt.Printf("Triggered build for %v [%v] %v\n", job.Name, id, build.GetUrl())
		}(job)
	}
	wg.Wait()
	return nil
}
