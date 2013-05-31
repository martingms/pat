package main

import (
	"time"
)

// TODO(mg): Is this really the best way to parse time?
var formats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano, // lol
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
}

func parseDate(val string) (t time.Time, err error) {
	for _, format := range formats {
		t, err = time.Parse(format, val)
		if err == nil {
			break
		}
	}

	return t, err
}
