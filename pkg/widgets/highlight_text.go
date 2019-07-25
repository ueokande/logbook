package widgets

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// HighlightText is a text widget with highlighted keyword
type HighlightText struct {
	view      views.View
	style     tcell.Style
	runes     []rune
	widths    []int
	width     int
	height    int
	keyword   string
	highlight []bool

	views.WidgetWatchers
}

// Draw draws the HighlightText.
func (t *HighlightText) Draw() {
	v := t.view
	if v == nil {
		return
	}

	t.view.Fill(' ', t.style)

	var x, y int
	for i, c := range t.runes {
		if c == '\n' {
			x = 0
			y++
			continue
		}
		style := t.style
		if t.highlight[i] {
			style = style.Reverse(true)
		}
		v.SetContent(x, y, c, nil, style)
		x += t.widths[i]
	}
}

// Size returns the width and height of the HighlightText
func (t *HighlightText) Size() (int, int) {
	return t.width, t.height
}

// SetView sets the view for the HighlightText
func (t *HighlightText) SetView(view views.View) {
	t.view = view
}

// HandleEvent implements a tcell.EventHandler
func (t *HighlightText) HandleEvent(tcell.Event) bool {
	return false
}

// SetText sets the text for the HighlightText
func (t *HighlightText) SetText(s string) {
	t.width = 0
	t.runes = []rune(s)
	t.widths = make([]int, len(t.runes))

	var x, y int
	for i, r := range t.runes {
		t.widths[i] = runewidth.RuneWidth(r)
		x += t.widths[i]

		if r == '\n' {
			y++
			if x > t.width {
				t.width = x
				x = 0
			}
		}
	}
	if x > t.width {
		t.width = x
	}
	t.height = y

	t.SetKeyword(t.keyword)

	t.PostEventWidgetContent(t)
}

// Text returns the current text of the HighlightText
func (t *HighlightText) Text() string {
	return string(t.runes)
}

// SetKeyword sets the keyword to be highlighted in the content
func (t *HighlightText) SetKeyword(keyword string) {
	t.keyword = keyword
	t.highlight = make([]bool, len(t.runes))
	if len(keyword) == 0 {
		return
	}

	var starts []int // start positions to highlight
	str := string(t.runes)
	var x int
	for {
		i := strings.Index(str, keyword)
		if i == -1 {
			break
		}
		starts = append(starts, len([]rune(str[:i]))+x)
		str = str[i+len(keyword):]
		x += i + len(keyword)
	}

	l := len([]rune(keyword))
	for _, s := range starts {
		for i := s; i < s+l; i++ {
			t.highlight[i] = true
		}
	}
}

// SetStyle sets the style of the content
func (t *HighlightText) SetStyle(style tcell.Style) {
	t.style = style
	t.PostEventWidgetContent(t)
}

// Resize is called when the View changes sizes.
func (t *HighlightText) Resize() {
	t.PostEventWidgetResize(t)
}
