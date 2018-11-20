package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bndr/gojenkins"
	"github.com/mre/riffraff/internal/commands"
	"github.com/spf13/cobra"
)

var (
// statusCommand        = kingpin.Command("status", "Show the status of all matching jobs")
// statusRegexArg       = statusCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()
// statusOnlyFailingArg = statusCommand.Flag("only-failing", "Show only failing jobs").Bool()

// buildCommand  = kingpin.Command("build", "Trigger build for all matching jobs")
// buildRegexArg = buildCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

// logsCommand = kingpin.Command("logs", "Show the logs of a job")
// logsJobArg  = logsCommand.Arg("job", "The name of the job to get logs for").Required().String()

// diffCommand   = kingpin.Command("diff", "Print a diff between two builds of a job")
// diffJobArg    = diffCommand.Arg("job", "The name of the job to get the diff for").Required().String()
// diffBuild1Arg = diffCommand.Arg("build1", "First build").Required().Int64()
// diffBuild2Arg = diffCommand.Arg("build2", "Second build").Required().Int64()

// queueCommand  = kingpin.Command("queue", "Show the queue of all matching jobs")
// queueRegexArg = queueCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

// nodesCommand = kingpin.Command("nodes", "Show the status of all Jenkins nodes")

// openCommand  = kingpin.Command("open", "Open a job in the browser")
// openRegexArg = openCommand.Arg("regex", "The regular expression to match for the job names").Default(".*").String()

// verbose = kingpin.Flag("verbose", "Verbose mode. Print full job output").Short('v').Bool()

// // TODO: Replace this with a custom formatter or so
// salt = kingpin.Flag("salt", "Show failed salt states").Bool()
)
var jenkins *gojenkins.Jenkins
var statusCmd = &cobra.Command{
	Use:   "status <regex>",
	Short: "Show the status of all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		statusRegexArg := args[0]
		statusOnlyFailingArg, err := cmd.Flags().GetBool("only-failing")
		if err != nil {
			return err
		}
		fmt.Println(statusRegexArg, statusOnlyFailingArg)
		// return commands.NewStatus(jenkins, statusRegexArg, statusOnlyFailingArg).Exec()
		return nil
	},
}
var buildCmd = &cobra.Command{
	Use:   "build <regex>",
	Short: "Trigger build for all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		buildRegexArg := args[0]
		fmt.Println(buildRegexArg)
		// return commands.NewBuild(jenkins, buildRegexArg).Exec()
		return nil
	},
}
var logsCmd = &cobra.Command{
	Use:   "log <job>",
	Short: "Show the logs of a job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logsJobArg := args[0]
		salt, err := cmd.Flags().GetBool("salt")
		if err != nil {
			return err
		}
		fmt.Println(logsJobArg, salt)
		// return commands.NewLogs(jenkins, logsJobArg, salt).Exec()
		// return commands.NewStatus(jenkins, statusRegexArg, statusOnlyFailingArg).Exec()
		return nil
	},
}
var diffCmd = &cobra.Command{
	Use:   "diff <job> <build1> <build2>",
	Short: "Show the logs of a job",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		diffJobArg := args[0]
		diffBuild1Arg, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		diffBuild2Arg, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		fmt.Println(diffJobArg, diffBuild1Arg, diffBuild2Arg)
		// return commands.NewDiff(jenkins, diffJobArg, diffBuild1Arg, diffBuild2Arg).Exec()
		return nil
	},
}
var queueCmd = &cobra.Command{
	Use:   "queue <regex>",
	Short: "Show the queue of all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		queueRegexArg := args[0]
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		salt, err := cmd.Flags().GetBool("salt")
		if err != nil {
			return err
		}
		fmt.Println(queueRegexArg, verbose, salt)
		// return commands.NewQueue(jenkins, queueRegexArg, verbose, salt).Exec()
		return nil
	},
}
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Show the status of all Jenkins nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.NewNodes(jenkins).Exec()
	},
}
var openCmd = &cobra.Command{
	Use:   "open <regex>",
	Short: "Open a job in the browser",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		openRegexArg := args[0]
		fmt.Println(openRegexArg)
		// return commands.NewOpen(jenkins, openRegexArg).Exec()
		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "riffraff",
	Short: "riffraff is a commandline interface for Jenkins",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hello here")
		return nil
		// return create()
	},
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	fmt.Println("rootCmd run")
	// 	return nil
	// },
}

func create() error {
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
	if jenkins == nil {
		log.Fatal("Cannot instantiate Jenkins connection: null pointer return")
	}
	jenkins, err := jenkins.Init()

	if err != nil {
		log.Printf("Cannot authenticate: %v", err)
	}
	return err
}

func main() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(openCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode. Print full job output")
	rootCmd.PersistentFlags().Bool("salt", false, "Show failed salt states")

	statusCmd.Flags().Bool("only-failing", false, "Show only failing jobs")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	// TODO: Replace with a plugin-based system
	// e.g. https://github.com/mitchellh/cli
	// switch kingpin.Parse() {
	// case "status":
	// 	err = commands.NewStatus(jenkins, *statusRegexArg, *statusOnlyFailingArg).Exec()
	// case "diff":
	// 	err = commands.NewDiff(jenkins, *diffJobArg, *diffBuild1Arg, *diffBuild2Arg).Exec()
	// case "build":
	// 	err = commands.NewBuild(jenkins, *buildRegexArg).Exec()
	// case "logs":
	// 	err = commands.NewLogs(jenkins, *logsJobArg, *salt).Exec()
	// case "queue":
	// 	err = commands.NewQueue(jenkins, *queueRegexArg, *verbose, *salt).Exec()
	// case "nodes":
	// 	err = commands.NewNodes(jenkins).Exec()
	// case "open":
	// 	err = commands.NewOpen(jenkins, *openRegexArg).Exec()
	// default:
	// 	kingpin.Usage()
	// }

	// if err != nil {
	// 	log.Fatalf("Cannot execute command: %v", err)
	// }

}
