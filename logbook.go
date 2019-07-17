package main

import (
	"bufio"
	"context"
	"net/url"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/pkg/errors"
	"github.com/ueokande/logbook/ui"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
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
	clientset *kubernetes.Clientset

	mainLayout   *views.BoxLayout
	detailLayout *views.BoxLayout
	statusbar    *ui.StatusBar
	tabs         *ui.Tabs
	podsView     *ui.ListView
	line         *ui.VerticalLine
	pager        *ui.Pager
	pagerEnabled bool

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

func NewApp(clientset *kubernetes.Clientset, config *AppConfig) *App {
	statusbar := ui.NewStatusBar()
	statusbar.SetContext(config.Cluster, config.Namespace)
	podsView := ui.NewListView()
	line := ui.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := ui.NewPager()
	tabs := ui.NewTabs()

	detailLayout := &views.BoxLayout{}
	detailLayout.SetOrientation(views.Vertical)
	detailLayout.AddWidget(tabs, 0)
	detailLayout.AddWidget(pager, 1)

	mainLayout := &views.BoxLayout{}
	mainLayout.SetOrientation(views.Horizontal)
	mainLayout.AddWidget(podsView, 0)

	app := &App{
		clientset: clientset,

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

	app.statusbar.SetMode(ui.ModeNormal)
	app.SetOrientation(views.Vertical)
	app.AddWidget(mainLayout, 1)
	app.AddWidget(statusbar, 0)
	app.SetRootWidget(app)

	return app
}

func (app *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEnter:
			app.ShowPager()
			return true
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
				if app.pagerEnabled {
					app.HidePager()
				} else {
					app.Quit()
				}
				return true
			case 'k':
				if app.follow {
					return false
				}
				if app.pagerEnabled {
					app.pager.ScrollUp()
					app.UpdateScrollStatus()
				} else {
					app.SelectPrevPod()
				}
				return true
			case 'j':
				if app.follow {
					return false
				}
				if app.pagerEnabled {
					app.pager.ScrollDown()
					app.UpdateScrollStatus()
				} else {
					app.SelectNextPod()
				}
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
				if app.pagerEnabled {
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
				return false
			}
		}
	}
	return app.BoxLayout.HandleEvent(ev)
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

	if app.pagerEnabled {
		pod := app.pods[app.selectedPod]
		app.containers = nil
		app.tabs.Clear()
		for _, c := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
			app.tabs.AddTab(c.Name)
			app.containers = append(app.containers, c.Name)
		}
		app.SelectContainerAt(0)
	}
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
	app.selectedContainer = index
	app.follow = false
	app.statusbar.SetMode(ui.ModeNormal)

	app.tabs.SelectAt(index)
	app.UpdateScrollStatus()

	pod := app.pods[app.selectedPod]
	container := app.containers[app.selectedContainer]
	app.StartTailLog(pod.Namespace, pod.Name, container)
}

func (app *App) ShowPager() {
	if app.pagerEnabled {
		return
	}
	app.mainLayout.AddWidget(app.line, 0)
	app.mainLayout.AddWidget(app.detailLayout, 1)
	app.pagerEnabled = true

	app.SelectPodAt(app.selectedPod)
}

func (app *App) HidePager() {
	app.mainLayout.RemoveWidget(app.line)
	app.mainLayout.RemoveWidget(app.detailLayout)
	app.pagerEnabled = false
	app.follow = false
	app.statusbar.SetMode(ui.ModeNormal)

	app.StopTailLog()
}

func (app *App) StartTailLog(namespace, pod, container string) {
	app.StopTailLog()

	app.pager.ClearText()
	app.logworker.Start(func(ctx context.Context) error {
		opts := &corev1.PodLogOptions{
			Container: container,
			Follow:    true,
		}
		req := app.clientset.CoreV1().Pods(namespace).GetLogs(pod, opts)
		req.Context(ctx)
		r, err := req.Stream()
		if err != nil {
			return err
		}
		defer r.Close()

		s := bufio.NewScanner(r)

		// make channel to guarantee line order of logs
		ch := make(chan string)
		defer close(ch)
		for s.Scan() {
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
			case ch <- s.Text():
			case <-ctx.Done():
				return nil
			}
		}
		return s.Err()
	})
}

func (app *App) StopTailLog() {
	err := app.logworker.Stop()
	err = errors.Cause(err)
	if err == context.Canceled {
		return
	}
	if uerr, ok := err.(*url.Error); ok && uerr.Err == context.Canceled {
		return
	}
	if err != nil {
		// TODO handle err
	}
}

func (app *App) StartTailPods() {
	app.StopTailLog()
	app.podworker.Start(func(ctx context.Context) error {
		result, err := app.clientset.CoreV1().Pods(app.namespace).Watch(metav1.ListOptions{})
		if err != nil {
			return err
		}

		for ev := range result.ResultChan() {
			ev := ev
			pod, ok := ev.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			app.PostFunc(func() {
				switch ev.Type {
				case watch.Added:
					app.pods = append(app.pods, pod)
					switch GetPodStatus(pod) {
					case PodRunning, PodSucceeded:
						app.podsView.AddItem(pod.Name, StylePodActive)
					case PodPending, PodInitializing, PodTerminating:
						app.podsView.AddItem(pod.Name, StylePodPending)
					default:
						app.podsView.AddItem(pod.Name, StylePodError)
					}
					if len(app.pods) == 1 {
						app.podsView.SelectAt(0)
					}
				case watch.Modified:
					switch GetPodStatus(pod) {
					case PodRunning, PodSucceeded:
						app.podsView.SetStyle(pod.Name, StylePodActive)
					case PodPending, PodInitializing, PodTerminating:
						app.podsView.SetStyle(pod.Name, StylePodPending)
					default:
						app.podsView.SetStyle(pod.Name, StylePodError)
					}
				case watch.Deleted:
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
	err := app.podworker.Stop()
	err = errors.Cause(err)
	if err == context.Canceled {
		return
	}
	if uerr, ok := err.(*url.Error); ok && uerr.Err == context.Canceled {
		return
	}
	if err != nil {
		// TODO handle err
	}
}

func (app *App) Run(ctx context.Context) error {
	app.StartTailPods()
	return app.Application.Run()
}
