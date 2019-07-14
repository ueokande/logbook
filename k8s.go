package main

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type PodStatus string

const (
	PodRunning      PodStatus = "Running"
	PodSucceeded              = "Succeeded"
	PodPending                = "Pending"
	PodTerminating            = "Terminating"
	PodInitializing           = "Initializing"
	PodFailed                 = "Failed"
	PodUnknown                = "Unknown"
)

func GetKubeConfig(kubeconfig string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(kubeconfig) > 0 {
		rules.Precedence = []string{kubeconfig}
	}
	overrides := &clientcmd.ConfigOverrides{}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
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
