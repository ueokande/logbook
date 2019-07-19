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

type Tabs struct {
	texts    []*views.Text
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

func (w *Tabs) AddTab(label string) {
	text := &views.Text{}
	text.SetText(" " + label + " ")
	text.SetStyle(styleTabInactive)

	w.AddWidget(text, 0)
	w.texts = append(w.texts, text)
}

func (w *Tabs) Clear() {
	for _, t := range w.texts {
		w.RemoveWidget(t)
	}
	w.texts = nil
	w.selected = -1
}

func (w *Tabs) SelectAt(index int) {
	if index == w.selected {
		return
	}
	if w.selected >= 0 {
		text := w.texts[w.selected]
		text.SetStyle(styleTabInactive)
	}
	if index < 0 || index >= len(w.texts) {
		return
	}
	w.selected = index
	text := w.texts[w.selected]
	text.SetStyle(styleTabActive)

	w.PostEventWidgetContent(w)

	ev := &EventItemSelected{
		Name:   text.Text(),
		Index:  index,
		widget: w,
	}
	ev.SetEventNow()

	w.PostEvent(ev)
}
