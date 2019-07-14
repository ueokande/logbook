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
	Namespace string
}

type App struct {
	clientset *kubernetes.Clientset

	mainLayout   *views.BoxLayout
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

	*views.Application
	views.BoxLayout
}

func NewApp(clientset *kubernetes.Clientset, config *AppConfig) *App {
	podsView := ui.NewListView()
	line := ui.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := ui.NewPager()
	tabs := ui.NewTabs()

	mainLayout := &views.BoxLayout{}
	mainLayout.SetOrientation(views.Vertical)
	mainLayout.AddWidget(tabs, 0)
	mainLayout.AddWidget(pager, 1)

	app := &App{
		clientset: clientset,

		mainLayout: mainLayout,
		tabs:       tabs,
		podsView:   podsView,
		line:       line,
		pager:      pager,

		namespace: config.Namespace,
		logworker: NewWorker(context.Background()),
		podworker: NewWorker(context.Background()),

		Application: new(views.Application),
	}

	app.SetOrientation(views.Horizontal)
	app.AddWidget(podsView, 0.1)
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
			app.pager.ScrollHalfPageDown()
			return true
		case tcell.KeyCtrlU:
			app.pager.ScrollHalfPageUp()
			return true
		case tcell.KeyCtrlB:
			app.pager.ScrollPageUp()
			return true
		case tcell.KeyCtrlF:
			app.pager.ScrollPageDown()
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
				if app.pagerEnabled {
					app.pager.ScrollUp()
				} else {
					app.SelectPrevPod()
				}
				return true
			case 'j':
				if app.pagerEnabled {
					app.pager.ScrollDown()
				} else {
					app.SelectNextPod()
				}
				return true
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

	app.tabs.SelectAt(index)

	pod := app.pods[app.selectedPod]
	container := app.containers[app.selectedContainer]
	app.StartTailLog(pod.Namespace, pod.Name, container)
}

func (app *App) AddPod(pod *corev1.Pod) {
	app.pods = append(app.pods, pod)
	app.podsView.AddItem(pod.Name, tcell.StyleDefault)
}

func (app *App) ShowPager() {
	if app.pagerEnabled {
		return
	}
	app.AddWidget(app.line, 0)
	app.AddWidget(app.mainLayout, 0.9)
	app.pagerEnabled = true

	app.SelectPodAt(app.selectedPod)
}

func (app *App) HidePager() {
	app.RemoveWidget(app.line)
	app.RemoveWidget(app.mainLayout)
	app.pagerEnabled = false

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
		panic(err)
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
				}
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
		panic(err)
	}
}

func (app *App) Run(ctx context.Context) error {
	app.StartTailPods()
	return app.Application.Run()
}
