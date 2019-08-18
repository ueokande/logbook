package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// Pager is a Widget with the text and its view port.  It provides a scrollable
// view if the content size is larger than the actual view.
type Pager struct {
	view      views.View
	viewport  views.ViewPort
	text      HighlightText
	highlight string

	views.WidgetWatchers
}

// NewPager returns a new Pager
func NewPager() *Pager {
	w := &Pager{}
	w.text.SetView(&w.viewport)
	return w
}

// AppendLines adds the lines into the pager
func (w *Pager) AppendLines(lines []string) {
	w.text.AppendLines(lines)

	width, height := w.text.Size()
	w.viewport.SetContentSize(width, height, true)
	w.viewport.ValidateView()
}

// ScrollDown scrolls down by one line on the pager.
func (w *Pager) ScrollDown() {
	w.viewport.ScrollDown(1)
}

// ScrollUp scrolls up by one line on the pager
func (w *Pager) ScrollUp() {
	w.viewport.ScrollUp(1)
}

// ScrollHalfPageDown scrolls down by half-height of the screen.
func (w *Pager) ScrollHalfPageDown() {
	_, vh := w.view.Size()
	w.viewport.ScrollDown(vh / 2)
}

// ScrollHalfPageUp scrolls up by half-height of the screen.
func (w *Pager) ScrollHalfPageUp() {
	_, vh := w.view.Size()
	w.viewport.ScrollUp(vh / 2)
}

// ScrollHalfPageLeft scrolls left by half-width of the screen
func (w *Pager) ScrollHalfPageLeft() {
	vw, _ := w.view.Size()
	w.viewport.ScrollLeft(vw / 2)
}

// ScrollHalfPageRight scrolls right by half-width of the screen.
func (w *Pager) ScrollHalfPageRight() {
	vw, _ := w.view.Size()
	w.viewport.ScrollRight(vw / 2)
}

// ScrollPageDown scrolls down by the height of the screen.
func (w *Pager) ScrollPageDown() {
	_, vh := w.view.Size()
	w.viewport.ScrollDown(vh)
}

// ScrollPageUp scrolls up by the height of the screen.
func (w *Pager) ScrollPageUp() {
	_, vh := w.view.Size()
	w.viewport.ScrollUp(vh)
}

// ScrollToTop scrolls to the top of the content
func (w *Pager) ScrollToTop() {
	_, h := w.text.Size()
	w.viewport.ScrollUp(h)
}

// ScrollToBottom scrolls to the bottom of the content.
func (w *Pager) ScrollToBottom() {
	_, h := w.text.Size()
	w.viewport.ScrollDown(h)
}

// GetScrollYPosition returns vertical position of the scroll on the pager.
// Its range is 0.0 to 1.0.
func (w *Pager) GetScrollYPosition() float64 {
	_, contenth := w.viewport.GetContentSize()
	_, viewh := w.viewport.Size()
	_, y, _, _ := w.viewport.GetVisible()
	return float64(y) / float64(contenth-viewh)
}

// ClearText clears current content on the pager.
func (w *Pager) ClearText() {
	w.text.ClearText()

	width, height := w.text.Size()
	w.viewport.SetContentSize(width, height, true)
	w.viewport.ValidateView()
}

// SetKeyword sets the keyword to be highlighted in the pager
func (w *Pager) SetKeyword(keyword string) {
	w.text.SetKeyword(keyword)
	w.PostEventWidgetContent(w)
}

// Keyword returns the current keyword in the content
func (w *Pager) Keyword() string {
	return w.text.Keyword()
}

// FindNext finds next keyword in the content.  It returns true if the keyword found.
func (w *Pager) FindNext() bool {
	count := w.text.HighlightCount()
	if count == 0 {
		return false
	}
	next := w.text.CurrentHighlight() + 1
	if next >= count {
		next = 0
	}
	w.text.ActivateHighlight(next)
	x, y := w.text.HighlightPos(next)
	w.viewport.Center(x, y)
	w.PostEventWidgetContent(w)
	return true
}

// FindPrev finds previous keyword in the content.  It returns true if the keyword found.
func (w *Pager) FindPrev() bool {
	count := w.text.HighlightCount()
	if count == 0 {
		return false
	}
	next := w.text.CurrentHighlight() - 1
	if next < 0 {
		next = count - 1
	}
	w.text.ActivateHighlight(next)
	x, y := w.text.HighlightPos(next)
	w.viewport.Center(x, y)
	w.PostEventWidgetContent(w)
	return true
}

// Draw draws the Pager
func (w *Pager) Draw() {
	if w.view == nil {
		return
	}
	w.text.Draw()
}

// Resize is called when our View changes sizes.
func (w *Pager) Resize() {
	width, height := w.view.Size()
	w.viewport.Resize(0, 0, width, height)
	w.viewport.ValidateView()
}

// HandleEvent handles events on tcell.
func (w *Pager) HandleEvent(ev tcell.Event) bool {
	return false
}

// SetView sets the View object used for the pager
func (w *Pager) SetView(view views.View) {
	w.view = view
	w.viewport.SetView(view)
	if view == nil {
		return
	}
	w.Resize()
}

// Size returns the width and height in vertical line.
func (w *Pager) Size() (int, int) {
	width, height := w.view.Size()
	if width > 2 {
		width = 2
	}
	if height > 2 {
		height = 2
	}
	return width, height
}
