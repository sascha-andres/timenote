package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"log"
)

// tokenDeleteCmd represents the current command
var tokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete token from keyring",
	Long:  `Delete the token from the local keyring`,
	Run: func(cmd *cobra.Command, args []string) {
		err := keyring.Delete("timenote", "token")
		if err != nil {
			log.Printf("error deleting token: %s", err)
		} else {
			log.Print("token successfully deleted")
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenDeleteCmd)
}
