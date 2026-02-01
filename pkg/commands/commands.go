package commands

import (
	"github.com/ereztaiar/agenter/pkg/agents"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var (
	agentFile  string
	rootConfig *agent.RootConfig
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}

	RootCmd.PersistentFlags().StringVar(&agentFile, "file", "", "agent file (default is ./agent.yaml)")

	RootCmd.AddCommand(listCmd)
	viper.SetDefault("file", filepath.Join(os.Getenv("PWD"), "agent.yaml"))
}

var RootCmd = &cobra.Command{
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

	log.Println("Loading file", file)

	data, err := os.ReadFile(file)

	replaced := os.ExpandEnv(string(data))

	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal([]byte(replaced), &rootConfig)
	if err != nil {
		log.Fatalln(err)
	}

	rootConfig.BuildAgents()

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

func Start(rc *agent.RootConfig) {
	rootConfig = rc
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)

	}
}
