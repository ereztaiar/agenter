package commands

import (
	agenterVersion "github.com/ereztaiar/agenter/pkg/version"
	"github.com/spf13/cobra"
	"log"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Get the current version",
	Long:  "The current build of agenter app",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(agenterVersion.Version)
	},
}
