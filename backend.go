package timenote

import (
	"errors"
	"time"
)

var (
	// ErrNoCurrentEntry is returned when no time entry is running
	ErrNoCurrentEntry = errors.New("no current time entry")
	// ErrProjectNotFound is returned when a project name or id does not exist
	ErrProjectNotFound = errors.New("project does not exist")
)

type (
	// Backend is a time tracking service (for example toggl). Caching is an
	// implementation detail of the backend.
	Backend interface {
		// New starts a new time entry with no description
		New() error
		// StartPrevious starts a new time entry based on the last finished time entry
		StartPrevious() error
		// Append adds line to the description of the running time entry, separated by separator
		Append(line, separator string) error
		// SetDescription replaces the description of the running time entry
		SetDescription(line string) error
		// Tag toggles the tag name on the running time entry
		Tag(name string) error
		// SetProjectForCurrentTimestamp sets the project of the running time entry.
		// An empty name removes the project.
		SetProjectForCurrentTimestamp(name string, autoCreateProject bool) error
		// Done stops the running time entry
		Done() error
		// Current returns the running time entry
		Current() (*TimeEntry, error)
		// ListForDay returns all time entries of today
		ListForDay() ([]TimeEntry, error)

		Projects() ([]Project, error)
		// CreateProject creates a project, if it does not exist
		CreateProject(name string) error
		// DeleteProject removes a project by name or id
		DeleteProject(nameOrID string) error

		Clients() ([]Client, error)
		// CreateClient creates a client, if it does not exist
		CreateClient(name string) error
	}

	// CachingBackend is a Backend that keeps a local cache.
	CachingBackend interface {
		Backend
		// UpdateCache reloads the cache from the service
		UpdateCache() error
		// CacheStatus returns the state of each part of the cache
		CacheStatus() ([]CacheStatus, error)
	}

	CacheStatus struct {
		Name       string    `json:"name"`
		Updated    time.Time `json:"updated"`
		NextUpdate time.Time `json:"next_update"`
	}
)
