package job

import (
	"regexp"

	"github.com/bndr/gojenkins"
)

// FindMatchingJobs finds all jobs matching the given regex
func FindMatchingJobs(jenkins *gojenkins.Jenkins, regex string) ([]gojenkins.InnerJob, error) {
	jobs, err := jenkins.GetAllJobNames()
	if err != nil {
		return nil, err
	}

	var matchingJobs []gojenkins.InnerJob
	for _, job := range jobs {
		match, _ := regexp.MatchString(regex, job.Name)
		if match {
			matchingJobs = append(matchingJobs, job)
		}
	}

	return matchingJobs, nil
}
