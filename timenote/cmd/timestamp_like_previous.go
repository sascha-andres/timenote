package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
)

// timestampAppendCmd represents the append command
var timestampLikePreviousCmd = &cobra.Command{
	Use:   "like-previous",
	Short: "Start the last finished entry again",
	Long: `Select the last finished entry and use it as a template
to start a new one`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		err = p.StartPrevious()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampLikePreviousCmd)
}
