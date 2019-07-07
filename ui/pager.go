package ui

import (
	"strings"

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
	offset   int
	lines    int

	views.WidgetWatchers
}

func NewPager() *Pager {
	w := &Pager{
		offset: 3,
	}
	w.text.SetView(&w.viewport)
	w.text.SetStyle(tcell.StyleDefault)
	return w
}

func (w *Pager) WriteText(content string) {
	prevLines := w.lines

	w.lines += strings.Count(content, "\n")
	w.content += content
	_, vh := w.view.Size()
	if prevLines-w.offset > vh {
		// Out of viewport
		return
	}
	w.text.SetText(tail(w.content, w.offset))
	w.text.Watch(w)

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *Pager) ScrollDown() {
	w.scrollBy(1)
}

func (w *Pager) ScrollUp() {
	w.scrollBy(-1)
}

func (w *Pager) ScrollHalfPageDown() {
	_, vh := w.view.Size()
	w.scrollBy(vh / 2)
}

func (w *Pager) ScrollHalfPageUp() {
	_, vh := w.view.Size()
	w.scrollBy(vh / -2)
}

func (w *Pager) ScrollPageDown() {
	_, vh := w.view.Size()
	w.scrollBy(vh)
}

func (w *Pager) ScrollPageUp() {
	_, vh := w.view.Size()
	w.scrollBy(-vh)
}

func (w *Pager) scrollBy(count int) {
	offset := w.offset + count
	_, vh := w.view.Size()
	if offset >= w.lines-vh {
		offset = w.lines - vh
	}
	if offset < 0 {
		offset = 0
	}
	if offset == w.offset {
		return
	}
	w.offset = offset
	w.text.SetText(tail(w.content, w.offset))

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *Pager) ClearText() {
	w.content = ""
	w.offset = 0
	w.lines = 0
	w.text.SetText("")

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

func tail(str string, line int) string {
	if line <= 0 {
		return str
	}
	idx := strings.Index(str, "\n")
	if idx < 0 {
		return ""
	}
	return tail(str[idx+1:], line-1)
}
