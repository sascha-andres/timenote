package cmd

import (
	"github.com/spf13/cobra"
)

// browserCmd represents the browser command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "manage token (set token or delete token)",
}

func init() {
	RootCmd.AddCommand(tokenCmd)
}
