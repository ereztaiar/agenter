package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"agenter/pkg/agent"
)

var (
	agentFile  string
	rootConfig *RootConfig
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	rootCmd.PersistentFlags().StringVar(&agentFile, "file", "", "agent file (default is ./agent.yaml)")

	rootCmd.AddCommand(listCmd)
	viper.SetDefault("file", filepath.Join(os.Getenv("PWD"), "agent.yaml"))
}

var rootCmd = &cobra.Command{
	Use:               "agenter",
	Short:             "An AI agent orchestrator",
	Long:              `agenter lets you run agents localy and quickly.`,
	Run:               rootTask,
	PersistentPreRun:  loadAgents,
	PersistentPostRun: clearAgents,
}

func getAgentsFile() string {
	if agentFile != "" {
		return agentFile
	}
	return viper.GetString("file")
}

func loadAgents(cmd *cobra.Command, args []string) {

	file := getAgentsFile()

	data, err := os.ReadFile(file)

	replaced := os.ExpandEnv(string(data))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(replaced), &rootConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ac := range rootConfig.AgentsConfig {
		if ac.Type == "root-agent" {
			continue
		}
		// fmt.Println(key)
		ac.GenerateAgent()
	}

}

func clearAgents(cmd *cobra.Command, args []string) {
	// fmt.Println("agents has been relesed")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running agents",
	Long:  "List agents running under current session",
	Run:   listAgents,
}

func rootTask(cmd *cobra.Command, args []string) {
	// fmt.Println("root task")
}

func listAgents(cmd *cobra.Command, args []string) {
	// fmt.Println("running agents list")
}
