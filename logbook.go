package main

import (
	"context"
	"time"

	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/pkg/k8s"
	"github.com/ueokande/logbook/pkg/types"
	"github.com/ueokande/logbook/pkg/ui"
	corev1 "k8s.io/api/core/v1"
)

// AppConfig is a config for Logbook App
type AppConfig struct {
	Cluster   string
	Namespace string
}

// App is an application of logbook
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

// NewApp returns new App instance
func NewApp(client *k8s.Client, config *AppConfig) *App {
	w := ui.NewUI()
	w.SetContext(config.Cluster, config.Namespace)
	w.SetStatusMode(ui.ModeNormal)

	app := &App{
		client: client,
		ui:     w,

		namespace: config.Namespace,
		logworker: NewWorker(context.TODO()),
		podworker: NewWorker(context.TODO()),

		Application: new(views.Application),
	}

	w.WatchUIEvents(app)
	app.SetRootWidget(w)

	return app
}

// OnContainerSelected handles events on container selected by UI
func (app *App) OnContainerSelected(name string, index int) {
	pod := app.currentPod
	app.ui.ClearPager()
	app.StartTailLog(pod.Namespace, pod.Name, name)
}

// OnPodSelected handles events on pod selected by UI
func (app *App) OnPodSelected(name string, index int) {
	app.currentPod = app.pods[index]
	pod := app.currentPod
	app.ui.ClearContainers()
	for _, c := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		app.ui.AddContainer(c.Name)
	}
	app.ui.SelectContainerAt(0)
}

// OnQuit handles events on quit is required by UI
func (app *App) OnQuit() {
	app.Quit()
}

// StartTailLog starts tailing logs for container of pod in namespace
func (app *App) StartTailLog(namespace, pod, container string) {
	app.StopTailLog()

	app.logworker.Start(func(ctx context.Context) error {
		logs, err := app.client.WatchLogs(ctx, namespace, pod, container)
		if err != nil {
			return err
		}

		b := NewTimeBuffer(100*time.Millisecond, 100)
		go func() {
			defer b.Close()

			for log := range logs {
				b.Write(log)
			}
		}()

		// make channel to guarantee line order of logs
		for lines := range b.Flushed() {
			lines := lines
			app.PostFunc(func() {
				app.ui.AddPagerTexts(lines)
			})
		}
		return nil
	})
}

// StopTailLog stops tailing logs
func (app *App) StopTailLog() {
	app.logworker.Stop()
	// TODO handle err
}

// StartTailPods tarts tailing pods on Kubernetes
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
					app.ui.AddPod(pod.Name, types.GetPodStatus(pod))
					if len(app.pods) == 1 {
						app.ui.SelectPodAt(0)
					}
				case k8s.PodModified:
					app.ui.SetPodStatus(pod.Name, types.GetPodStatus(pod))
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

// StopTailPods stops tailing pods
func (app *App) StopTailPods() {
	app.podworker.Stop()
	// TODO handle err
}

// Run starts logbook application
func (app *App) Run(ctx context.Context) error {
	app.StartTailPods()
	return app.Application.Run()
}
