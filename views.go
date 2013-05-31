package main

import (
	"./maildir"
	"fmt"
	"github.com/nsf/termbox-go"
)

var (
	// Predefined views.
	TitleView         = &titleView{}
	DirectoryListView = &directoryListView{0, 1, 0}

	// Shortcuts to those views.
	shortcuts = map[rune]view{
		't': TitleView,
		'd': DirectoryListView,
	}
)

type view interface {
	render()
	handleEvent(*termbox.Event)
}

///////////
// Views //
///////////

// TitleView
type titleView struct{} // TODO(mg): What if I don't need state?

func (v *titleView) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Start")

	drawPatLogo()

	width, height := termbox.Size()
	verStr := "v." + VERSION
	drawString(int((width-len(verStr))/2), int(height/4)+5, verStr)

	termbox.Flush()
}

func (v *titleView) handleEvent(ev *termbox.Event) {
	return
}

// DirectoryListView
type directoryListView struct {
	listPos int
	listMin int
	listMax int
}

func (v *directoryListView) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Directories")

	list := [][]string{}
	colWidths := []int{len(string(len(mdirs))) + 2, 0, 3}
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
	v.listMax = len(list)

	drawList(1, 3, colWidths, list, v.listPos)

	termbox.Flush()
}

func (v *directoryListView) handleEvent(ev *termbox.Event) {
	switch {
	// Common list operations.
	case ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown:
		if v.listPos < v.listMax {
			v.listPos += 1
		} else {
			v.listPos = v.listMin
		}
		v.render()
	case ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp:
		if v.listPos > v.listMin {
			v.listPos -= 1
		} else {
			v.listPos = v.listMax
		}
		v.render()
	case ev.Key == termbox.KeyEnter:
		return
	}
}

// DirectoryView
type directoryView struct {
	dir     *maildir.Maildir
	listPos int
	listMin int
	listMax int
}

func (v *directoryView) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine(v.dir.Name)

	termbox.Flush()
}

func (v *directoryView) handleEvent(ev *termbox.Event) {
	return
}
