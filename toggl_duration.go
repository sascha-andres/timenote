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

func (td *TogglDuration) String() string {
	if td.days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", td.days, td.hours, td.minutes, td.seconds)
	}
	if td.hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", td.hours, td.minutes, td.seconds)
	}
	if td.minutes > 0 {
		return fmt.Sprintf("%dm %ds", td.minutes, td.seconds)
	}
	return fmt.Sprintf("%ds", td.seconds)
}

// NewTogglDuration creates a new instance
func NewTogglDuration(duration int64) (*TogglDuration, error) {
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

func (td *TogglDuration) FromTime(t time.Time) string {
	if t.Hour() > 0 {
		return fmt.Sprintf("%dh %dm %ds", t.Hour(), t.Minute(), t.Second())
	}
	if t.Minute() > 0 {
		return fmt.Sprintf("%dm %ds", t.Minute(), t.Second())
	}

	return fmt.Sprintf("%ds", t.Second())
}
