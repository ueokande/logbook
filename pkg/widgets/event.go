package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// EventItemSelected represents an event on the item selected
type EventItemSelected struct {
	// The name of the item
	Name string

	// The index of the item
	Index int

	widget views.Widget
	tcell.EventTime
}

// Widget returns a target widget of the event
func (e *EventItemSelected) Widget() views.Widget {
	return e.widget
}
