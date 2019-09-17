package persistence

import "livingit.de/code/timenote"

// Persistor is used to store the received data
type Persistor interface {
	// New starts a new timenote
	New() error
	// Append adds a line to the note
	Append(line string) error
	// Tag sets a tag for the note
	Tag(name string) error
	// Clients sets a client for the note
	Clients() ([]timenote.Client, error)
	// NewClient is used to create a new client
	NewClient(name string) error
	// Project adds time entry to a project
	Project(name string) error
	// Done writes the stop timestamp
	Done() error
	// Close the connection to the persistence backend
	Close() error
	// List of entries for the current day
	ListForDay() ([]timenote.TimeEntry, error)
	// Get currently running time entry
	Current() (*timenote.TimeEntry, error)
	// GetWebsite returns a URL to the time management system
	GetWebsite() (bool, string, error)
}
