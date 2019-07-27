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
	keyword    string

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

// SetText sets the text for the HighlightText
func (t *HighlightText) SetText(s string) {
	t.text.SetText(s)
	t.SetKeyword(t.keyword)
	t.PostEventWidgetContent(t)
}

// Text returns the current text of the HighlightText
func (t *HighlightText) Text() string {
	return string(t.text.Text())
}

// SetKeyword sets the keyword to be highlighted in the content
func (t *HighlightText) SetKeyword(keyword string) {
	t.keyword = keyword
	t.current = -1
	t.text.SetStyle(t.text.Style())
	t.highlights = nil
	if len(keyword) == 0 {
		return
	}

	style := t.text.Style().Reverse(true)
	str := t.text.Text()
	var x int
	for {
		i := strings.Index(str, keyword)
		if i == -1 {
			break
		}
		start := len([]rune(str[:i])) + x
		t.highlights = append(t.highlights, start)
		str = str[i+len(keyword):]
		x += i + len(keyword)
	}

	runes := []rune(keyword)
	for _, start := range t.highlights {
		for i := range runes {
			t.text.SetStyleAt(start+i, style)
		}
	}
}

// Keyword returns the current keyword in the content
func (t *HighlightText) Keyword() string {
	return t.keyword
}

// SetStyle sets the style of the content
func (t *HighlightText) SetStyle(style tcell.Style) {
	t.text.SetStyle(style)
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
	t.SetKeyword(t.keyword)
	t.current = index

	pos := t.highlights[index]
	for i := pos; i < pos+len(t.keyword); i++ {
		t.text.SetStyleAt(i, styleHighlightCurrent)
	}
	t.PostEventWidgetContent(t)
}

// CurrentHighlight returns the index of the current highlight. It returns -1 if no highlights are active.
func (t *HighlightText) CurrentHighlight() int {
	return t.current
}
