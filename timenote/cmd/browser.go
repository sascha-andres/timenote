package cmd

import (
	"log/slog"

	"github.com/pkg/browser"
)

// openBrowser opens the toggl dashboard in the default browser.
func openBrowser() {
	if err := browser.OpenURL("https://toggl.com/app/timer"); err != nil {
		slog.Error("error executing browser", "error", err)
	}
}
