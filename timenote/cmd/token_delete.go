package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log/slog"
)

// tokenDeleteCmd represents the current command
var tokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete token from keyring",
	Long:  `Delete the token from the local keyring`,
	Run: func(cmd *cobra.Command, args []string) {
		err := keyring.Delete("timenote", "token")
		if err != nil {
			slog.Error("error deleting token", "error", err)
		} else {
			slog.Info("token successfully deleted")
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenDeleteCmd)
}
