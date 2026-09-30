package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jason0x43/go-toggl"
)

// dispatchClients routes "clients" sub-verbs; with none given it lists
// clients.
func dispatchClients(verbs []string) {
	if len(verbs) == 0 {
		clientsList()
		return
	}
	switch verbs[0] {
	case "create":
		clientsCreate()
	default:
		fatal(fmt.Errorf("unknown clients command %q", verbs[0]))
	}
}

func clientsList() {
	p := newPersistor()
	clients, err := p.Clients()
	if err != nil {
		fatal(err)
	}

	if *outputFormat != "json" {
		writeClientsTable(clients)
	} else {
		writeClientsJson(clients)
	}
}

func writeClientsJson(clients []toggl.Client) {
	data, err := json.Marshal(clients)
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(data))
}

func writeClientsTable(clients []toggl.Client) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "ID\tName\t")
	for _, c := range clients {
		fmt.Fprintf(w, "%d\t%s\t\n", c.ID, c.Name)
	}
	_ = w.Flush()
}
