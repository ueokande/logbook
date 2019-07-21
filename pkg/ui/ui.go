package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/pkg/types"
	"github.com/ueokande/logbook/pkg/widgets"
)

var (
	stylePodActive  = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	stylePodError   = tcell.StyleDefault.Foreground(tcell.ColorRed)
	stylePodPending = tcell.StyleDefault.Foreground(tcell.ColorYellow)
)

type UIEventListener interface {
	OnQuit()
	OnContainerSelected(name string, index int)
	OnPodSelected(name string, index int)
}

type nopListener struct{}

func (l nopListener) OnContainerSelected(name string, index int) {}

func (l nopListener) OnPodSelected(name string, index int) {}

func (l nopListener) OnQuit() {}

type UI struct {
	pods       *widgets.ListView
	containers *widgets.Tabs
	pager      *widgets.Pager
	statusbar  *StatusBar

	follow            bool
	podIndex          int
	selectedContainer int

	listener UIEventListener

	views.BoxLayout
}

func NewUI() *UI {
	statusbar := NewStatusBar()
	pods := widgets.NewListView()
	line := widgets.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := widgets.NewPager()
	containers := widgets.NewTabs()

	detailLayout := &views.BoxLayout{}
	detailLayout.SetOrientation(views.Vertical)
	detailLayout.AddWidget(containers, 0)
	detailLayout.AddWidget(pager, 1)

	mainLayout := &views.BoxLayout{}
	mainLayout.SetOrientation(views.Horizontal)
	mainLayout.AddWidget(pods, 0)
	mainLayout.AddWidget(line, 0)
	mainLayout.AddWidget(detailLayout, 1)

	ui := &UI{
		pods:       pods,
		containers: containers,
		pager:      pager,
		statusbar:  statusbar,
		listener:   &nopListener{},
	}

	ui.SetOrientation(views.Vertical)
	ui.AddWidget(mainLayout, 1)
	ui.AddWidget(statusbar, 0)

	pods.Watch(ui)
	containers.Watch(ui)

	return ui
}

func (ui *UI) WatchUIEvents(l UIEventListener) {
	ui.listener = l
}

func (ui *UI) AddPod(name string, status types.PodStatus) {
	ui.pods.AddItem(name, podStatusStyle(status))
	ui.statusbar.SetPodCount(ui.pods.ItemCount())
}

func (ui *UI) DeletePod(name string) {
	ui.pods.DeleteItem(name)
	ui.statusbar.SetPodCount(ui.pods.ItemCount())
}

func (ui *UI) SetPodStatus(name string, status types.PodStatus) {
	ui.pods.SetStyle(name, podStatusStyle(status))
}

func (ui *UI) AddContainer(name string) {
	ui.containers.AddTab(name)
}

func (ui *UI) ClearContainers() {
	ui.containers.Clear()
}

func (ui *UI) SetContext(cluster, context string) {
	ui.statusbar.SetContext(cluster, context)
}

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *widgets.EventItemSelected:
		switch ev.Widget() {
		case ui.containers:
			ui.listener.OnContainerSelected(ev.Name, ev.Index)
			return true
		case ui.pods:
			ui.listener.OnPodSelected(ev.Name, ev.Index)
			return true
		}
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			ui.selectPrevPod()
			return true
		case tcell.KeyCtrlN:
			ui.selectNextPod()
			return true
		case tcell.KeyCtrlD:
			ui.scrollHalfPageDown()
			return true
		case tcell.KeyCtrlU:
			ui.scrollHalfPageUp()
		case tcell.KeyCtrlB:
			ui.scrollPageDown()
		case tcell.KeyCtrlF:
			ui.scrollPageUp()
			return true
		case tcell.KeyTab:
			ui.selectNextContainer()
			return true
		case tcell.KeyCtrlC:
			ui.listener.OnQuit()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'k':
				ui.scrollUp()
				return true
			case 'j':
				ui.scrollDown()
				return true
			case 'g':
				ui.scrollToTop()
				return true
			case 'G':
				ui.scrollToBottom()
				return true
			case 'f':
				ui.toggleFollowMode()
				return true
			case 'q':
				ui.listener.OnQuit()
				return true
			}
		}
	}
	return false
}

func (ui *UI) AddPagerText(line string) {
	ui.pager.AppendLine(line)
	if ui.follow {
		ui.pager.ScrollToBottom()
	}
	ui.updateScrollStatus()
}

func (ui *UI) ClearPager() {
	ui.pager.ClearText()
	ui.updateScrollStatus()
	ui.DisableFollowMode()
}

func (ui *UI) SetStatusMode(mode Mode) {
	ui.statusbar.SetMode(mode)
}

func (ui *UI) scrollDown() {
	if ui.follow {
		return
	}
	ui.pager.ScrollDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollUp() {
	if ui.follow {
		return
	}
	ui.pager.ScrollUp()
	ui.updateScrollStatus()
}

func (ui *UI) scrollPageDown() {
	if ui.follow {
		return
	}
	ui.pager.ScrollPageDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollPageUp() {
	if ui.follow {
		return
	}
	ui.pager.ScrollPageUp()
	ui.updateScrollStatus()
}

func (ui *UI) scrollHalfPageDown() {
	if ui.follow {
		return
	}
	ui.pager.ScrollHalfPageDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollToTop() {
	if ui.follow {
		return
	}
	ui.pager.ScrollToTop()
	ui.updateScrollStatus()
}

func (ui *UI) scrollToBottom() {
	if ui.follow {
		return
	}
	ui.pager.ScrollToBottom()
	ui.updateScrollStatus()
}

func (ui *UI) scrollHalfPageUp() {
	if ui.follow {
		return
	}
	ui.pager.ScrollHalfPageUp()
	ui.updateScrollStatus()
}

func (ui *UI) toggleFollowMode() {
	if ui.follow {
		ui.DisableFollowMode()
	} else {
		ui.EnableFollowMode()
	}
}

func (ui *UI) EnableFollowMode() {
	ui.follow = true
	ui.statusbar.SetMode(ModeFollow)
	ui.pager.ScrollToBottom()
	ui.updateScrollStatus()
}

func (ui *UI) DisableFollowMode() {
	ui.follow = false
	ui.statusbar.SetMode(ModeNormal)
}

func (ui *UI) selectNextPod() {
	count := ui.pods.ItemCount()
	if ui.podIndex+1 >= count {
		ui.SelectPodAt(0)
	} else {
		ui.SelectPodAt(ui.podIndex + 1)
	}
}

func (ui *UI) selectPrevPod() {
	if ui.podIndex == 0 {
		count := ui.pods.ItemCount()
		ui.SelectPodAt(count - 1)
	} else {
		ui.SelectPodAt(ui.podIndex - 1)
	}
}

func (ui *UI) SelectPodAt(index int) {
	ui.podIndex = index
	ui.pods.SelectAt(index)
}

func (ui *UI) selectNextContainer() {
	index := ui.selectedContainer + 1
	if index >= ui.containers.TabCount() {
		index = 0
	}
	ui.selectedContainer = index
	ui.containers.SelectAt(index)
}

func (ui *UI) SelectContainerAt(index int) {
	ui.containers.SelectAt(index)
}

func (ui *UI) updateScrollStatus() {
	y := ui.pager.GetScrollYPosition()
	ui.statusbar.SetScroll(int(y * 100))
}

func podStatusStyle(status types.PodStatus) tcell.Style {
	switch status {
	case types.PodRunning, types.PodSucceeded:
		return stylePodActive
	case types.PodPending, types.PodInitializing, types.PodTerminating:
		return stylePodPending
	default:
		return stylePodError
	}
}
