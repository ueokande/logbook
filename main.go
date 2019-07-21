package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ueokande/logbook/pkg/k8s"
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
	params := Params{}

	cmd := &cobra.Command{}
	cmd.Short = "View logs on multiple pods and containers from Kubernetes"

	cmd.Flags().StringVarP(&params.namespace, "namespace", "n", params.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	cmd.Flags().StringVarP(&params.kubeconfig, "kubeconfig", "", params.kubeconfig, " Path to kubeconfig file to use")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		context, err := k8s.LoadCurrentContext(params.kubeconfig)
		if err != nil {
			return err
		}

		client, err := k8s.NewClient(params.kubeconfig)
		if err != nil {
			return err
		}

		config := &AppConfig{
			Cluster:   context.Cluster,
			Namespace: "default",
		}
		if len(context.Namespace) > 0 {
			config.Namespace = context.Namespace
		}
		if len(params.namespace) > 0 {
			config.Namespace = context.Namespace
		}

		app := NewApp(client, config)
		err = app.Run(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
