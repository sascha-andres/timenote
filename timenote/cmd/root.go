// Copyright © 2021 Sascha Andres <sascha.andres@outlook.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gops/agent"
	"github.com/zalando/go-keyring"
	"go.livingit.de/reuse/flag"
	"go.livingit.de/timenote/internal/cache"
	"go.livingit.de/timenote/internal/persistence"
)

var (
	token   string
	caching *cache.Cache
)

func init() {
	token, _ = keyring.Get("timenote", "token")
}

// fatal logs err and terminates the process.
func fatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

// fatalf logs msg with err and terminates the process.
func fatalf(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

// requireString exits the process if value is empty, naming flagName in the error.
func requireString(flagName, value string) string {
	if value == "" {
		fatal(fmt.Errorf("--%s is required", flagName))
	}
	return value
}

// newPersistor builds the persistence layer for the configured workspace, or
// terminates the process if that fails.
func newPersistor() *persistence.TogglPersistor {
	p, err := persistence.NewToggl(token, *workspace, caching)
	if err != nil {
		fatal(err)
	}
	return p
}

// Execute parses flags and dispatches to the requested command.
func Execute() {
	flag.Parse()
	configureLogging()

	if err := agent.Listen(agent.Options{}); err != nil {
		fatal(err)
	}

	c, err := cache.NewCache(*cacheMaxAge, resolveCacheDir())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	caching = c

	dispatch(flag.GetVerbs())
}

// resolveCacheDir returns the directory cache files are stored in. When
// --cache-path is left at its default ($HOME/.config/timenote), a "cache"
// subdirectory is used so cache files don't clutter the config directory; an
// explicitly provided --cache-path is used as-is.
func resolveCacheDir() string {
	if *cachePath == defaultCachePath {
		return filepath.Join(*cachePath, "cache")
	}
	return *cachePath
}

// dispatch routes the leading verb to its command group, falling back to
// treating the whole verb list as a new timestamp's description.
func dispatch(verbs []string) {
	if len(verbs) == 0 {
		return
	}
	switch verbs[0] {
	case "timestamp":
		dispatchTimestamp(verbs[1:])
	case "projects":
		dispatchProjects(verbs[1:])
	case "clients":
		dispatchClients(verbs[1:])
	case "cache":
		dispatchCache(verbs[1:])
	case "token":
		dispatchToken(verbs[1:])
	case "today":
		timestampToday()
	case "browser":
		openBrowser()
	case "i":
		if err := run(); err != nil {
			fatal(err)
		}
	default:
		newFromDescription(verbs)
	}
}

// newFromDescription starts a new timestamp using the given words joined as
// its description; this is the default action for the root command.
func newFromDescription(words []string) {
	p := newPersistor()
	if err := p.New(); err != nil {
		fatal(err)
	}
	_ = p.Append(strings.Join(words, " "), *separator)
}
