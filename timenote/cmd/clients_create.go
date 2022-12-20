package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
)

var clientsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a client",
	Long:  `Create a client in the current workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		name := viper.GetString("clients.create.name")

		err = p.CreateClient(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	clientsCmd.AddCommand(clientsCreateCmd)

	clientsCreateCmd.Flags().StringP("name", "", "", "name for client")
	_ = projectsCreateCmd.MarkFlagRequired("name")

	_ = viper.BindPFlag("clients.create.name", clientsCreateCmd.Flags().Lookup("name"))
}
