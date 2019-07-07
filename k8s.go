package main

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type PodStatus string
type ContainerStatus string

const (
	PodRunning                = "Running"
	PodSucceeded              = "Succeeded"
	PodPending      PodStatus = "Pending"
	PodTerminating            = "Terminating"
	PodInitializing           = "Initializing"
	PodFailed                 = "Failed"
	PodUnknown                = "Unknown"
)

func NewClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{}
	configLoadingRules.ExplicitPath = kubeconfig
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoadingRules,
		&clientcmd.ConfigOverrides{},
	)
	cconfig, err := config.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client config")
	}

	clientset, err := kubernetes.NewForConfig(cconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}

	return clientset, nil
}

func GetPodStatus(pod *corev1.Pod) PodStatus {
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		return PodSucceeded
	case corev1.PodPending:
		return PodPending
	case corev1.PodFailed:
		return PodFailed
	case corev1.PodUnknown:
		return PodUnknown
	}

	for _, container := range pod.Status.InitContainerStatuses {
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			return PodFailed
		default:
			return PodInitializing
		}
		break
	}

	hasCompleted := false
	hasRunning := false
	for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
		container := pod.Status.ContainerStatuses[i]
		if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
			return PodInitializing
		} else if container.State.Terminated != nil && container.State.Terminated.Reason == "Completed" {
			hasCompleted = true
		} else if container.State.Terminated != nil {
			return PodFailed
		} else if container.Ready && container.State.Running != nil {
			hasRunning = true
		}
	}
	if hasCompleted && hasRunning {
		return PodRunning
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		return PodUnknown
	} else if pod.DeletionTimestamp != nil {
		return PodTerminating
	}
	return PodRunning
}
