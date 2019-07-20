package main

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/pkg/k8s"
	"github.com/ueokande/logbook/pkg/ui"
	"github.com/ueokande/logbook/pkg/widgets"
	corev1 "k8s.io/api/core/v1"
)

var (
	StylePodActive  = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	StylePodError   = tcell.StyleDefault.Foreground(tcell.ColorRed)
	StylePodPending = tcell.StyleDefault.Foreground(tcell.ColorYellow)
)

type AppConfig struct {
	Cluster   string
	Namespace string
}

type App struct {
	client *k8s.Client

	mainLayout   *views.BoxLayout
	detailLayout *views.BoxLayout
	statusbar    *ui.StatusBar
	tabs         *widgets.Tabs
	podsView     *widgets.ListView
	line         *ui.VerticalLine
	pager        *ui.Pager

	namespace         string
	pods              []*corev1.Pod
	containers        []string
	selectedPod       int
	selectedContainer int
	podworker         *Worker
	logworker         *Worker
	follow            bool

	*views.Application
	views.BoxLayout
}

func NewApp(client *k8s.Client, config *AppConfig) *App {
	statusbar := ui.NewStatusBar()
	statusbar.SetContext(config.Cluster, config.Namespace)
	podsView := widgets.NewListView()
	line := ui.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := ui.NewPager()
	tabs := widgets.NewTabs()

	detailLayout := &views.BoxLayout{}
	detailLayout.SetOrientation(views.Vertical)
	detailLayout.AddWidget(tabs, 0)
	detailLayout.AddWidget(pager, 1)

	mainLayout := &views.BoxLayout{}
	mainLayout.SetOrientation(views.Horizontal)
	mainLayout.AddWidget(podsView, 0)
	mainLayout.AddWidget(line, 0)
	mainLayout.AddWidget(detailLayout, 1)

	app := &App{
		client: client,

		mainLayout:   mainLayout,
		detailLayout: detailLayout,
		statusbar:    statusbar,
		tabs:         tabs,
		podsView:     podsView,
		line:         line,
		pager:        pager,

		namespace: config.Namespace,
		logworker: NewWorker(context.Background()),
		podworker: NewWorker(context.Background()),

		Application: new(views.Application),
	}

	tabs.Watch(app)
	podsView.Watch(app)

	app.statusbar.SetMode(ui.ModeNormal)
	app.SetOrientation(views.Vertical)
	app.AddWidget(mainLayout, 1)
	app.AddWidget(statusbar, 0)
	app.SetRootWidget(app)

	return app
}

func (app *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *widgets.EventItemSelected:
		switch ev.Widget() {
		case app.tabs:
			app.HandleContainerSelected(ev.Name, ev.Index)
			return true
		case app.podsView:
			app.HandlePodSelected(ev.Name, ev.Index)
			return true
		}
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			app.Quit()
			return true
		case tcell.KeyCtrlP:
			app.SelectPrevPod()
			return true
		case tcell.KeyCtrlN:
			app.SelectNextPod()
			return true
		case tcell.KeyCtrlD:
			if app.follow {
				return false
			}
			app.pager.ScrollHalfPageDown()
			app.UpdateScrollStatus()
			return true
		case tcell.KeyCtrlU:
			if app.follow {
				return false
			}
			app.pager.ScrollHalfPageUp()
			app.UpdateScrollStatus()
			return true
		case tcell.KeyCtrlB:
			if app.follow {
				return false
			}
			app.pager.ScrollPageUp()
			app.UpdateScrollStatus()
			return true
		case tcell.KeyCtrlF:
			if app.follow {
				return false
			}
			app.pager.ScrollPageDown()
			app.UpdateScrollStatus()
			return true
		case tcell.KeyTab:
			app.SelectNextContainer()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				app.Quit()
				return true
			case 'k':
				if app.follow {
					return false
				}
				app.pager.ScrollUp()
				app.UpdateScrollStatus()
				return true
			case 'j':
				if app.follow {
					return false
				}
				app.pager.ScrollDown()
				app.UpdateScrollStatus()
				return true
			case 'g':
				if app.follow {
					return false
				}
				app.pager.ScrollToTop()
				app.UpdateScrollStatus()
				return true
			case 'G':
				if app.follow {
					return false
				}
				app.pager.ScrollToBottom()
				app.UpdateScrollStatus()
				return true
			case 'f':
				app.follow = !app.follow
				if app.follow {
					app.statusbar.SetMode(ui.ModeFollow)
					app.pager.ScrollToBottom()
					app.UpdateScrollStatus()
				} else {
					app.statusbar.SetMode(ui.ModeNormal)
				}
				return true
			}
		}
	}
	return app.BoxLayout.HandleEvent(ev)
}

func (app *App) HandleContainerSelected(name string, index int) {
	app.selectedContainer = index
	app.follow = false
	app.statusbar.SetMode(ui.ModeNormal)

	app.UpdateScrollStatus()

	pod := app.pods[app.selectedPod]
	container := app.containers[app.selectedContainer]
	app.StartTailLog(pod.Namespace, pod.Name, container)
}

func (app *App) HandlePodSelected(name string, index int) {
	pod := app.pods[app.selectedPod]
	app.containers = nil
	app.tabs.Clear()
	for _, c := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		app.tabs.AddTab(c.Name)
		app.containers = append(app.containers, c.Name)
	}
	app.SelectContainerAt(0)
}

func (app *App) SelectNextPod() {
	app.SelectPodAt(app.selectedPod + 1)
}

func (app *App) SelectPrevPod() {
	app.SelectPodAt(app.selectedPod - 1)
}

func (app *App) SelectPodAt(index int) {
	if index < 0 {
		index = 0
	}
	if index > len(app.pods)-1 {
		index = len(app.pods) - 1
	}
	app.selectedPod = index

	app.podsView.SelectAt(app.selectedPod)
}

func (app *App) UpdateScrollStatus() {
	y := app.pager.GetScrollYPosition()
	app.statusbar.SetScroll(int(y * 100))
}

func (app *App) SelectNextContainer() {
	index := app.selectedContainer + 1
	if index > len(app.containers)-1 {
		index = 0
	}
	app.SelectContainerAt(index)
}

func (app *App) SelectContainerAt(index int) {
	if index < 0 {
		index = 0
	}
	if index > len(app.pods)-1 {
		index = len(app.pods) - 1
	}
	app.tabs.SelectAt(index)
}

func (app *App) StartTailLog(namespace, pod, container string) {
	app.StopTailLog()

	app.pager.ClearText()
	app.logworker.Start(func(ctx context.Context) error {
		logs, err := app.client.WatchLogs(ctx, namespace, pod, container)
		if err != nil {
			return err
		}

		// make channel to guarantee line order of logs
		ch := make(chan string)
		defer close(ch)
		for log := range logs {
			app.PostFunc(func() {
				for line := range ch {
					app.pager.AppendLine(line)
					if app.follow {
						app.pager.ScrollToBottom()
					}
					app.UpdateScrollStatus()
					break
				}
			})
			select {
			case ch <- log:
			case <-ctx.Done():
				return nil
			}
		}
		return nil
	})
}

func (app *App) StopTailLog() {
	app.logworker.Stop()
	// TODO handle err
}

func (app *App) StartTailPods() {
	app.StopTailLog()
	app.podworker.Start(func(ctx context.Context) error {
		events, err := app.client.WatchPods(ctx, app.namespace)
		if err != nil {
			return err
		}
		for ev := range events {
			ev := ev
			app.PostFunc(func() {
				pod := ev.Pod
				switch ev.Type {
				case k8s.PodAdded:
					app.pods = append(app.pods, pod)
					switch k8s.GetPodStatus(pod) {
					case k8s.PodRunning, k8s.PodSucceeded:
						app.podsView.AddItem(pod.Name, StylePodActive)
					case k8s.PodPending, k8s.PodInitializing, k8s.PodTerminating:
						app.podsView.AddItem(pod.Name, StylePodPending)
					default:
						app.podsView.AddItem(pod.Name, StylePodError)
					}
					if len(app.pods) == 1 {
						app.podsView.SelectAt(0)
					}
				case k8s.PodModified:
					switch k8s.GetPodStatus(pod) {
					case k8s.PodRunning, k8s.PodSucceeded:
						app.podsView.SetStyle(pod.Name, StylePodActive)
					case k8s.PodPending, k8s.PodInitializing, k8s.PodTerminating:
						app.podsView.SetStyle(pod.Name, StylePodPending)
					default:
						app.podsView.SetStyle(pod.Name, StylePodError)
					}
				case k8s.PodDeleted:
					for i, p := range app.pods {
						if p.Name == pod.Name {
							app.pods = append(app.pods[:i], app.pods[i+1:]...)
							break
						}
					}
					app.podsView.DeleteItem(pod.Name)
				}
				app.statusbar.SetPodCount(len(app.pods))
			})

		}
		return nil
	})
}

func (app *App) StopTailPods() {
	app.podworker.Stop()
	// TODO handle err
}

func (app *App) Run(ctx context.Context) error {
	app.StartTailPods()
	return app.Application.Run()
}
