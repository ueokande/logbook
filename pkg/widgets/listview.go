package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type item struct {
	name string
	text *views.Text
	view *views.ViewPort
}

// ListView is a Widget with containing multiple items as a list
type ListView struct {
	view     views.View
	items    []item
	selected int
	changed  bool
	width    int
	height   int

	views.WidgetWatchers
}

// NewListView returns a new ListView
func NewListView() *ListView {
	return &ListView{
		selected: -1,
	}
}

// AddItem adds a new item with the text and its style.  The text must be
// unique in the list view.  It panics when the text is already exists in the
// list
func (w *ListView) AddItem(text string, style tcell.Style) {
	if w.getItemIndex(text) != -1 {
		panic("item " + text + " already exists")
	}

	v := &views.ViewPort{}
	v.SetView(w.view)

	t := &views.Text{}
	t.SetText(text)
	t.SetStyle(style)
	t.SetView(v)
	t.Watch(w)

	item := item{name: text, text: t, view: v}
	w.items = append(w.items, item)

	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

// SetStyle updates the style of the text.  It panics when the text does not
// exist in the list.
func (w *ListView) SetStyle(text string, style tcell.Style) {
	idx := w.getItemIndex(text)
	if idx == -1 {
		panic("item " + text + " not fount")
	}

	w.items[idx].text.SetStyle(style)
	w.changed = true
	w.layout()
	w.PostEventWidgetContent(w)
}

// DeleteItem deletes a item with the text.  It panics when the text does not
// exist in the list
func (w *ListView) DeleteItem(text string) {
	idx := w.getItemIndex(text)
	if idx == -1 {
		panic("item " + text + " not fount")
	}
	item := w.items[idx]
	item.text.Unwatch(w)
	w.items = append(w.items[:idx], w.items[idx+1:]...)
}

// ItemCount returns the count of the items.
func (w *ListView) ItemCount() int {
	return len(w.items)
}

func (w *ListView) getItemIndex(name string) int {
	for i, item := range w.items {
		if item.name == name {
			return i
		}
	}
	return -1
}

// SelectAt selects nth items by the index.
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

	ev := &EventItemSelected{
		Name:   i.name,
		Index:  index,
		widget: w,
	}
	ev.SetEventNow()

	w.PostEvent(ev)
}

// Draw draws the ListView
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

// Resize is called when our View changes sizes.
func (w *ListView) Resize() {
	w.layout()
	w.PostEventWidgetResize(w)
}

// HandleEvent handles events on tcell
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

// SetView sets the View object used for the list view
func (w *ListView) SetView(view views.View) {
	w.view = view
	for _, item := range w.items {
		item.view.SetView(view)
	}
	w.changed = true
}

// Size returns the width and height in vertical line.
func (w *ListView) Size() (int, int) {
	return w.width, w.height
}

func (w *ListView) layout() {
	vieww, _ := w.view.Size()
	w.width, w.height = 0, 0
	for y, item := range w.items {
		textw, texth := item.text.Size()
		if textw > w.width {
			w.width = textw
		}
		if textw < vieww {
			textw = vieww
		}
		item.view.Resize(0, y, textw, texth)
		item.text.Resize()
	}
	w.height = len(w.items)
	w.changed = false
}
