package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type item struct {
	text *views.Text
	view *views.ViewPort
}

type ListView struct {
	view     views.View
	items    []item
	selected int
	changed  bool
	width    int
	height   int

	views.WidgetWatchers
}

func NewListView() *ListView {
	return &ListView{}
}

func (w *ListView) AddItem(text string, style tcell.Style) {
	v := &views.ViewPort{}
	v.SetView(w.view)

	t := &views.Text{}
	t.SetText(text)
	t.SetStyle(style)
	t.SetView(v)
	t.Watch(w)

	i := item{text: t, view: v}
	w.items = append(w.items, i)

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *ListView) SelectAt(index int) {
	if index == w.selected {
		return
	}
	if w.selected >= 0 {
		i := w.items[w.selected]
		i.text.SetStyle(i.text.Style().Reverse(false))
	}
	if index < 0 || index >= len(w.items) {
		return
	}

	w.selected = index
	i := w.items[index]
	i.text.SetStyle(i.text.Style().Reverse(true))

	w.PostEventWidgetContent(w)
}

func (w *ListView) Draw() {
	if w.view == nil {
		return
	}
	if w.changed {
		w.layout()
	}
	for _, i := range w.items {
		i.text.Draw()
	}
}

func (w *ListView) Resize() {
	w.layout()
	w.PostEventWidgetResize(w)
}

func (w *ListView) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.changed = true
		w.PostEventWidgetContent(w)
		return true
	}
	for _, item := range w.items {
		if item.text.HandleEvent(ev) {
			return true
		}
	}
	return false

}

func (w *ListView) SetView(view views.View) {
	w.view = view
	for _, item := range w.items {
		item.view.SetView(view)
	}
	w.changed = true
}

func (w *ListView) Size() (int, int) {
	return w.width, w.height
}

func (w *ListView) layout() {
	w.width, w.height = 0, 0
	for y, item := range w.items {
		textw, texth := item.text.Size()
		if textw > w.width {
			w.width = textw
		}
		item.view.Resize(0, y, textw, texth)
		item.text.Resize()
	}
	w.height = len(w.items)
	w.changed = false
}
