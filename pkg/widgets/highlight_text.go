package widgets

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type point struct {
	x int
	y int
}

var styleHighlightCurrent = tcell.StyleDefault.Background(tcell.ColorYellow)

// HighlightText is a text widget with highlighted keyword
type HighlightText struct {
	keyword string

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
	t.text.SetStyle(t.text.Style())
	if len(keyword) == 0 {
		return
	}

	style := t.text.Style().Reverse(true)
	str := t.text.Text()
	keywordRunes := []rune(keyword)
	var x int
	for {
		i := strings.Index(str, keyword)
		if i == -1 {
			break
		}
		start := len([]rune(str[:i])) + x
		for i := range keywordRunes {
			t.text.SetStyleAt(start+i, style)
		}
		str = str[i+len(keyword):]
		x += i + len(keyword)
	}
}

// SetStyle sets the style of the content
func (t *HighlightText) SetStyle(style tcell.Style) {
	t.text.SetStyle(style)
}

// Resize is called when the View changes sizes.
func (t *HighlightText) Resize() {
	t.text.Resize()
}
