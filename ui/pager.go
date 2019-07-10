package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Pager struct {
	view     views.View
	viewport views.ViewPort
	text     views.Text

	views.WidgetWatchers
}

func NewPager() *Pager {
	w := &Pager{}
	w.text.SetView(&w.viewport)
	w.text.SetStyle(tcell.StyleDefault)
	return w
}

func (w *Pager) AppendLine(line string) {
	text := w.text.Text()
	if len(text) > 0 {
		text += "\n"
	}
	text += line
	w.text.SetText(text)

	width, height := w.text.Size()
	w.viewport.SetContentSize(width, height, true)
	w.viewport.ValidateView()
}

func (w *Pager) ScrollDown() {
	w.viewport.ScrollDown(1)
}

func (w *Pager) ScrollUp() {
	w.viewport.ScrollUp(1)
}

func (w *Pager) ScrollHalfPageDown() {
	_, vh := w.view.Size()
	w.viewport.ScrollDown(vh / 2)
}

func (w *Pager) ScrollHalfPageUp() {
	_, vh := w.view.Size()
	w.viewport.ScrollUp(vh / -2)
}

func (w *Pager) ScrollPageDown() {
	_, vh := w.view.Size()
	w.viewport.ScrollDown(vh)
}

func (w *Pager) ScrollPageUp() {
	_, vh := w.view.Size()
	w.viewport.ScrollUp(vh)
}

func (w *Pager) ClearText() {
	w.text.SetText("")

	width, height := w.text.Size()
	w.viewport.SetContentSize(width, height, true)
	w.viewport.ValidateView()
}

func (w *Pager) Draw() {
	if w.view == nil {
		return
	}
	w.text.Draw()
}

func (w *Pager) Resize() {
	width, height := w.view.Size()
	w.viewport.Resize(0, 0, width, height)
	w.viewport.ValidateView()
}

func (w *Pager) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *Pager) SetView(view views.View) {
	w.view = view
	w.viewport.SetView(view)
	if view == nil {
		return
	}
	w.Resize()
}

func (w *Pager) Size() (int, int) {
	width, height := w.view.Size()
	if width > 2 {
		width = 2
	}
	if height > 2 {
		height = 2
	}
	return width, height
}
