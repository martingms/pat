package main

import (
	// TODO(mg): I don't like these kinds of imports, seems magical and prone
	// to breakage. Host my own fork or something, preferably no remote imports
	// at all.
	"github.com/nsf/termbox-go"

	"./maildir"
	"runtime"
)

const (
	VERSION = "0.1.12"
)

var (
	cv view // The current view
	// TODO(mg): Abstract this.
	mdirs []*maildir.Maildir
)

func main() {
	// We use all available cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := termbox.Init()
	if err != nil {
		// TODO(mg): Gracefully quit, don't panic.
		panic(err)
	}
	defer termbox.Close()

	// TODO(mg): Find out which one we want to use.
	termbox.SetInputMode(termbox.InputEsc)

	cv = STARTUP_VIEW
	cv.render()

	// Initialize maildirs.
	// TODO(mg): Abstract this to allow imap, other mailbox-specs etc.
	mdir, err := maildir.NewMaildir(MAILDIR_PATH)
	if err != nil {
		// TODO(mg): Gracefully quit, don't panic.
		panic(err)
	}
	// We want relative names.
	mdir.Name = "."

	mdirs, err = mdir.ListMaildirs()
	if err != nil {
		panic(err)
	}
	// There might be mail in the root directory as well.
	// TODO(mg): Sort this list.
	mdirs = append(mdirs, mdir)

	// Main loop.
main_loop:
	for {
	event_switch:
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			// TODO(mg): Can we do something to break up this if/else chain nicely?
			// Global keys.
			if ev.Key == termbox.KeyCtrlC {
				break main_loop
			}

			// Common list operations.
			switch {
			case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
				if listPos < listMax {
					listPos += 1
				} else {
					listPos = listMin
				}
				cv.render()
				break event_switch

			case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
				if listPos > listMin {
					listPos -= 1
				} else {
					listPos = listMax
				}
				cv.render()
				break event_switch
			}

			// Shortcuts
			if v, ok := shortcuts[ev.Ch]; ok {
				cv = v
				cv.render()
			}

			// All other keys should be handled by the current view.
			cv.handleEvent(&ev)
		case termbox.EventResize:
			cv.render()
			break event_switch
		case termbox.EventError:
			// TODO(mg): Probably shouldn't panic.
			panic(ev.Err)
		}
	}
}
