package cmd

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"go.livingit.de/reuse/flag"
)

var (
	workspace        *int
	outputFormat     *string
	separator        *string
	cacheMaxAge      *int
	cachePath        *string
	excludedProjects func() []string

	description       *string
	name              *string
	autoCreateProject *bool
	includeSeconds    *bool
	tokenFlag         *string
	sumOnly           *bool
	group             *bool
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fatal(err)
	}

	workspace = flag.Int("workspace", 0, "Set to work within this workspace, leave to zero to have it guessed (first workspace)")
	outputFormat = flag.String("output-format", "text", "text or json")
	separator = flag.String("separator", ";", "Separator for existing value and new value")
	cacheMaxAge = flag.Int("cache-max-age", 360, "Maximum age of cache in minutes")
	cachePath = flag.String("cache-path", path.Join(home, ".config/timenote"), "Where to store cache")
	excludedProjects = flag.StringSlice("excluded-projects", []string{}, "exclude projects from the list by name (comma separated)")

	description = flag.String("description", "", "Description for timestamp")
	name = flag.String("name", "", "Name value (tag, project or client)")
	autoCreateProject = flag.Bool("auto-create-project", false, "automatically create project if it does not exist")
	includeSeconds = flag.Bool("include-seconds", true, "Include seconds when writing out time entry")
	tokenFlag = flag.String("token", "", "toggl token to use")
	sumOnly = flag.Bool("sum-only", false, "Just print sum of timestamps")
	group = flag.Bool("group", false, "Print grouped by name")
}
