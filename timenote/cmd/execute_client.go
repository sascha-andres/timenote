package cmd

import (
	"fmt"
	"github.com/mgutz/str"
	"livingit.de/code/timenote/persistence"
	"strings"
)

func executeClient(persistence persistence.Persistor, commandline string) error {
	fmt.Println(commandline)
	if commandline == "" {
		clients, err := persistence.Client()
		if err != nil {
			return err
		}
		for _, client := range clients {
			fmt.Printf("%s", client)
		}
	}

	if strings.HasPrefix(commandline, "new ") {
		tokenize := str.ToArgv(commandline)
		clientName := strings.Join(tokenize[1:], " ")
	}
	return nil
}
