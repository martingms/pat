package main

import (
  "github.com/nsf/termbox-go"
	"errors"
)

var cv string // The current view.
var views = map[string]*view{
	"titleView": &view{shortcut: 'd', renderFunc: titleView, keyHandlerFunc: titleKeyHandler},
}

type view struct {
	shortcut rune
	renderFunc func()
	keyHandlerFunc func(*termbox.Event)
	// TODO(mg): May need some state?
}

func titleView() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawTopLine("Start")

	drawPatLogo()
	
	width, height := termbox.Size()
	verStr := "v."+VERSION
	drawString(int((width-len(verStr))/2), int(height/4)+5, verStr)

	termbox.Flush()
}

func titleKeyHandler(ev *termbox.Event) {
	return
}

func mainView() error {
	return errors.New("mainView not yet implemented!")	
}

// Helper functions

func cvRender() {
	if view, ok := views[cv]; ok {
		view.renderFunc()
	}
	//TODO(mg): Show some error message perhaps?
}

func cvKeyHandler(ev *termbox.Event) {
	if view, ok := views[cv]; ok {
		view.keyHandlerFunc(ev)
	}
	//TODO(mg): Show some error message perhaps?
}
