package cmd

import (
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"livingit.de/code/timenote/persistence"
	"livingit.de/code/timenote/persistence/factory"
	"github.com/spf13/viper"
)

var (
	completer = readline.NewPrefixCompleter(
		readline.PcItem("new"),
		readline.PcItem("append"),
		readline.PcItem("tag"),
		readline.PcItem("done"),
		readline.PcItem("current"),
	)
)

func run() error {
	persistence, err := factory.CreatePersistence(viper.GetString("persistor"), viper.GetString("dsn"))
	if err != nil {
		return errors.Wrap(err, "Could not create persistence layer")
	}
	defer func() {
		err := persistence.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return runInputLoop(persistence)
}

func runInputLoop(persistence persistence.Persistor) error {
	l, err := getReadlineConfig()
	if err != nil {
		return errors.Wrap(err, "Error creating readline config")
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Fatalf("Error closing readline: " + err.Error())
		}
	}()

	for {
		line, doBreak := getLine(l)
		if doBreak {
			break
		}
		line = strings.TrimSpace(line)
		switch line {
		case "quit", "q":
			return nil
		default:
			err := executeLine(persistence, line)
			if err != nil {
				log.Printf("Error: %#v\n", err)
			}
			break
		}
	}

	return nil
}

func getLine(l *readline.Instance) (string, bool) {
	line, err := l.Readline()
	if err == readline.ErrInterrupt {
		if len(line) == 0 {
			return "", true
		}
		return line, false
	} else if err == io.EOF {
		return "", true
	}
	return line, false
}

func getReadlineConfig() (*readline.Instance, error) {
	return readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/timenote.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
}
