package main

import (
	"fmt"
	"log"
	"os"

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
	Use:   "status",
	Short: "Show the status of all matching jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		statusRegexArg, _ := cmd.Flags().GetString("regex")
		statusOnlyFailingArg, _ := cmd.Flags().GetBool("only-failing")
		fmt.Println(statusRegexArg, statusOnlyFailingArg)
		// return commands.NewStatus(jenkins, statusRegexArg, statusOnlyFailingArg).Exec()
		return nil
	},
}
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Trigger build for all matching jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		buildRegexArg, _ := cmd.Flags().GetString("regex")
		fmt.Println(buildRegexArg)
		// return commands.NewBuild(jenkins, buildRegexArg).Exec()
		return nil
	},
}
var logsCmd = &cobra.Command{
	Use:   "log",
	Short: "Show the logs of a job",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.MarkFlagRequired("job")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logsJobArg, _ := cmd.Flags().GetString("job")
		fmt.Println(logsJobArg)
		// return commands.NewStatus(jenkins, statusRegexArg, statusOnlyFailingArg).Exec()
		return nil
	},
}
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show the logs of a job",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check error
		cmd.MarkFlagRequired("build1")
		cmd.MarkFlagRequired("build2")
		cmd.MarkFlagRequired("job")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// 	err = commands.NewDiff(jenkins, *diffJobArg, *diffBuild1Arg, *diffBuild2Arg).Exec()
		diffJobArg, _ := cmd.Flags().GetString("job")
		diffBuild1Arg, _ := cmd.Flags().GetInt64("build1")
		diffBuild2Arg, _ := cmd.Flags().GetInt64("build2")

		fmt.Println(diffJobArg, diffBuild1Arg, diffBuild2Arg)
		// return commands.NewDiff(jenkins, diffJobArg, diffBuild1Arg, diffBuild2Arg).Exec()
		return nil
	},
}
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Show the queue of all matching jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		queueRegexArg, _ := cmd.Flags().GetString("regex")
		fmt.Println(queueRegexArg)
		// return commands.NewQueue(jenkins, queueRegexArg).Exec()
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
	Use:   "open",
	Short: "Open a job in the browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		openRegexArg, _ := cmd.Flags().GetString("regex")
		fmt.Println(openRegexArg)
		// return commands.NewOpen(jenkins, openRegexArg).Exec()
		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "riffraff",
	Short: "riffraff is a commandline interface for Jenkins",
	Long:  "riffraff long description",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hello here")
		return nil
		// return create()
	},
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

	statusCmd.Flags().String("regex", ".*", "The regular expression to match for the job names")
	statusCmd.Flags().Bool("only-failing", false, "Show only failing jobs")
	buildCmd.Flags().String("regex", ".*", "The regular expression to match for the job names")

	// unsure about the initial values
	logsCmd.Flags().String("job", "", "The name of the job to get logs for")
	diffCmd.Flags().String("job", "", "The name of the job to get the diff for")
	diffCmd.Flags().Int64("build1", 0, "First build")
	diffCmd.Flags().Int64("build2", 0, "Second build")

	queueCmd.Flags().String("regex", ".*", "The regular expression to match for the job names")
	openCmd.Flags().String("regex", ".*", "The regular expression to match for the job names")
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
