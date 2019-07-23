package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

// InputLine is a single-line input widget
type InputLine struct {
	view    views.View
	style   tcell.Style
	value   []rune
	prompt  []rune
	content string
	cursor  int

	views.WidgetWatchers
}

// NewInputLine returns new InputLine
func NewInputLine() *InputLine {
	return &InputLine{}
}

// SetPrompt sets the prompt of the input
func (w *InputLine) SetPrompt(prompt string) {
	w.prompt = []rune(prompt)
	w.PostEventWidgetContent(w)
}

// SetValue sets the value of the input
func (w *InputLine) SetValue(value string) {
	w.value = []rune(value)
	w.cursor = len(w.value)
	w.PostEventWidgetContent(w)
}

// SetStyle sets the style of the input
func (w *InputLine) SetStyle(style tcell.Style) {
	w.style = style
	w.PostEventWidgetContent(w)
}

// SetCursorAt sets the pos of the cursor
func (w *InputLine) SetCursorAt(pos int) {
	w.cursor = pos
	if w.cursor < 0 {
		w.cursor = 0
	} else if w.cursor > len(w.value) {
		w.cursor = len(w.value)
	}

	w.PostEventWidgetContent(w)
}

// Draw draws the input with the cursor
func (w *InputLine) Draw() {
	if w.view == nil {
		return
	}

	var x int
	for _, c := range w.prompt {
		w.view.SetContent(x, 0, c, nil, w.style)
		x += runewidth.RuneWidth(c)
	}
	for i, c := range w.value {
		style := w.style
		if i == w.cursor {
			style = style.Reverse(true)
		}
		w.view.SetContent(x, 0, c, nil, style)
		x += runewidth.RuneWidth(c)
	}
	if w.cursor == len(w.value) {
		w.view.SetContent(x, 0, ' ', nil, w.style.Reverse(true))
	}
}

// SetView sets the view
func (w *InputLine) SetView(view views.View) {
	w.view = view
	w.PostEventWidgetContent(w)
}

// Resize is called when our View changes sizes.
func (w *InputLine) Resize() {
	w.PostEventWidgetResize(w)
}

// HandleEvent handles events on tcell
func (w *InputLine) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyLeft:
			w.SetCursorAt(w.cursor - 1)
			return true
		case tcell.KeyRight:
			w.SetCursorAt(w.cursor + 1)
			return true
		case tcell.KeyDelete:
			copy(w.value[w.cursor:], w.value[w.cursor+1:])
			w.value = w.value[:len(w.value)-1]
			return true
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if w.cursor == 0 {
				break
			}
			copy(w.value[w.cursor-1:], w.value[w.cursor:])
			w.value = w.value[:len(w.value)-1]
			w.cursor--
			return true
		case tcell.KeyRune:
			runes := make([]rune, len(w.value)+1)
			copy(runes, w.value[:w.cursor])
			copy(runes[w.cursor+1:], w.value[w.cursor:])
			runes[w.cursor] = ev.Rune()
			w.value = runes
			w.cursor++
			return true
		}
	}
	return false
}

// Size returns the width and height in vertical line.
func (w *InputLine) Size() (int, int) {
	width1 := runewidth.StringWidth(string(w.prompt))
	width2 := runewidth.StringWidth(string(w.value))
	return width1 + width2 + 1, 1
}
