package commands

import (
	"log"

	"github.com/mre/riffraff/job"

	"github.com/bndr/gojenkins"
	"github.com/skratchdot/open-golang/open"
)

type Open struct {
	jenkins *gojenkins.Jenkins
	regex   string
}

func NewOpen(jenkins *gojenkins.Jenkins, regex string) *Open {
	return &Open{}
}

func (o Open) Exec() error {
	jobs, err := job.FindMatchingJobs(o.jenkins, o.regex)
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
