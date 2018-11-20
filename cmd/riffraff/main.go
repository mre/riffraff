package main

import (
	"log"
	"os"
	"strconv"

	"github.com/bndr/gojenkins"
	"github.com/mre/riffraff/internal/commands"
	"github.com/spf13/cobra"
)

var jenkins *gojenkins.Jenkins
var statusCmd = &cobra.Command{
	Use:   "status [REGEX]",
	Short: "Show the status of all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var statusRegexArg = ".*"
		if len(args) > 0 {
			statusRegexArg = args[0]
		}
		statusOnlyFailingArg, err := cmd.Flags().GetBool("only-failing")
		if err != nil {
			return err
		}
		return commands.NewStatus(jenkins, statusRegexArg, statusOnlyFailingArg).Exec()
	},
}
var buildCmd = &cobra.Command{
	Use:   "build [<regex>]",
	Short: "Trigger build for all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var buildRegexArg = ".*"
		if len(args) > 0 {
			buildRegexArg = args[0]
		}
		return commands.NewBuild(jenkins, buildRegexArg).Exec()
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
		return commands.NewLogs(jenkins, logsJobArg, salt).Exec()
	},
}
var diffCmd = &cobra.Command{
	Use:   "diff <job> <build1> <build2>",
	Short: "Print a diff between two builds of a job",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		diffJobArg := args[0]
		diffBuild1Arg, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		diffBuild2Arg, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}
		return commands.NewDiff(jenkins, diffJobArg, diffBuild1Arg, diffBuild2Arg).Exec()
	},
}
var queueCmd = &cobra.Command{
	Use:   "queue [<regex>]",
	Short: "Show the queue of all matching jobs",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var queueRegexArg = ".*"
		if len(args) > 0 {
			queueRegexArg = args[0]
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		salt, err := cmd.Flags().GetBool("salt")
		if err != nil {
			return err
		}
		return commands.NewQueue(jenkins, queueRegexArg, verbose, salt).Exec()
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
	Use:   "open [<regex>]",
	Short: "Open a job in the browser",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var openRegexArg = ".*"
		if len(args) > 0 {
			openRegexArg = args[0]
		}
		return commands.NewOpen(jenkins, openRegexArg).Exec()
	},
}

var rootCmd = &cobra.Command{
	Use:   "riffraff",
	Short: "riffraff is a commandline interface for Jenkins",
}

func main() {
	authenticate()
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(openCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode. Print full job output")

	// TODO: Replace this with a custom formatter or so
	rootCmd.PersistentFlags().Bool("salt", false, "Show failed salt states")

	statusCmd.Flags().Bool("only-failing", false, "Show only failing jobs")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Cannot execute command: %v", err)
	}
}

func authenticate() {
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
}
