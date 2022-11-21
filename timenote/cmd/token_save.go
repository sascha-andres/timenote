package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
	"log"
)

// timestampAppendCmd represents the append command
var tokenSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "save token to local keyring",
	Run: func(cmd *cobra.Command, args []string) {
		t := viper.GetString("token.save.token")
		err := keyring.Set("timenote", "token", t)
		if err != nil {
			log.Printf("could not set token: %s", err)
		} else {
			log.Print("token set")
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenSaveCmd)

	tokenSaveCmd.Flags().StringP("token", "t", "", "toggl token to use")
	err := tokenSaveCmd.MarkFlagRequired("token")
	if err != nil {
		log.Fatalf("error marking as required flag: %s", err)
	}

	err = viper.BindPFlag("token.save.token", tokenSaveCmd.Flags().Lookup("token"))
	if err != nil {
		log.Fatalf("error adding flag: %s", err)
	}
}
