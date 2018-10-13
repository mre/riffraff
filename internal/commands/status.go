package commands

import (
	"fmt"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
	"github.com/mre/riffraff/internal/job"
)

type Status struct {
	jenkins *gojenkins.Jenkins
	regex   string
}

func NewStatus(jenkins *gojenkins.Jenkins, regex string) *Status {
	return &Status{jenkins, regex}
}

func (s Status) Exec() error {
	jobs, err := job.FindMatchingJobs(s.jenkins, s.regex)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job gojenkins.InnerJob) {
			defer wg.Done()
			s.print(job)
		}(job)
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
	var result string
	if err != nil {
		result = fmt.Sprintf("UNKNOWN (%v)", err)
	} else {
		if lastBuild.IsRunning() {
			result = StatusRunning
		} else {
			result = lastBuild.GetResult()
		}
	}

	marker := GetMarker(result)

	fmt.Printf("%v %v (%v)\n", marker, job.Name, job.Url)
	return nil
}
