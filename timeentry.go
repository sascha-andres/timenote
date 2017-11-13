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
		// Id is a systerm id which may be set from a persistor
		Id int
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
