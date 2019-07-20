package k8s

import (
	"github.com/ueokande/logbook/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

func GetPodStatus(pod *corev1.Pod) types.PodStatus {
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		return types.PodSucceeded
	case corev1.PodPending:
		return types.PodPending
	case corev1.PodFailed:
		return types.PodFailed
	case corev1.PodUnknown:
		return types.PodUnknown
	}

	for _, container := range pod.Status.InitContainerStatuses {
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			return types.PodFailed
		default:
			return types.PodInitializing
		}
		break
	}

	hasCompleted := false
	hasRunning := false
	for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
		container := pod.Status.ContainerStatuses[i]
		if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
			return types.PodInitializing
		} else if container.State.Terminated != nil && container.State.Terminated.Reason == "Completed" {
			hasCompleted = true
		} else if container.State.Terminated != nil {
			return types.PodFailed
		} else if container.Ready && container.State.Running != nil {
			hasRunning = true
		}
	}
	if hasCompleted && hasRunning {
		return types.PodRunning
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == "NodeLost" {
		return types.PodUnknown
	} else if pod.DeletionTimestamp != nil {
		return types.PodTerminating
	}
	return types.PodRunning
}
