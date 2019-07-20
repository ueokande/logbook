package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type VerticalLine struct {
	view views.View
	views.WidgetWatchers
	char  rune
	style tcell.Style
}

func NewVerticalLine(char rune, style tcell.Style) *VerticalLine {
	return &VerticalLine{
		char:  char,
		style: style,
	}
}

func (w *VerticalLine) Draw() {
	if w.view == nil {
		return
	}
	w.view.Fill(w.char, w.style)
}

func (w *VerticalLine) Resize() {
	w.PostEventWidgetResize(w)
}

func (w *VerticalLine) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *VerticalLine) SetView(view views.View) {
	w.view = view
}

func (w *VerticalLine) Size() (int, int) {
	return 1, 1
}
