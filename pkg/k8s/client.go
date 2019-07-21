package k8s

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient(kubeconfig string) (*Client, error) {
	kubeConfig := getKubeConfig(kubeconfig)
	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func LoadCurrentContext(kubeconfig string) (*api.Context, error) {
	kubeConfig := getKubeConfig(kubeconfig)
	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get raw config")
	}
	return rawConfig.Contexts[rawConfig.CurrentContext], nil
}

func getKubeConfig(kubeconfig string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(kubeconfig) > 0 {
		rules.Precedence = []string{kubeconfig}
	}
	overrides := &clientcmd.ConfigOverrides{}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
