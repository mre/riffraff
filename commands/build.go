package commands

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/mre/riffraff/job"
)

func BuildExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := job.FindMatchingJobs(jenkins, regex)
	if err != nil {
		return err
	}

	// TODO
	// var waitGroup sync.WaitGroup
	// waitGroup.Add(len(jobs))
	// defer waitGroup.Wait()
	for _, job := range jobs {
		id, err := jenkins.BuildJob(job.Name)
		if err != nil {
			return err
		}
		build, err := jenkins.GetBuild(job.Name, id)
		if err != nil {
			return err
		}
		fmt.Printf("Triggered build for %v [%v] %v\n", job.Name, id, build.GetUrl())
	}
	return nil
}
