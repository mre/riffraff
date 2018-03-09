package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/bndr/gojenkins"
	"github.com/fatih/color"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// TODO: Replace this with a custom formatter or so
	statusCommand  = kingpin.Command("status", "Show the status of all matching jobs")
	statusRegexArg = statusCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()
	verbose        = statusCommand.Flag("verbose", "Verbose mode. Print full job output").Short('v').Bool()
	salt           = statusCommand.Flag("salt", "Show failed salt states").Bool()

	queueCommand = kingpin.Command("queue", "Show the queue of all matching jobs")
	nodesCommand = kingpin.Command("nodes", "Show the status of all Jenkins nodes")
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

func printStatus(waitGroup *sync.WaitGroup, jenkins *gojenkins.Jenkins, job gojenkins.InnerJob, verbose, salt bool) error {
	// Buffer full output to avoid race conditions between jobs
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
		if verbose || salt {
			output += fmt.Sprintf("Jenkins result code: %v\n", result)
			consoleOutput := lastBuild.GetConsoleOutput()
			if salt {
				output += fmt.Sprintf(consoleOutput)
			} else {
				for _, stateOutput := range getFailedSaltStates(consoleOutput) {
					output += stateOutput
				}
			}
			output += fmt.Sprintf("%v/consoleText\n", lastBuild.GetUrl())
		}
	}
	fmt.Print(output)
	return nil
}

// Find all jobs matching the given regex
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

func queue(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) {
	queue, err := jenkins.GetQueue()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Println(queue.Raw)
	// for _, task := range tasks {
	// 	fmt.Println(task.GetWhy())
	// }

}

func printNodeStatus(waitGroup *sync.WaitGroup, node gojenkins.Node) error {
	defer waitGroup.Done()
	// Fetch Node Data
	node.Poll()
	online, err := node.IsOnline()
	if err != nil {
		return err
	}
	if online {
		fmt.Printf("%v: Online\n", node.GetName())
	} else {
		fmt.Printf("%v: Offline\n", node.GetName())
	}
	return nil
}

func nodes(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) {
	nodes, err := jenkins.GetAllNodes()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(nodes))
	defer waitGroup.Wait()
	for _, node := range nodes {
		go printNodeStatus(&waitGroup, *node)
	}
}

func status(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) {
	jobs, err := findMatchingJobs(jenkins, regex)
	if err != nil {
		panic(fmt.Sprintf("Cannot retrieve jobs: %v", err))
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(jobs))
	defer waitGroup.Wait()
	for _, job := range jobs {
		go printStatus(&waitGroup, jenkins, job, verbose, salt)
	}
}

func main() {
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

	// TODO: Replace with a plugin-based system
	switch kingpin.Parse() {
	case "status":
		status(jenkins, *statusRegexArg, *verbose, *salt)
	case "queue":
		queue(jenkins, *statusRegexArg, *verbose, *salt)
	case "nodes":
		nodes(jenkins, *statusRegexArg, *verbose, *salt)
	default:
		kingpin.Usage()
	}
}
