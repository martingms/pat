package main

import (
	"github.com/nsf/termbox-go"
	//"./maildir"
	"fmt"
)

// TODO(mg): This way of doing this is cute but prone to fail. Find a new strategy.
// It's getting hard to avoid circular definitions.
//
// The current view
var cv string
var views = map[string]*view{
	// We don't really need a shortcut here.
	"titleView": &view{shortcut: 't', renderFunc: titleView, keyHandlerFunc: titleKeyHandler},
	"directoryListView": &view{shortcut: 'd', renderFunc: directoryListView, keyHandlerFunc: directoryListKeyHandler},
}

type view struct {
	shortcut       rune
	renderFunc     func()
	keyHandlerFunc func(*termbox.Event)
	// TODO(mg): May need some state?
}

// TODO(mg): Keeping some state here for the time being.
// Find a better way as this will fuck shit up.
var (
	listPos = 0
	listMin = 1
	listMax = 0
)

func resetList() {
	listPos = 0
	listMin = 1
	listMax = 0
}

// Views
func titleView() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Start")

	drawPatLogo()

	width, height := termbox.Size()
	verStr := "v." + VERSION
	drawString(int((width-len(verStr))/2), int(height/4)+5, verStr)

	termbox.Flush()
}

func titleKeyHandler(ev *termbox.Event) {
	return
}

func directoryListView() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Directories")

	list := [][]string{}
	colWidths := []int{len(string(len(mdirs)))+2, 0, 3}
	for i, dir := range mdirs {
		hasnew := ""
		if dir.HasNewMail() {
			hasnew = "N"
		}
		list = append(list, []string{fmt.Sprint(i), dir.Name, hasnew})
		if len(dir.Name) > colWidths[1] {
			colWidths[1] = len(dir.Name)
		}
	}
	colWidths[1] += 2
	listMax = len(list)

	drawList(1, 3, colWidths, list, listPos)

	termbox.Flush()
}

func directoryListKeyHandler(ev *termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEnter:
		// TODO(mg): Get what place in mdirs we are.
		cv = "directoryView"
		directoryView()
	}
}

func directoryView() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// TODO(mg): Get what place in mdirs we are.
	drawTopLine(".personlig")

	termbox.Flush()
}

func directoryKeyHandler(ev *termbox.Event) {

}

/////////////////////////////////////////////////
// Helper functions
func cvRender() {
	if view, ok := views[cv]; ok {
		view.renderFunc()
	}
}

func cvKeyHandler(ev *termbox.Event) {
	if view, ok := views[cv]; ok {
		view.keyHandlerFunc(ev)
	}
}
