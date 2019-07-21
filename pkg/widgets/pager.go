package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// Pager is a Widget with the text and its view port.  It provides a scrollable
// view if the content size is larger than the actual view.
type Pager struct {
	view     views.View
	viewport views.ViewPort
	text     views.Text

	views.WidgetWatchers
}

// NewPager returns a new Pager
func NewPager() *Pager {
	w := &Pager{}
	w.text.SetView(&w.viewport)
	w.text.SetStyle(tcell.StyleDefault)
	return w
}

// AppendLine adds the line into the pager
func (w *Pager) AppendLine(line string) {
	text := w.text.Text()
	if len(text) > 0 {
		text += "\n"
	}
	text += line
	w.text.SetText(text)

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
	w.text.SetText("")

	width, height := w.text.Size()
	w.viewport.SetContentSize(width, height, true)
	w.viewport.ValidateView()
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
