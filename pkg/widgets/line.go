package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// VerticalLine is a Widget with containing a vertical line.
type VerticalLine struct {
	view views.View
	views.WidgetWatchers
	char  rune
	style tcell.Style
}

// NewVerticalLine creates a new VerticalLine
func NewVerticalLine(char rune, style tcell.Style) *VerticalLine {
	return &VerticalLine{
		char:  char,
		style: style,
	}
}

// Draw draws the VerticalLine
func (w *VerticalLine) Draw() {
	if w.view == nil {
		return
	}
	w.view.Fill(w.char, w.style)
}

// Resize is called when our View changes sizes.
func (w *VerticalLine) Resize() {
	w.PostEventWidgetResize(w)
}

// HandleEvent handles events on tcell
func (w *VerticalLine) HandleEvent(ev tcell.Event) bool {
	return false
}

// SetView sets the View object used for the vertical line
func (w *VerticalLine) SetView(view views.View) {
	w.view = view
}

// Size returns the width and height in vertical line, which always (1, 1).
func (w *VerticalLine) Size() (int, int) {
	return 1, 1
}
