package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

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
		diff := time.Now().Sub(entry.Start)
		fmt.Println(humanizeDuration(diff))
		break
	case "append":
		return persistence.Append(strings.Join(tokenize[1:], " "))
	case "project":
		return persistence.Project(strings.Join(tokenize[1:], " "))
	case "tag":
		return persistence.Tag(strings.Join(tokenize[1:], " "))
	case "client":
		return executeClient(persistence, strings.Join(tokenize[1:], " "))
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
