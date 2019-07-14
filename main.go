package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
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

		config, clientset, err := loadConfig(&params)
		if err != nil {
			return err
		}

		app := NewApp(clientset, config)
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

func loadConfig(params *Params) (*AppConfig, *kubernetes.Clientset, error) {
	kubeConfig := GetKubeConfig(params.kubeconfig)
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get raw config")
	}
	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get client config")
	}

	context := rawConfig.Contexts[rawConfig.CurrentContext]
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
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create clientset")
	}

	return config, clientset, nil
}
