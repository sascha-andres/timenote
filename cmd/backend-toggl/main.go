// Command backend-toggl exposes the toggl timenote.Backend as a REST API.
package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"

	"go.livingit.de/reuse/flag"

	"go.livingit.de/timenote/internal/persistence"
)

func main() {
	// env fallback: BACKEND_TOGGL_<FLAG>, except TOGGL_TOKEN for -toggl-token
	flag.SetEnvPrefix("BACKEND_TOGGL")
	flag.SetEnvPrefixForFlag("toggl-token", "")
	listen := flag.String("listen", "127.0.0.1:8141", "address to listen on")
	token := flag.String("token", "", "token clients must send as \"Authorization: Bearer <token>\" (required)")
	togglToken := flag.String("toggl-token", "", "toggl api token (default: the timenote keyring entry)")
	workspace := flag.Int("workspace", 0, "toggl workspace, 0 uses the first workspace")
	cacheDir := flag.String("cache-path", "", "cache directory (default $HOME/.config/timenote/cache)")
	cacheMaxAge := flag.Int("cache-max-age", 360, "maximum age of cache in minutes")
	readCacheTTL := flag.Duration("read-cache-ttl", 30*time.Second, "how long GET results are kept, 0 disables")
	flag.Parse()

	if err := run(*listen, *token, *togglToken, *workspace, *cacheDir, *cacheMaxAge, *readCacheTTL); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(listen, token, togglToken string, workspace int, cacheDir string, cacheMaxAge int, readCacheTTL time.Duration) error {
	if token == "" {
		return errors.New("-token or $BACKEND_TOGGL_TOKEN is required")
	}
	if togglToken == "" {
		togglToken, _ = keyring.Get("timenote", "token")
	}
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		cacheDir = filepath.Join(home, ".config", "timenote", "cache")
	}
	backend, err := persistence.NewToggl(togglToken, workspace, cacheDir, cacheMaxAge)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              listen,
		Handler:           newHandler(backend, token, readCacheTTL),
		ReadHeaderTimeout: 10 * time.Second,
	}
	slog.Info("listening", "address", listen)
	return srv.ListenAndServe()
}
