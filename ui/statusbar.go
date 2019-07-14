package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleStatusBarDefault = tcell.StyleDefault.Background(tcell.ColorGray)
	styleStatusBarLeft    = tcell.StyleDefault.Background(tcell.ColorSilver).Foreground(tcell.ColorWhite)
	styleStatusBarRight   = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
)

type StatusBar struct {
	views.TextBar
}

func NewStatusBar() *StatusBar {
	w := &StatusBar{}
	w.SetStyle(styleStatusBarDefault)
	return w
}

func (w *StatusBar) SetCenterStatus(text string) {
	w.SetCenter(text, styleStatusBarDefault)
}

func (w *StatusBar) SetLeftStatus(text string) {
	w.SetLeft(" "+text+" ", styleStatusBarLeft)
}

func (w *StatusBar) SetRightStatus(text string) {
	w.SetRight(" "+text+" ", styleStatusBarRight)
}
