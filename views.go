package main

import (
	"github.com/nsf/termbox-go"
	"./maildir"
	"fmt"
)

var (
	// Predefined views.
	TitleView = &view{renderFunc: titleView, keyHandlerFunc: titleKeyHandler}
	DirectoryListView = &view{renderFunc: directoryListView, keyHandlerFunc: directoryListKeyHandler}

	// Shortcuts to those views.
	shortcuts = map[rune]*view{
		't' : TitleView,
		'd' : DirectoryListView,
	}
)

// TODO(mg): Keeping some state here for the time being.
// Find a better way as this will fuck shit up.
var (
	listPos = 0
	listMin = 1
	listMax = 0
)

type view struct {
	renderFunc     func(interface{})
	keyHandlerFunc func(interface{}, *termbox.Event)
	state          interface{}
}

func (v *view) render() {
	v.renderFunc(v.state)
}

func (v *view) keyHandler(ev *termbox.Event) {
	v.keyHandlerFunc(v.state, ev)
}

// Views
func titleView(state interface{}) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Start")

	drawPatLogo()

	width, height := termbox.Size()
	verStr := "v." + VERSION
	drawString(int((width-len(verStr))/2), int(height/4)+5, verStr)

	termbox.Flush()
}

func titleKeyHandler(state interface{}, ev *termbox.Event) {
	return
}

func directoryListView(state interface{}) {
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

func directoryListKeyHandler(state interface{}, ev *termbox.Event) {
	switch {
	case ev.Key == termbox.KeyEnter:
		// TODO(mg): Get what place in mdirs we are.
	}
}

func directoryViewBuilder(dir *maildir.Maildir) *view {
	return &view{
		renderFunc: directoryView,
		keyHandlerFunc: directoryKeyHandler,
		state: dir,
	}
}

func directoryView(state interface{}) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// TODO(mg): Get what place in mdirs we are.
	drawTopLine(".personlig")

	termbox.Flush()
}

func directoryKeyHandler(state interface{}, ev *termbox.Event) {

}
