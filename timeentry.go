package timenote

import (
	"fmt"
	"time"
)

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
		// TimeEntry belongs to project
		Project string
		// // TimeEntry belongs to client
		Client string
	}
)

func (te *TimeEntry) String() string {
	humanTime := te.getHumanTime()
	result := fmt.Sprintf("client, project: %s, %s\n", te.Client, te.Project)
	result = result + fmt.Sprintf("duration: %s\n", humanTime)
	if "[]" != te.Tag {
		result = result + fmt.Sprintf("tags: %s\n", te.Tag)
	}
	result = result + fmt.Sprintf("note: %s\n", te.Note)
	return result
}

func (te *TimeEntry) getHumanTime() string {
	humanTime := ""
	if te.Duration < 0 {
		t := time.Now().UTC().Add(time.Duration(te.Duration) * time.Second)
		humanTime = t.Format("15:04:05")
	} else {
		td, err := NewTogglDuration(te.Duration)
		if err != nil {
			panic(err)
		}
		humanTime = td.String()
	}
	return humanTime
}
