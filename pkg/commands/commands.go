package commands

import (
	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
)

var (
	rootConfig *agent.RootConfig
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}

	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(versionCommand)
	runAgents.
		Flags().
		StringP(
			"file",
			"f",
			"agent.yaml",
			"agent file to run")
	RootCmd.AddCommand(runAgents)
	RootCmd.AddCommand(secrets)
}

var RootCmd = &cobra.Command{
	Use:   "agenter",
	Short: "An AI agent orchestrator",
	Long:  `agenter lets you run agents localy and quickly.`,
}

func Start(rc *agent.RootConfig) {
	rootConfig = rc
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)

	}
}
