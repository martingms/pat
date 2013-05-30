package main

// TODO(mg): This should be a parsed config file.
import (
	"github.com/nsf/termbox-go"
)

const (
	STARTUP_VIEW = "titleView" // Which view the program starts in.
)

// Colorscheme
// Should be pretty self explanatory.
const (
	BAR_BG_COLOR = termbox.ColorBlue
	BAR_TEXT_COLOR = termbox.ColorWhite
)
