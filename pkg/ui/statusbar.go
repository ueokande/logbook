package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeFollow
)

var (
	styleStatusBarModeNormal = tcell.StyleDefault.Background(tcell.ColorYellowGreen).Foreground(tcell.ColorDarkGreen).Bold(true)
	styleStatusBarModeFollow = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite).Bold(true)
	styleStatusBarContext    = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSilver)
	styleStatusBarPods       = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	styleStatusBarScroll     = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
)

type StatusBar struct {
	mode    *views.Text
	pods    *views.Text
	context *views.Text
	scroll  *views.Text
	views.BoxLayout
}

func NewStatusBar() *StatusBar {
	mode := &views.Text{}
	mode.SetStyle(styleStatusBarPods)
	pods := &views.Text{}
	pods.SetStyle(styleStatusBarPods)
	context := &views.Text{}
	context.SetAlignment(views.AlignMiddle)
	context.SetStyle(styleStatusBarContext)
	scroll := &views.Text{}
	scroll.SetStyle(styleStatusBarScroll)

	w := &StatusBar{
		mode:    mode,
		pods:    pods,
		context: context,
		scroll:  scroll,
	}
	w.AddWidget(mode, 0)
	w.AddWidget(pods, 0)
	w.AddWidget(context, 1)
	w.AddWidget(scroll, 0)
	return w
}

func (w *StatusBar) SetMode(mode Mode) {
	switch mode {
	case ModeNormal:
		w.mode.SetText(" NORMAL ")
		w.mode.SetStyle(styleStatusBarModeNormal)
	case ModeFollow:
		w.mode.SetText(" FOLLOW ")
		w.mode.SetStyle(styleStatusBarModeFollow)
	default:
		panic("unsupported mode")
	}
}

func (w *StatusBar) SetContext(cluster, namespace string) {
	w.context.SetText(fmt.Sprintf("%s/%s", cluster, namespace))
}

func (w *StatusBar) SetPodCount(count int) {
	w.pods.SetText(fmt.Sprintf(" %d Pods ", count))
}

func (w *StatusBar) SetScroll(percent int) {
	w.scroll.SetText(fmt.Sprintf(" %d%% ", percent))
}
