package commands

import (
	// sec "github.com/ereztaiar/agenter/pkg/secrets"
	"github.com/spf13/cobra"
)

var secrets = &cobra.Command{
	Use:   "secrets",
	Short: "handles secrets",
	Long:  "Keeps secure version of the API Key",
	RunE: func(cmd *cobra.Command, args []string) error {
		// file, err := cmd.Flags().GetString("set")

		// if err != nil {
		// 	return err
		// }
		// sec.GetAPIKey()

		return nil
	},
}
