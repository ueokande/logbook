package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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
	if len(pods.Items) == 0 {
		return nil
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	title := &views.TextBar{}
	title.SetCenter(fmt.Sprintf("Got %d pods", len(pods.Items)), tcell.StyleDefault)

	app.SetOrientation(views.Vertical)
	app.AddWidget(title, 0)

	app.SetRootWidget(app)
	if e := app.Application.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	return nil

}
