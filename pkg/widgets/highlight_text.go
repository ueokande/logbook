package widgets

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

type point struct {
	x int
	y int
}

var styleHighlightCurrent = tcell.StyleDefault.Background(tcell.ColorYellow)

// HighlightText is a text widget with highlighted keyword
type HighlightText struct {
	highlights []int
	current    int
	keyword    []rune

	text views.Text
	views.WidgetWatchers
}

// Draw draws the HighlightText.
func (t *HighlightText) Draw() {
	t.text.Draw()
}

// Size returns the width and height of the HighlightText
func (t *HighlightText) Size() (int, int) {
	return t.text.Size()
}

// SetView sets the view for the HighlightText
func (t *HighlightText) SetView(view views.View) {
	t.text.SetView(view)
}

// HandleEvent implements a tcell.EventHandler
func (t *HighlightText) HandleEvent(ev tcell.Event) bool {
	return t.text.HandleEvent(ev)
}

// AppendLines appends the lines into the content
func (t *HighlightText) AppendLines(lines []string) {
	if len(lines) == 0 {
		return
	}
	text := t.text.Text()
	if len(text) > 0 {
		text += "\n"
	}
	for _, l := range lines {
		text += l + "\n"
	}
	text = text[:len(text)-1]
	t.text.SetText(text)

	t.resetHighlights()

	t.PostEventWidgetContent(t)
}

// ClearText clears current content and highlights
func (t *HighlightText) ClearText() {
	t.text.SetText("")
	t.keyword = nil
	t.current = -1
	t.highlights = nil
}

// SetKeyword sets the keyword to be highlighted in the content
func (t *HighlightText) SetKeyword(keyword string) {
	t.keyword = []rune(keyword)
	t.current = -1
	t.text.SetStyle(t.text.Style())
	if len(keyword) == 0 {
		return
	}

	t.resetHighlights()
	t.PostEventWidgetContent(t)
}

func (t *HighlightText) resetHighlights() {
	t.highlights = nil

	str := t.text.Text()
	keyword := string(t.keyword)
	var x int
	for len(keyword) > 0 {
		i := strings.Index(str, keyword)
		if i == -1 {
			break
		}
		start := len([]rune(str[:i])) + x
		t.highlights = append(t.highlights, start)
		str = str[i+len(keyword):]
		x += i + len(keyword)
	}

	for i, start := range t.highlights {
		style := t.text.Style().Reverse(true)
		if i == t.current {
			style = styleHighlightCurrent
		}
		for offset := range t.keyword {
			t.text.SetStyleAt(start+offset, style)
		}
	}
}

// Keyword returns the current keyword in the content
func (t *HighlightText) Keyword() string {
	return string(t.keyword)
}

// Resize is called when the View changes sizes.
func (t *HighlightText) Resize() {
	t.text.Resize()
}

// HighlightPos returns the position (x, y) of highlighted text in the content
func (t *HighlightText) HighlightPos(index int) (int, int) {
	if index < 0 || index >= len(t.highlights) {
		panic("index out of range")
	}

	start := t.highlights[index]
	var x, y int
	for i, c := range []rune(t.text.Text()) {
		if i == start {
			break
		}
		if c == '\n' {
			x = 0
			y++
			continue
		}
		x += runewidth.RuneWidth(c)
	}
	return x, y
}

// HighlightCount returns the count of the highlighted keywords
func (t *HighlightText) HighlightCount() int {
	return len(t.highlights)
}

// ActivateHighlight makes the highlighted keyword active (focused)
func (t *HighlightText) ActivateHighlight(index int) {
	if index < 0 || index > len(t.highlights) {
		panic("index out of range")
	}
	t.current = index
	t.resetHighlights()
	t.PostEventWidgetContent(t)
}

// CurrentHighlight returns the index of the current highlight. It returns -1 if no highlights are active.
func (t *HighlightText) CurrentHighlight() int {
	return t.current
}
