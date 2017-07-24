package timenote

import (
	"fmt"
	"time"
)

type (
	// TimeEntry represents a simple note
	TimeEntry struct {
		// Tag is used for grouping
		Tag string
		// Some text attached to entry
		Note string
		// When has the author started working on the note
		Start time.Time
		// End time
		Stop *time.Time
	}
)

func (te *TimeEntry) ToString() string {
	return fmt.Sprintf("%v, %s,\n%s\n", te.Start, te.Tag, te.Note)
}
