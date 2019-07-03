package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/ueokande/logbook/ui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppConfig struct {
	Namespace  string
	KubeConfig string
}

type App struct {
	views.BoxLayout
	*views.Application

	namespace  string
	kubeconfig string
}

func NewApp(config *AppConfig) *App {
	return &App{
		namespace:  config.Namespace,
		kubeconfig: config.KubeConfig,

		Application: new(views.Application),
	}
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
			}
		}
	}
	return app.BoxLayout.HandleEvent(ev)
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

	l := ui.NewListView()
	for _, p := range pods.Items {
		item := ui.ListItem{
			Text: p.Name,
		}
		l.AddItem(item)
	}

	app.SetOrientation(views.Vertical)
	app.AddWidget(l, 0)

	app.SetRootWidget(app)
	if e := app.Application.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return nil

}
