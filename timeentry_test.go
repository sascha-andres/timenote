package timenote_test

import (
	"testing"
	"time"

	"go.livingit.de/timenote"
)

func TestTimeEntrySeconds(t *testing.T) {
	running := timenote.TimeEntry{Start: time.Now().Add(-90 * time.Minute), Duration: -1}
	if got := running.Seconds(); got < 5400 || got > 5401 {
		t.Errorf("running Seconds() = %d, want 5400", got)
	}
	done := timenote.TimeEntry{Start: time.Now().Add(-5 * time.Hour), Duration: 60}
	if got := done.Seconds(); got != 60 {
		t.Errorf("stopped Seconds() = %d, want 60", got)
	}
}
