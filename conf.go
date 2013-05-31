package main

// TODO(mg): This should be a parsed config file.
import (
	"github.com/nsf/termbox-go"
)

// Mail setup
const (
	MAILDIR_PATH = "/home/mg/dev/pat/maildir/test_mails/Mail"
)

// Misc
var (
	STARTUP_VIEW = TitleView // Which view the program starts in.
)

// Colorscheme
const (
	BAR_BG_COLOR   = termbox.ColorBlue
	BAR_TEXT_COLOR = termbox.ColorWhite

	FOCUS_BAR_BG_COLOR = termbox.ColorCyan
	FOCUS_BAR_TEXT_COLOR = termbox.ColorBlack
)
