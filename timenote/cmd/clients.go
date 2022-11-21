package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jason0x43/go-toggl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
	"os"
	"text/tabwriter"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "clients management",
	Long:  `List clients and manage clients using sub commands`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		clients, err := p.Clients()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if viper.GetString("output-format") != "json" {
			writeClientsTable(filterClients(clients))
		} else {
			writeClientsJson(filterClients(clients))
		}
	},
}

func writeClientsJson(clients []toggl.Client) {
	data, err := json.Marshal(clients)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	_, _ = fmt.Println(string(data))
}

func writeClientsTable(clients []toggl.Client) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	_, _ = fmt.Fprintln(w, "ID\tName\t")
	for _, prj := range clients {
		_, _ = fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t", prj.ID, prj.Name))
	}
	_ = w.Flush()
}

func filterClients(clients []toggl.Client) []toggl.Client {
	return clients
}

func init() {
	RootCmd.AddCommand(clientsCmd)
}
