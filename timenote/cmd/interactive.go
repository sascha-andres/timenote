package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// timestampCmd represents the timestamp command
var interactiveCmd = &cobra.Command{
	Use:   "i",
	Short: "interactive shell",
	Long: `Interact with toggl interactively with a "shell"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(interactiveCmd)
}