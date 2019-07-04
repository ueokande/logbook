package main

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/ui"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppConfig struct {
	Namespace  string
	KubeConfig string
}

type App struct {
	podsView *ui.ListView

	namespace   string
	kubeconfig  string
	pods        []corev1.Pod
	selectedPod int

	*views.Application
	views.BoxLayout
}

func NewApp(config *AppConfig) *App {
	podsView := ui.NewListView()

	app := &App{
		namespace:  config.Namespace,
		kubeconfig: config.KubeConfig,

		Application: new(views.Application),

		podsView: podsView,
	}

	app.SetOrientation(views.Vertical)
	app.AddWidget(app.podsView, 0)
	app.SetRootWidget(app)

	return app
}

func (app *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			app.Quit()
			return true
		}

		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				app.Quit()
				return true
			case 'k':
				app.SelectPrevPod()
				return true
			case 'j':
				app.SelectNextPod()
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
}

func (app *App) SelectPrevPod() {
	app.selectedPod -= 1
	if app.selectedPod < 0 {
		app.selectedPod = len(app.pods) - 1
	}
	app.podsView.SelectAt(app.selectedPod)
}

func (app *App) AddPod(pod corev1.Pod) {
	app.pods = append(app.pods, pod)
	app.podsView.AddItem(pod.Name, tcell.StyleDefault)
}

func (app *App) Run(ctx context.Context) error {
	clientset, err := NewClientset(app.kubeconfig)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(app.namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		app.AddPod(p)
	}
	app.podsView.SelectAt(0)

	return app.Application.Run()
}
