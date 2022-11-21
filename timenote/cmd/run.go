package cmd

import (
	"io"
	"strings"

	"log"

	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
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
	p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
	if err != nil {
		return errors.Wrap(err, "Could not create p layer")
	}

	return runInputLoop(p)
}

func runInputLoop(p *persistence.TogglPersistor) error {
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
			err := executeLine(p, line)
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
