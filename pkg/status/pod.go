package status

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

func TestPodStatus(pod *v1.Pod) RolloutStatus {
	log.Printf("    checking status for pod %v", pod.Name)

	for _, containerStatus := range pod.Status.ContainerStatuses {
		// https://github.com/kubernetes/kubernetes/blob/4fda1207e347af92e649b59d60d48c7021ba0c54/pkg/kubelet/container/sync_result.go#L37
		// fail if the pod is containerruntimeerror (misconfigured env, missing image, etc)
		// fail if the pod is in crashloop backoff
		if containerStatus.State.Waiting != nil {
			switch containerStatus.State.Waiting.Reason {
			case "CrashLoopBackOff":
				// TODO this should retry but have a deadline, all restarts fall to CrashLoopBackOff
				fallthrough
			case "ErrImagePull":
				fallthrough
			case "ImagePullBackOff":
				fallthrough
			case "RunContainerError":
				err := MakeRolloutErorr("container %v is in %v: %v", containerStatus.Name, containerStatus.State.Waiting.Reason, containerStatus.State.Waiting.Message)
				return RolloutFatal(err)
			}
		}
	}

	if pod.Status.Phase == v1.PodPending {
		for _, condition := range pod.Status.Conditions {
			// fail if the pod is pending for X time
			if condition.Type == v1.PodScheduled {
				deadline := metav1.NewTime(time.Now().Add(time.Minute * -3)) // TODO configure
				err := MakeRolloutErorr("failed to scheduled pod: %v", condition.Message)

				if condition.LastTransitionTime.Before(&deadline) {
					return RolloutFatal(err)
				} else {
					return RolloutErrorProgressing(err)
				}
			}
		}
	}

	return RolloutOk()
}
