package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleTabActive     = tcell.StyleDefault.Background(tcell.ColorSilver).Foreground(tcell.ColorWhite).Bold(true)
	styleTabInactive   = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	styleTabBackground = tcell.StyleDefault.Background(tcell.ColorWhite)
)

type tabsItem struct {
	name string
	text *views.Text
}

type Tabs struct {
	items    []tabsItem
	selected int

	views.BoxLayout
	views.WidgetWatchers
}

func NewTabs() *Tabs {
	w := &Tabs{
		selected: -1,
	}
	w.SetStyle(styleTabBackground)
	w.SetOrientation(views.Horizontal)
	return w
}

func (w *Tabs) AddTab(name string) {
	text := &views.Text{}
	text.SetText(" " + name + " ")
	text.SetStyle(styleTabInactive)

	w.AddWidget(text, 0)
	w.items = append(w.items, tabsItem{
		name: name,
		text: text,
	})
}

func (w *Tabs) TabCount() int {
	return len(w.items)
}

func (w *Tabs) Clear() {
	for _, item := range w.items {
		w.RemoveWidget(item.text)
	}
	w.items = nil
	w.selected = -1
}

func (w *Tabs) SelectAt(index int) {
	if index == w.selected {
		return
	}
	if w.selected >= 0 {
		item := w.items[w.selected]
		item.text.SetStyle(styleTabInactive)
	}
	if index < 0 || index >= len(w.items) {
		return
	}
	w.selected = index
	item := w.items[w.selected]
	item.text.SetStyle(styleTabActive)

	w.PostEventWidgetContent(w)

	ev := &EventItemSelected{
		Name:   item.name,
		Index:  index,
		widget: w,
	}
	ev.SetEventNow()

	w.PostEvent(ev)
}
