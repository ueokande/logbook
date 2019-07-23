package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/pkg/types"
	"github.com/ueokande/logbook/pkg/widgets"
)

// Mode represents a mode on UI
type Mode int

// UI mode
const (
	ModeNormal    Mode = iota // Normal mode
	ModeFollow                // Follow mode
	ModeInputFind             // Input find mode
)

var (
	stylePodActive  = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	stylePodError   = tcell.StyleDefault.Foreground(tcell.ColorRed)
	stylePodPending = tcell.StyleDefault.Foreground(tcell.ColorYellow)
)

// EventListener is a listener interface for UI events
type EventListener interface {
	// OnQuit is invoked on the quit is required
	OnQuit()

	// OnPodSelected is invoked when the selected pod is changed
	OnPodSelected(name string, index int)

	// OnContainerSelected is invoked when the selected container is changed
	OnContainerSelected(name string, index int)
}

type nopListener struct{}

func (l nopListener) OnContainerSelected(name string, index int) {}

func (l nopListener) OnPodSelected(name string, index int) {}

func (l nopListener) OnQuit() {}

// UI is an user interface for the logbook
type UI struct {
	input      *widgets.InputLine
	pods       *widgets.ListView
	containers *widgets.Tabs
	pager      *widgets.Pager
	statusbar  *StatusBar

	mode     Mode
	keyword  string
	listener EventListener

	views.BoxLayout
}

// NewUI returns new UI
func NewUI() *UI {
	statusbar := NewStatusBar()
	input := widgets.NewInputLine()
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
		input:      input,
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

// WatchUIEvents registers an EventListener for the UI
func (ui *UI) WatchUIEvents(l EventListener) {
	ui.listener = l
}

// AddPod adds a pod by the name and its status to the list view.
func (ui *UI) AddPod(name string, status types.PodStatus) {
	ui.pods.AddItem(name, podStatusStyle(status))
	ui.statusbar.SetPodCount(ui.pods.ItemCount())
}

// DeletePod deletes pod by the name on the list view.
func (ui *UI) DeletePod(name string) {
	ui.pods.DeleteItem(name)
	ui.statusbar.SetPodCount(ui.pods.ItemCount())
}

// SetPodStatus updates the pod status by name to the status
func (ui *UI) SetPodStatus(name string, status types.PodStatus) {
	ui.pods.SetStyle(name, podStatusStyle(status))
}

// AddContainer adds container by the name into the tabs
func (ui *UI) AddContainer(name string) {
	ui.containers.AddTab(name)
}

// ClearContainers clears containers in the tabs
func (ui *UI) ClearContainers() {
	ui.containers.Clear()
}

// SetContext sets kubenetes context (the cluster name and the namespace)
func (ui *UI) SetContext(cluster, namespace string) {
	ui.statusbar.SetContext(cluster, namespace)
}

// HandleEvent handles events on tcell
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
		return ui.handleEventKey(ev)
	}
	return false
}

// AddPagerText adds text line into the pager
func (ui *UI) AddPagerText(line string) {
	ui.pager.AppendLine(line)
	if ui.mode == ModeFollow {
		ui.pager.ScrollToBottom()
	}
	ui.updateScrollStatus()
}

// ClearPager clears the pager
func (ui *UI) ClearPager() {
	ui.pager.ClearText()
	ui.updateScrollStatus()
	ui.DisableFollowMode()
}

// SetStatusMode sets the mode in the status bar
func (ui *UI) SetStatusMode(mode Mode) {
	ui.statusbar.SetMode(mode)
}

func (ui *UI) scrollDown() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollUp() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollUp()
	ui.updateScrollStatus()
}

func (ui *UI) scrollPageDown() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollPageDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollPageUp() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollPageUp()
	ui.updateScrollStatus()
}

func (ui *UI) scrollHalfPageDown() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollHalfPageDown()
	ui.updateScrollStatus()
}

func (ui *UI) scrollToTop() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollToTop()
	ui.updateScrollStatus()
}

func (ui *UI) scrollToBottom() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollToBottom()
	ui.updateScrollStatus()
}

func (ui *UI) scrollHalfPageUp() {
	if ui.mode == ModeFollow {
		return
	}
	ui.pager.ScrollHalfPageUp()
	ui.updateScrollStatus()
}

func (ui *UI) toggleFollowMode() {
	if ui.mode == ModeFollow {
		ui.DisableFollowMode()
	} else {
		ui.EnableFollowMode()
	}
}

// EnableFollowMode enabled follow mode on the pager
func (ui *UI) EnableFollowMode() {
	ui.mode = ModeFollow
	ui.statusbar.SetMode(ModeFollow)
	ui.pager.ScrollToBottom()
	ui.updateScrollStatus()
}

// DisableFollowMode disables follow mode on the pager
func (ui *UI) DisableFollowMode() {
	ui.mode = ModeNormal
	ui.statusbar.SetMode(ModeNormal)
}

// SelectPodAt selects a pod by the index
func (ui *UI) SelectPodAt(index int) {
	ui.pods.SelectAt(index)
}

// SelectContainerAt selects a container by the index
func (ui *UI) SelectContainerAt(index int) {
	ui.containers.SelectAt(index)
}

func (ui *UI) updateScrollStatus() {
	y := ui.pager.GetScrollYPosition()
	ui.statusbar.SetScroll(int(y * 100))
}

func (ui *UI) enterFindInputMode() {
	ui.input.SetPrompt("/")
	ui.mode = ModeInputFind
	ui.RemoveWidget(ui.statusbar)
	ui.AddWidget(ui.input, 0)
}

func (ui *UI) cancelInput() {
	ui.mode = ModeNormal
	ui.RemoveWidget(ui.statusbar)
	ui.AddWidget(ui.input, 0)
}

func (ui *UI) startFind() {
	ui.keyword = ui.input.Value()
	ui.mode = ModeNormal
	ui.AddWidget(ui.statusbar, 0)
	ui.RemoveWidget(ui.input)
	ui.findNext()
}

func (ui *UI) findNext() {
	panic("TODO: implement")
}

func (ui *UI) findPrev() {
	panic("TODO: implement")
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
