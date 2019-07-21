package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type EventItemSelected struct {
	Name  string
	Index int

	widget views.Widget
	tcell.EventTime
}

func (e *EventItemSelected) Widget() views.Widget {
	return e.widget
}
