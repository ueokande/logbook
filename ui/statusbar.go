package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleStatusBarContext = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSilver)
	styleStatusBarPods    = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	styleStatusBarScroll  = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
)

type StatusBar struct {
	pods    *views.Text
	context *views.Text
	scroll  *views.Text
	views.BoxLayout
}

func NewStatusBar() *StatusBar {
	pods := &views.Text{}
	pods.SetStyle(styleStatusBarPods)
	context := &views.Text{}
	context.SetAlignment(views.AlignMiddle)
	context.SetStyle(styleStatusBarContext)
	scroll := &views.Text{}
	scroll.SetStyle(styleStatusBarScroll)

	w := &StatusBar{
		pods:    pods,
		context: context,
		scroll:  scroll,
	}
	w.AddWidget(pods, 0)
	w.AddWidget(context, 1)
	w.AddWidget(scroll, 0)
	return w
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
