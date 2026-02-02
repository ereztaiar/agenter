package commands

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var runAgents = &cobra.Command{
	Use:   "run",
	Short: "Runs an agents file",
	Long:  "Runs agents file by path (default: ./agent.yaml)",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("file")

		if err != nil {
			return err
		}

		return nil
	},
	PersistentPostRun: loadAgents,
}

func loadAgents(cmd *cobra.Command, args []string) {
	file, err := cmd.Flags().GetString("file")

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
