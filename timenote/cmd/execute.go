package cmd

import (
	"fmt"
	"strings"

	"github.com/mgutz/str"
	"github.com/sascha-andres/timenote/persistence"
)

func executeLine(persistence persistence.Persistor, commandline string) error {
	if strings.TrimSpace(commandline) == "" {
		return nil
	}
	tokenize := str.ToArgv(commandline)
	switch tokenize[0] {
	case "new":
		return persistence.New()
		break
	case "done":
		return persistence.Done()
		break
	case "current":
		entry, err := persistence.Current()
		if err != nil {
			return err
		}
		fmt.Println(entry)
		break
	case "append":
		return persistence.Append(strings.Join(tokenize[1:], " "))
		break
	case "project":
		return persistence.Project(strings.Join(tokenize[1:], " "))
		break
	case "tag":
		return persistence.Tag(strings.Join(tokenize[1:], " "))
		break
	}
	return nil
}
