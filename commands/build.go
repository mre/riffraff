package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/mre/riffraff/job"
)

type Build struct {
	jenkins *gojenkins.Jenkins
	regex   string
}

func NewBuild(jenkins *gojenkins.Jenkins, regex string) *Build {
	return &Build{jenkins, regex}
}

func (b Build) Exec() error {
	jobs, err := job.FindMatchingJobs(b.jenkins, b.regex)
	if err != nil {
		return err
	}

	// TODO
	// var waitGroup sync.WaitGroup
	// waitGroup.Add(len(jobs))
	// defer waitGroup.Wait()
	for _, job := range jobs {
		id, err := b.jenkins.BuildJob(job.Name)
		if err != nil {
			return err
		}
		build, err := b.jenkins.GetBuild(job.Name, id)
		if err != nil {
			return err
		}
		fmt.Printf("Triggered build for %v [%v] %v\n", job.Name, id, build.GetUrl())
	}
	return nil
}
