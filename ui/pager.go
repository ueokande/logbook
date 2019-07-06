package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Pager struct {
	view     views.View
	viewport views.ViewPort
	content  string
	changed  bool
	width    int
	height   int
	text     views.Text

	views.WidgetWatchers
}

func NewPager() *Pager {
	w := &Pager{}
	w.text.SetView(&w.viewport)
	w.text.SetStyle(tcell.StyleDefault)
	return w
}

func (w *Pager) WriteText(content string) {
	w.content += content
	w.text.SetText(w.content)

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *Pager) ClearText() {
	w.content = ""
	w.text.SetText(w.content)

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *Pager) Draw() {
	if w.view == nil {
		return
	}
	if w.changed {
		w.layout()
	}
	w.text.Draw()
}

func (w *Pager) Resize() {
	w.layout()
	w.PostEventWidgetResize(w)
}

func (w *Pager) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *Pager) SetView(view views.View) {
	w.view = view
	w.viewport.SetView(view)
	w.changed = true
}

func (w *Pager) Size() (int, int) {
	return w.width, w.height
}

func (w *Pager) layout() {
	w.width, w.height = w.text.Size()
	w.viewport.Resize(0, 0, w.width, w.height)
	w.text.Resize()

	w.changed = false
}
