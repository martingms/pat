package main

import (
	// TODO(mg): I don't like these kinds of imports, seems magical.
	"github.com/nsf/termbox-go"

	"runtime"
)

func drawLine(y, cl int) {
	termw, _ := termbox.Size()

	for i := 0; i < termw; i++ {
		termbox.SetCell(i, y, 'l', termbox.Attribute(cl), termbox.Attribute(cl))
	}
}

func main() {
	// We use all available cores.
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// TODO(mg): Find out which one we want to use.
	termbox.SetInputMode(termbox.InputEsc)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawLine(2, 3)
	termbox.Flush()

	// Main loop.
main_loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				break main_loop
			}
		}
	}
}
