package timenote_test

import (
	"testing"

	"go.livingit.de/timenote"
)

func TestTogglDurationString(t *testing.T) {
	cases := []struct {
		seconds int64
		want    string
	}{
		{10, "10s"},
		{60, "1m 00s"},
		{3600, "1h 00m 00s"},
		{3666, "1h 01m 06s"},
		{3600 * 24, "1d 0h 00m 00s"},
	}
	for _, c := range cases {
		td, err := timenote.NewTogglDuration(c.seconds)
		if err != nil {
			t.Fatalf("NewTogglDuration(%d): %v", c.seconds, err)
		}
		if got := td.String(); got != c.want {
			t.Errorf("NewTogglDuration(%d).String() = %q, want %q", c.seconds, got, c.want)
		}
	}
}

func TestTogglDurationStringOmitSeconds(t *testing.T) {
	cases := []struct {
		seconds int64
		want    string
	}{
		{10, "<1m"},
		{60, "1m"},
		{3600, "1h 00m"},
		{3666, "1h 01m"},
		{3600 * 24, "1d 0h 00m"},
	}
	for _, c := range cases {
		td, err := timenote.NewTogglDuration(c.seconds)
		if err != nil {
			t.Fatalf("NewTogglDuration(%d): %v", c.seconds, err)
		}
		td.OmitSeconds()
		if got := td.String(); got != c.want {
			t.Errorf("NewTogglDuration(%d) omit-seconds String() = %q, want %q", c.seconds, got, c.want)
		}
	}
}

func TestNewTogglDurationNegative(t *testing.T) {
	for _, seconds := range []int64{-10, -60, -3600, -3666, -3600 * 24} {
		_, err := timenote.NewTogglDuration(seconds)
		if err == nil {
			t.Fatalf("NewTogglDuration(%d): expected error, got nil", seconds)
		}
		if got, want := err.Error(), "negative values not allowed"; got != want {
			t.Errorf("NewTogglDuration(%d) error = %q, want %q", seconds, got, want)
		}
	}
}
