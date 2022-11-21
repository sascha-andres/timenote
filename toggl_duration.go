package timenote

import (
	"fmt"
	"time"
)
import "errors"

// TogglDuration is used to pretty print a toggl duration
type TogglDuration struct {
	duration int64

	calculated bool

	omitSeconds bool

	days    int64
	hours   int64
	minutes int64
	seconds int64
}

const (
	dayInSeconds    int64 = 24 * 60 * 60
	hourInSeconds   int64 = 60 * 60
	minuteInSeconds int64 = 60
)

// String returns a properly readable string for a duration in seconds
func (td *TogglDuration) String() string {
	if !td.omitSeconds {
		return td.withSeconds()
	}
	return td.withoutSeconds()
}

// OmitSeconds tells the converter to not print seconds
func (td *TogglDuration) OmitSeconds() {
	td.omitSeconds = true
}

// ShowSeconds tells the converter to print seconds
func (td *TogglDuration) ShowSeconds() {
	td.omitSeconds = false
}

func (td *TogglDuration) withSeconds() string {
	if td.days > 0 {
		return fmt.Sprintf("%dd %dh %02dm %02ds", td.days, td.hours, td.minutes, td.seconds)
	}
	if td.hours > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", td.hours, td.minutes, td.seconds)
	}
	if td.minutes > 0 {
		return fmt.Sprintf("%dm %02ds", td.minutes, td.seconds)
	}
	return fmt.Sprintf("%ds", td.seconds)
}

func (td *TogglDuration) withoutSeconds() string {
	if td.days > 0 {
		return fmt.Sprintf("%dd %dh %02dm", td.days, td.hours, td.minutes)
	}
	if td.hours > 0 {
		return fmt.Sprintf("%dh %02dm", td.hours, td.minutes)
	}
	if td.minutes > 0 {
		return fmt.Sprintf("%dm", td.minutes)
	}
	return "<1m"
}

// NewTogglDuration creates a new instance
func NewTogglDuration(duration int64) (*TogglDuration, error) {
	if duration < 0 {
		return nil, errors.New("negative values not allowed")
	}
	result := TogglDuration{
		duration: duration,
	}
	return &result, result.calculate()
}

func (td *TogglDuration) calculate() error {
	if td.duration < 0 {
		return errors.New("negative values not allowed")
	}
	td.calculateDone()
	return nil
}

func (td *TogglDuration) calculateDone() {
	localDuration := td.duration
	td.days = localDuration / dayInSeconds
	localDuration = localDuration - (td.days * dayInSeconds)
	td.hours = localDuration / hourInSeconds
	localDuration = localDuration - (td.hours * hourInSeconds)
	td.minutes = localDuration / minuteInSeconds
	localDuration = localDuration - (td.minutes * minuteInSeconds)
	td.seconds = localDuration
}

// TogglDurationFromTime formats a time like the duration
//
// Known caveat: this forgets about days
func TogglDurationFromTime(t time.Time) (*TogglDuration, error) {
	return NewTogglDuration(int64(t.Hour())*hourInSeconds + int64(t.Minute())*minuteInSeconds + int64(t.Second()))
}

func (td *TogglDuration) GetDuration() int64 {
	return td.duration
}
