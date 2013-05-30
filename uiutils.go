package main

import (
	"github.com/nsf/termbox-go"
)

// Top line should be consistent throughout.
func drawTopLine(view string) {
	width, _ := termbox.Size()
	ver := "pat v." + VERSION

	drawBlankLine(0, BAR_BG_COLOR)

	for i := 0; i < len(ver); i++ {
		termbox.SetCell(i+1, 0, rune(ver[i]), BAR_TEXT_COLOR, BAR_BG_COLOR)
	}

	for i := 0; i < len(view); i++ {
		termbox.SetCell(width-i-2, 0, rune(view[len(view)-i-1]),
			BAR_TEXT_COLOR, BAR_BG_COLOR)
	}
}

func drawBlankLine(y int, color termbox.Attribute) {
	width, _ := termbox.Size()
	for i := 0; i < width; i++ {
		termbox.SetCell(i, y, ' ', color, color)
	}
}

func drawString(x, y int, str string) {
	for i, char := range str {
		termbox.SetCell(x+i, y, char, termbox.ColorDefault, termbox.ColorDefault)	
	}
}

// Draws pat logo as close to center as it can.
// TODO(mg): Make a cooler logo.
func drawPatLogo() {
	width, height := termbox.Size()

	startx := int((width-8)/2)
	starty := int(height/4)

	// Row 1
	for i := startx; i < startx+9; i++ {
		termbox.SetCell(i, starty, ' ', termbox.ColorDefault, termbox.ColorRed)
	}

	// Row 2
	termbox.SetCell(startx, starty+1, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+2, starty+1, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+3, starty+1, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+5, starty+1, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+7, starty+1, ' ', termbox.ColorDefault, termbox.ColorRed)

	// Row 3
	termbox.SetCell(startx, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+1, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+2, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+3, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+4, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+5, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+7, starty+2, ' ', termbox.ColorDefault, termbox.ColorRed)

	// Row 4
	termbox.SetCell(startx, starty+3, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+3, starty+3, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+5, starty+3, ' ', termbox.ColorDefault, termbox.ColorRed)
	termbox.SetCell(startx+7, starty+3, ' ', termbox.ColorDefault, termbox.ColorRed)
}
