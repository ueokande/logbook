package main

import (
	"bufio"
	"context"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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

	podsView     *ui.ListView
	line         *ui.VerticalLine
	pager        *ui.Pager
	pagerEnabled bool

	namespace   string
	pods        []*corev1.Pod
	selectedPod int
	podworker   *Worker
	logworker   *Worker

	*views.Application
	views.BoxLayout
}

func NewApp(clientset *kubernetes.Clientset, config *AppConfig) *App {
	podsView := ui.NewListView()
	line := ui.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := ui.NewPager()

	app := &App{
		clientset: clientset,

		podsView: podsView,
		line:     line,
		pager:    pager,

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
	app.selectedPod += 1
	if app.selectedPod >= len(app.pods) {
		app.selectedPod = 0
	}
	app.podsView.SelectAt(app.selectedPod)

	if app.pagerEnabled {
		pod := app.pods[app.selectedPod]
		app.StartTailLog(pod)
	}
}

func (app *App) SelectPrevPod() {
	app.selectedPod -= 1
	if app.selectedPod < 0 {
		app.selectedPod = len(app.pods) - 1
	}
	app.podsView.SelectAt(app.selectedPod)

	if app.pagerEnabled {
		pod := app.pods[app.selectedPod]
		app.StartTailLog(pod)
	}
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
	app.AddWidget(app.pager, 0.9)
	app.pagerEnabled = true

	pod := app.pods[app.selectedPod]
	app.StartTailLog(pod)
}

func (app *App) HidePager() {
	app.RemoveWidget(app.line)
	app.RemoveWidget(app.pager)
	app.pagerEnabled = false

	app.StopTailLog()
}

func (app *App) StartTailLog(pod *corev1.Pod) {
	app.StopTailLog()

	app.pager.ClearText()
	app.logworker.Start(func(ctx context.Context) error {
		opts := &corev1.PodLogOptions{
			Container: pod.Spec.Containers[0].Name,
			Follow:    true,
		}
		req := app.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts)
		req.Context(ctx)
		r, err := req.Stream()
		if err != nil {
			return err
		}
		defer r.Close()

		s := bufio.NewScanner(r)
		for s.Scan() {
			app.pager.WriteText(s.Text() + "\n")
			app.Update()
		}
		return s.Err()
	})
}

func (app *App) StopTailLog() {
	err := app.logworker.Stop()
	if err != nil && err != context.Canceled {
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
			pod, ok := ev.Object.(*corev1.Pod)
			if !ok {
				continue
			}

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
			app.Update()
		}
		return nil
	})
}

func (app *App) StopTailPods() {
	err := app.podworker.Stop()
	if err != nil && err != context.Canceled {
		panic(err)
	}
}

func (app *App) Run(ctx context.Context) error {
	app.StartTailPods()
	return app.Application.Run()
}
