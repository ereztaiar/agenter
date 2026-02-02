package commands

import (
	"github.com/spf13/cobra"
	"log"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running agents",
	Long:  "List agents running under current session",
	Run:   listAgents,
}

func listAgents(cmd *cobra.Command, args []string) {
	log.Println("Command reserved for daemon")

}
