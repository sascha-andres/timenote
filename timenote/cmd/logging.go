package cmd

import (
	"log/slog"
	"os"
)

// configureLogging sets up the default slog logger from the -log-level and
// -log-json flags. Must be called after flag.Parse().
func configureLogging() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		fatal(err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	}

	var handler slog.Handler
	if *logJSON {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}
