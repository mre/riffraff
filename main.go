package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
)

func getFailedSaltStates(output string) []string {
	saltStates := strings.Split(output, "----------")
	var failedStates []string
	for _, state := range saltStates {
		if strings.Contains(state, "Result: False") {
			failedStates = append(failedStates, state)
		}
	}
	return failedStates
}

func runJob(waitGroup *sync.WaitGroup, jenkins *gojenkins.Jenkins, job gojenkins.InnerJob, verbose, raw bool) error {
	output := ""
	defer waitGroup.Done()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	build, err := jenkins.GetJob(job.Name)
	if err != nil {
		return err
	}

	lastBuild, err := build.GetLastBuild()
	var result string
	if err != nil {
		result = fmt.Sprintf("UNKNOWN (%v)", err)
	} else {
		result = lastBuild.GetResult()
	}

	var marker string
	switch result {
	case "SUCCESS":
		marker = green("✓")
	case "FAILURE":
		marker = red("✗")
	default:
		marker = yellow("?")
	}

	output += fmt.Sprintf("%v %v (%v)\n", marker, job.Name, job.Url)

	if result != "SUCCESS" && lastBuild != nil {
		if verbose || raw {
			output += fmt.Sprintf("Jenkins result code: %v\n", result)
			if raw {
				output += fmt.Sprintf(lastBuild.GetConsoleOutput())
			} else {
				for _, stateOutput := range getFailedSaltStates(lastBuild.GetConsoleOutput()) {
					output += stateOutput
				}
			}
			output += fmt.Sprintf("%v/consoleText\n", lastBuild.GetUrl())
		}
	}
	fmt.Print(output)
	return nil
}

func findMatchingJobs(jenkins *gojenkins.Jenkins, regex string) ([]gojenkins.InnerJob, error) {
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

func main() {

	verbose := flag.Bool("v", false, "verbose output")
	raw := flag.Bool("raw", false, "raw output")
	flag.Parse()

	var regex string
	if len(flag.Args()) > 0 {
		regex = flag.Args()[0]
	} else {
		regex = "*"
	}

	jenkinsURL := os.Getenv("JENKINS_URL")
	jenkinsUser := os.Getenv("JENKINS_USER")
	jenkinsPw := os.Getenv("JENKINS_PW")

	if len(jenkinsURL) == 0 {
		log.Fatal("Please set JENKINS_URL")
	}
	if len(jenkinsUser) == 0 {
		log.Fatal("Please set JENKINS_USER")
	}

	jenkins := gojenkins.CreateJenkins(nil, jenkinsURL, jenkinsUser, jenkinsPw)
	_, err := jenkins.Init()

	if err != nil {
		panic(fmt.Sprintf("Cannot authenticate: %v", err))
	}

	jobs, err := findMatchingJobs(jenkins, regex)
	if err != nil {
		panic(fmt.Sprintf("Cannot retrieve jobs: %v", err))
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(jobs))
	defer waitGroup.Wait()
	for _, job := range jobs {
		go runJob(&waitGroup, jenkins, job, *verbose, *raw)
	}
}
