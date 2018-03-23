package commands

import (
	"log"

	"github.com/mre/riffraff/job"

	"github.com/bndr/gojenkins"
	"github.com/skratchdot/open-golang/open"
)

func OpenExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := job.FindMatchingJobs(jenkins, regex)
	if err != nil {
		return err
	}
	if len(jobs) > 3 {
		log.Fatalf("More than three jobs match your criteria. This is probably not what you expected. Please narrow down your search\n")
	}

	for _, job := range jobs {
		open.Run(job.Url)
	}
	return nil
}
