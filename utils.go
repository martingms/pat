package main

import (
	"time"
	"net/mail"
)

// TODO(mg): This is copied to maildir/util. Find a better home to avoid repetition.
// TODO(mg): Is this really the best way to parse time?
// These formats seems to cover most mail.
// Others should be added as they are found in the wild.
var formats = []string{
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
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
}

func parseDate(h *mail.Header) (t time.Time, err error) {
	// We prefer Date header to Delivery-date.
	// TODO(mg): Should we?
	dateHeaders := []string{h.Get("Date"), h.Get("Delivery-date")}

header_loop:
	for _, val := range dateHeaders {
		for _, format := range formats {
			t, err = time.Parse(format, val)
			if err == nil {
				break header_loop
			}
		}
	}

	return t, err
}
