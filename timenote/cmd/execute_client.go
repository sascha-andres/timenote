package cmd

import (
	"fmt"
	"github.com/mgutz/str"
	"livingit.de/code/timenote/internal/persistence"
	"strings"
)

func executeClient(p *persistence.TogglPersistor, commandline string) error {
	if commandline == "" {
		clients, err := p.Clients()
		if err != nil {
			return err
		}
		for _, client := range clients {
			fmt.Printf("%s", client.Name)
		}
	}

	if strings.HasPrefix(commandline, "new ") {
		tokenize := str.ToArgv(commandline)
		clientName := strings.Join(tokenize[1:], " ")
		return p.NewClient(clientName)
	}

	return nil
}
