package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type ListItem struct {
	Text  string
	Style tcell.Style
}

type ListView struct {
	view  views.View
	texts []*views.Text
	views []*views.ViewPort

	views.WidgetWatchers
}

func NewListView() *ListView {
	return &ListView{}
}

func (w *ListView) AddItem(item ListItem) {
	v := &views.ViewPort{}

	t := &views.Text{}
	t.SetText(item.Text)
	t.SetStyle(item.Style)
	t.SetView(v)
	t.Watch(w)

	w.views = append(w.views, v)
	w.texts = append(w.texts, t)

	w.PostEventWidgetContent(w)
}

func (w *ListView) Draw() {
	for i, t := range w.texts {
		textw, texth := t.Size()

		v := w.views[i]
		v.Resize(0, i, textw, texth)

		t.Draw()
	}
}

func (w *ListView) Resize() {
	for _, t := range w.texts {
		t.Resize()
	}
	w.PostEventWidgetResize(w)
}

func (w *ListView) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *ListView) SetView(view views.View) {
	w.view = view
	for _, v := range w.views {
		v.SetView(view)
	}
}

func (w *ListView) Size() (int, int) {
	var width int
	for _, t := range w.texts {
		tw, _ := t.Size()
		if tw > width {
			width = tw
		}
	}
	return width, len(w.texts)
}
