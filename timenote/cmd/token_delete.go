package cmd

import (
	"log/slog"

	"github.com/zalando/go-keyring"
)

// tokenDelete deletes the token from the local keyring.
func tokenDelete() {
	if err := keyring.Delete("timenote", "token"); err != nil {
		slog.Error("error deleting token", "error", err)
		return
	}
	slog.Info("token successfully deleted")
}
