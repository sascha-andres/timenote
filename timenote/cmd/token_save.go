package cmd

import (
	"log/slog"

	"github.com/zalando/go-keyring"
)

// tokenSave saves the toggl token to the local keyring.
func tokenSave() {
	t := requireString("token", *tokenFlag)
	if err := keyring.Set("timenote", "token", t); err != nil {
		slog.Error("could not set token", "error", err)
		return
	}
	slog.Info("token set")
}
