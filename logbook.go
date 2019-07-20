package main

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/pkg/k8s"
	"github.com/ueokande/logbook/pkg/ui"
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
	ui     *ui.UI

	namespace  string
	pods       []*corev1.Pod
	currentPod *corev1.Pod
	podworker  *Worker
	logworker  *Worker

	*views.Application
}

func NewApp(client *k8s.Client, config *AppConfig) *App {
	w := ui.NewUI()
	w.SetContext(config.Cluster, config.Namespace)
	w.SetStatusMode(ui.ModeNormal)

	app := &App{
		client: client,
		ui:     w,

		namespace: config.Namespace,
		logworker: NewWorker(context.Background()),
		podworker: NewWorker(context.Background()),

		Application: new(views.Application),
	}

	w.WatchUIEvents(app)
	app.SetRootWidget(w)

	return app
}

func (app *App) OnContainerSelected(name string, index int) {
	pod := app.currentPod
	app.ui.ClearPager()
	app.StartTailLog(pod.Namespace, pod.Name, name)
}

func (app *App) OnPodSelected(name string, index int) {
	app.currentPod = app.pods[index]
	pod := app.currentPod
	app.ui.ClearContainers()
	for _, c := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		app.ui.AddContainer(c.Name)
	}
	app.ui.SelectContainerAt(0)
}

func (app *App) OnQuit() {
	app.Quit()
}

func (app *App) StartTailLog(namespace, pod, container string) {
	app.StopTailLog()

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
					app.ui.AddPagerText(line)
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
					app.ui.AddPod(pod.Name, k8s.GetPodStatus(pod))
					if len(app.pods) == 1 {
						app.ui.SelectPodAt(0)
					}
				case k8s.PodModified:
					app.ui.SetPodStatus(pod.Name, k8s.GetPodStatus(pod))
				case k8s.PodDeleted:
					for i, p := range app.pods {
						if p.Name == pod.Name {
							app.pods = append(app.pods[:i], app.pods[i+1:]...)
							break
						}
					}
					app.ui.DeletePod(pod.Name)
				}
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
