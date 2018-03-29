package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bndr/gojenkins"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/mre/riffraff/commands"
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
		err = commands.NewStatus(jenkins, *statusRegexArg).Exec()
	case "diff":
		err = commands.NewDiff(jenkins, *diffJobArg, *diffBuild1Arg, *diffBuild2Arg).Exec()
	case "build":
		err = commands.NewBuild(jenkins, *buildRegexArg).Exec()
	case "logs":
		err = commands.NewLogs(jenkins, *logsJobArg, *salt).Exec()
	case "queue":
		err = commands.NewQueue(jenkins, *queueRegexArg, *verbose, *salt).Exec()
	case "nodes":
		err = commands.NewNodes(jenkins).Exec()
	case "open":
		err = commands.NewOpen(jenkins, *openRegexArg).Exec()
	default:
		kingpin.Usage()
	}

	if err != nil {
		log.Fatalf("Cannot execute command: %v", err)
	}
}
