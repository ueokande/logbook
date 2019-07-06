package main

import (
	"bufio"
	"context"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/ui"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
	pods        []corev1.Pod
	selectedPod int

	*views.Application
	views.BoxLayout
}

func NewApp(clientset *kubernetes.Clientset, config *AppConfig) *App {
	podsView := ui.NewListView()
	line := ui.NewVerticalLine(tcell.RuneVLine, tcell.StyleDefault)
	pager := ui.NewPager()

	app := &App{
		clientset: clientset,

		namespace: config.Namespace,

		Application: new(views.Application),

		podsView: podsView,
		line:     line,
		pager:    pager,
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
					// TODO scroll up pager
				} else {
					app.SelectPrevPod()
				}
				return true
			case 'j':
				if app.pagerEnabled {
					// TODO scroll down pager
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

func (app *App) AddPod(pod corev1.Pod) {
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

func (app *App) StartTailLog(pod corev1.Pod) {
	app.pager.ClearText()
	go func() {
		opts := &corev1.PodLogOptions{
			Container: pod.Spec.Containers[0].Name,
			Follow:    true,
		}
		req := app.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, opts)
		r, err := req.Stream()
		if err != nil {
			return
		}
		defer r.Close()

		s := bufio.NewScanner(r)
		for s.Scan() {
			app.pager.WriteText(s.Text() + "\n")
			app.Refresh()
		}
	}()
}

func (app *App) StopTailLog() {
}

func (app *App) Run(ctx context.Context) error {
	pods, err := app.clientset.CoreV1().Pods(app.namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		app.AddPod(p)
	}
	app.podsView.SelectAt(0)

	return app.Application.Run()
}
