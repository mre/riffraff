package commands

import (
	"log"

	"github.com/mre/riffraff/internal/job"

	"github.com/bndr/gojenkins"
	"github.com/skratchdot/open-golang/open"
)

// Open is a type that can be used to open a job or set of jobs in a browser.
type Open struct {
	jenkins *gojenkins.Jenkins
	regex   string
}

// NewOpen is a convenience method for initializing a new Open instance.
func NewOpen(jenkins *gojenkins.Jenkins, regex string) *Open {
	return &Open{
		jenkins,
		regex,
	}
}

// Exec will find a list of jobs with a name matching the specified regular expression and will
// open the page for that job in the system's default web browser.
func (o Open) Exec() error {
	jobs, err := job.FindMatchingJobs(o.jenkins, o.regex)
	if err != nil {
		return err
	}
	if len(jobs) > 3 {
		log.Fatalf("More than three jobs match your criteria. This is probably not what you expected. Please narrow down your search\n")
	}

	for _, job := range jobs {
		if err = open.Run(job.Url); err != nil {
			return err
		}
	}
	return nil
}
