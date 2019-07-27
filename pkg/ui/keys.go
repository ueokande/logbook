package ui

import (
	"github.com/gdamore/tcell"
)

func (ui *UI) handleEventKey(ev *tcell.EventKey) bool {
	var handles []func(ev *tcell.EventKey) bool
	switch ui.mode {
	case ModeNormal:
		handles = []func(ev *tcell.EventKey) bool{
			ui.handleKeyInputFind,
			ui.handleKeySelectContainer,
			ui.handleKeyToggleFollowMode,
			ui.handleKeyScroll,
			ui.handleKeyFind,
			ui.handleKeyQuit,
		}
	case ModeFollow:
		handles = []func(ev *tcell.EventKey) bool{
			ui.handleKeySelectContainer,
			ui.handleKeyToggleFollowMode,
			ui.handleKeyQuit,
		}
	case ModeInputFind:
		handles = []func(ev *tcell.EventKey) bool{
			ui.handleEventKeyInput,
			ui.handleKeyQuit,
		}
	}

	for _, h := range handles {
		if h(ev) {
			return true
		}
	}

	return false
}

func (ui *UI) handleKeyQuit(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlC:
		ui.listener.OnQuit()
		return true
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'q':
			ui.listener.OnQuit()
			return true
		}
	}
	return false
}

func (ui *UI) handleKeyInputFind(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case '/':
			ui.enterFindInputMode()
			return true
		}
	}
	return false
}

func (ui *UI) handleKeyToggleFollowMode(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'f':
			ui.toggleFollowMode()
			return true
		}
	}
	return false
}

func (ui *UI) handleKeySelectContainer(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlP:
		ui.pods.SelectPrev()
		return true
	case tcell.KeyCtrlN:
		ui.pods.SelectNext()
		return true
	case tcell.KeyTab:
		ui.containers.SelectNext()
		return true
	}
	return false
}

func (ui *UI) handleKeyScroll(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlD:
		ui.scrollHalfPageDown()
		return true
	case tcell.KeyCtrlU:
		ui.scrollHalfPageUp()
		return true
	case tcell.KeyCtrlB:
		ui.scrollPageDown()
		return true
	case tcell.KeyCtrlF:
		ui.scrollPageUp()
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
		}
	}
	return false
}

func (ui *UI) handleKeyFind(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'n':
			ui.pager.FindNext()
			return true
		case 'N':
			ui.pager.FindPrev()
			return true
		}
	}
	return false
}

func (ui *UI) handleEventKeyInput(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyEnter:
		ui.startFind()
		return true
	case tcell.KeyEscape:
		ui.startFind()
		return true
	}
	return ui.input.HandleEvent(ev)
}
