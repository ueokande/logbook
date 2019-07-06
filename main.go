package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var homedir string

func init() {
	if h := os.Getenv("HOME"); h != "" {
		homedir = h
	} else {
		homedir = os.Getenv("USERPROFILE") // windows
	}
}

type Params struct {
	namespace  string
	kubeconfig string
}

func main() {
	params := Params{
		namespace:  "default",
		kubeconfig: filepath.Join(homedir, ".kube", "config"),
	}

	cmd := &cobra.Command{}
	cmd.Short = "View logs on multiple pods and containers from Kubernetes"

	cmd.Flags().StringVarP(&params.namespace, "namespace", "n", params.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	cmd.Flags().StringVarP(&params.kubeconfig, "kubeconfig", "", params.kubeconfig, " Path to kubeconfig file to use")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := AppConfig{
			Namespace: params.namespace,
		}

		clientset, err := NewClientset(params.kubeconfig)
		if err != nil {
			return err
		}

		app := NewApp(clientset, &config)
		err = app.Run(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
