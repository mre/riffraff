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
	"github.com/pmezard/go-difflib/difflib"
	"github.com/skratchdot/open-golang/open"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	statusCommand  = kingpin.Command("status", "Show the status of all matching jobs")
	statusRegexArg = statusCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

	buildCommand  = kingpin.Command("build", "Trigger build for all matching jobs")
	buildRegexArg = buildCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

	logsCommand = kingpin.Command("logs", "Show the logs of a job")
	logsJobArg  = logsCommand.Arg("job", "The name of the job to get logs for").Required().String()

	diffCommand   = kingpin.Command("diff", "Print a diff between two builds of a job")
	diffJobArg    = diffCommand.Arg("job", "The name of the job to get the diff for").Required().String()
	diffBuild1Arg = diffCommand.Arg("build1", "First build").Required().Int64()
	diffBuild2Arg = diffCommand.Arg("build2", "Second build").Required().Int64()

	queueCommand  = kingpin.Command("queue", "Show the queue of all matching jobs")
	queueRegexArg = queueCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

	nodesCommand = kingpin.Command("nodes", "Show the status of all Jenkins nodes")

	openCommand  = kingpin.Command("open", "Open a job in the browser")
	openRegexArg = openCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

	verbose = kingpin.Flag("verbose", "Verbose mode. Print full job output").Short('v').Bool()

	// TODO: Replace this with a custom formatter or so
	salt = kingpin.Flag("salt", "Show failed salt states").Bool()
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

func printStatus(waitGroup *sync.WaitGroup, jenkins *gojenkins.Jenkins, job gojenkins.InnerJob) error {
	defer waitGroup.Done()
	// Buffer full output to avoid race conditions between jobs
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

	fmt.Printf("%v %v (%v)\n", marker, job.Name, job.Url)
	return nil
}

func logsExec(jenkins *gojenkins.Jenkins, jobName string, salt bool) error {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	build, err := jenkins.GetJob(jobName)
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

	fmt.Printf("%v %v (%v)\n", marker, jobName, lastBuild.GetUrl())

	fmt.Printf("Jenkins result code: %v\n", result)
	consoleOutput := lastBuild.GetConsoleOutput()
	if salt {
		for _, stateOutput := range getFailedSaltStates(consoleOutput) {
			fmt.Println(stateOutput)
		}
	} else {
		fmt.Printf(consoleOutput)
	}
	fmt.Printf("%v/consoleText\n", lastBuild.GetUrl())
	return nil
}

func diffExec(jenkins *gojenkins.Jenkins, jobName string, build1, build2 int64) error {
	build, err := jenkins.GetJob(jobName)
	if err != nil {
		return err
	}

	build1Logs, err := build.GetBuild(build1)
	if err != nil {
		return err
	}

	build2Logs, err := build.GetBuild(build2)
	if err != nil {
		return err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(build1Logs.GetConsoleOutput()),
		B:        difflib.SplitLines(build2Logs.GetConsoleOutput()),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Printf(text)

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

func queueExec(jenkins *gojenkins.Jenkins, regex string, verbose, salt bool) error {
	queue, err := jenkins.GetQueue()
	if err != nil {
		return err
	}
	fmt.Println(queue.Raw)
	// for _, task := range tasks {
	// 	fmt.Println(task.GetWhy())
	// }
	return nil
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

func nodesExec(jenkins *gojenkins.Jenkins) error {
	nodes, err := jenkins.GetAllNodes()
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(nodes))
	defer waitGroup.Wait()
	for _, node := range nodes {
		go printNodeStatus(&waitGroup, *node)
	}
	return nil
}

func openExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := findMatchingJobs(jenkins, regex)
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

func statusExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := findMatchingJobs(jenkins, regex)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(jobs))
	defer waitGroup.Wait()
	for _, job := range jobs {
		go printStatus(&waitGroup, jenkins, job)
	}
	return nil
}

func buildExec(jenkins *gojenkins.Jenkins, regex string) error {
	jobs, err := findMatchingJobs(jenkins, regex)
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
		err = statusExec(jenkins, *statusRegexArg)
	case "diff":
		err = diffExec(jenkins, *diffJobArg, *diffBuild1Arg, *diffBuild2Arg)
	case "build":
		err = buildExec(jenkins, *buildRegexArg)
	case "logs":
		err = logsExec(jenkins, *logsJobArg, *salt)
	case "queue":
		err = queueExec(jenkins, *queueRegexArg, *verbose, *salt)
	case "nodes":
		err = nodesExec(jenkins)
	case "open":
		err = openExec(jenkins, *openRegexArg)
	default:
		kingpin.Usage()
	}

	if err != nil {
		log.Fatalf("Cannot execute command: %v", err)
	}
}
