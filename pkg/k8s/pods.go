package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// PodEventType represents an event type of the pod
type PodEventType int

// The event type of the pods
const (
	PodAdded    PodEventType = iota // The pod is added
	PodModified                     // The pod is updated
	PodDeleted                      // The pod is deleted
)

// PodEvent represents an event of the pods in Kubernetes API
type PodEvent struct {
	Type PodEventType
	Pod  *corev1.Pod
}

// WatchPods watches pods from Kubernetes API in namespace.  It returns a
// channel to subscribe pods.
func (c *Client) WatchPods(ctx context.Context, namespace string) (<-chan *PodEvent, error) {
	r, err := c.clientset.CoreV1().Pods(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	ch := make(chan *PodEvent)
	go func() {
		<-ctx.Done()
		r.Stop()
	}()
	go func() {
		for ev := range r.ResultChan() {
			pod, ok := ev.Object.(*corev1.Pod)
			if !ok {
				continue
			}
			var t PodEventType
			switch ev.Type {
			case watch.Added:
				t = PodAdded
			case watch.Modified:
				t = PodModified
			case watch.Deleted:
				t = PodDeleted
			default:
				continue
			}

			ch <- &PodEvent{
				Type: t,
				Pod:  pod,
			}
		}
		close(ch)
	}()
	return ch, nil
}
