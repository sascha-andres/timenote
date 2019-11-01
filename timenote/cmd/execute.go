package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"strings"
	"time"

	"github.com/mgutz/str"
	"github.com/pkg/browser"
	"livingit.de/code/timenote/persistence"
)

func executeLine(p *persistence.TogglPersistor, commandline string) error {
	if strings.TrimSpace(commandline) == "" {
		return nil
	}
	tokenize := str.ToArgv(commandline)
	switch tokenize[0] {
	case "new":
		return p.New()
	case "done":
		return p.Done()
	case "current":
		entry, err := p.Current()
		if err != nil {
			return err
		}
		diff := time.Now().Sub(entry.Start)
		fmt.Println(humanizeDuration(diff))
		break
	case "append":
		return p.Append(strings.Join(tokenize[1:], " "), viper.GetString("separator"))
	case "project":
		return p.SetProjectForCurrentTimestamp(strings.Join(tokenize[1:], " "))
	case "tag":
		return p.Tag(strings.Join(tokenize[1:], " "))
	case "client":
		return executeClient(p, strings.Join(tokenize[1:], " "))
	case "open":
		return browser.OpenURL("https://toggl.com/app/timer")
	}
	return nil
}

// humanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default ``time.Duration`` output is badly formatted and unreadable.
func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}
