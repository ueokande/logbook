package types

import (
	corev1 "k8s.io/api/core/v1"
)

// PodStatus represents pod's status
type PodStatus string

// PodStatus is a pod status of the pod.  They determined the pod's phase,
// status and its containers.
const (
	PodRunning      PodStatus = "Running"      // The pod is running successfully,
	PodSucceeded              = "Succeeded"    // The pod created by job is completed successfully,
	PodPending                = "Pending"      // The Pod has been accepted, but some of the container images has not been created.
	PodTerminating            = "Terminating"  // The pod is terminating
	PodInitializing           = "Initializing" // The init containers are running, or some of the container is initializing.
	PodFailed                 = "Failed"       // The pod fails on the init ocntaienrs, or one or more of the container exits without exit code 0.
	PodUnknown                = "Unknown"      // Unknown phase or status
)

// GetPodStatus returns the status of the pod as PodStatus
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
