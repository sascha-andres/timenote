package timenote

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ErrNoCurrentTimeEntry should be returned in case no running timeentry is found
var ErrNoCurrentTimeEntry = errors.New("timenote: no current timeentry")

type (
	// TimeEntry represents a simple note
	TimeEntry struct {
		// ID is a system id which may be set from a persistor
		ID int
		// Tag is used for grouping
		Tag string
		// Some text attached to entry
		Note string
		// When has the author started working on the note
		Start time.Time
		// End time
		Stop *time.Time
		// time entry duration in seconds. If the time entry is currently running, the duration attribute contains a negative value, denoting the start of the time entry in seconds since epoch (Jan 1 1970). The correct duration can be calculated as current_time + duration, where current_time is the current time in seconds since epoch.
		Duration int64
	}
)

func (te *TimeEntry) String() string {
	if "[]" == te.Tag {
		if "" == te.Note {
			return fmt.Sprintf("duration: %s\n", te.Start.Format("15:04:05"))
		}
		return fmt.Sprintf("duration: %s\nnote: %s\n", te.Start.Format("15:04:05"), te.Note)
	}
	if "" == te.Note {
		return fmt.Sprintf("duration: %s - tags:%s\n", te.Start.Format("15:04:05"), te.Tag)
	}
	return fmt.Sprintf("duration: %s - tags:%s\n%s\n", te.Start.Format("15:04:05"), te.Tag, te.Note)
}
