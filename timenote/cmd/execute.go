package cmd

import (
	"fmt"
	"strings"

	"github.com/mgutz/str"
	"github.com/pkg/browser"
	"livingit.de/code/timenote/persistence"
)

func executeLine(persistence persistence.Persistor, commandline string) error {
	if strings.TrimSpace(commandline) == "" {
		return nil
	}
	tokenize := str.ToArgv(commandline)
	switch tokenize[0] {
	case "new":
		return persistence.New()
	case "done":
		return persistence.Done()
	case "current":
		entry, err := persistence.Current()
		if err != nil {
			return err
		}
		fmt.Println(entry)
		break
	case "append":
		return persistence.Append(strings.Join(tokenize[1:], " "))
	case "project":
		return persistence.Project(strings.Join(tokenize[1:], " "))
	case "tag":
		return persistence.Tag(strings.Join(tokenize[1:], " "))
	case "client":
		return persistence.Client(strings.Join(tokenize[1:], " "))
	case "open":
		hasOne, url, err := persistence.GetWebsite()
		if err != nil {
			return err
		}
		if hasOne {
			return browser.OpenURL(url)
		}
		break
	}
	return nil
}
