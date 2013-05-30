package main

import (
	// TODO(mg): I don't like these kinds of imports, seems magical and prone
  // to breakage. Host my own fork or something, preferably no remote imports
  // at all.
	"github.com/nsf/termbox-go"

	"runtime"
)

const(
  VERSION = "0.1.12"
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
  views[cv].renderFunc()

  // TODO(mg): Do loop in each view instead of single mainloop?
  // Specific stuff, then default: Handle menus etc.
	// Main loop.
main_loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
      // TODO(mg): Can we do something to break up this if/else chain nicely?
      // Global keys like changing view
			if ev.Key == termbox.KeyCtrlC {
				break main_loop
			} else if ev.Ch == 'd' {
        cv = "directoryView"
        cvRender()

      // All other keys should be handled by the current view.
      } else {
        cvKeyHandler(&ev)
      }
    case termbox.EventResize:
      cvRender()
    case termbox.EventError:
      panic(ev.Err)
		}
	}
}
