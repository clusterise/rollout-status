package status

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

func TestPodStatus(pod *v1.Pod) RolloutStatus {
	log.Printf("    checking status for pod %q", pod.Name)

	for _, initStatus := range pod.Status.InitContainerStatuses {
		if initStatus.State.Waiting!= nil {
			reason := initStatus.State.Waiting.Reason
			switch reason {
				case "CrashLoopBackOff":
					// TODO this should retry but have a deadline, all restarts fall to CrashLoopBackOff
					err := MakeRolloutErorr("init container %q is in %q", initStatus.Name, reason)
					return RolloutFatal(err)
			}
		}
		if initStatus.State.Terminated != nil {
			reason := initStatus.State.Terminated.Reason
			switch reason {
				case "Error":
					// TODO this should retry but have a deadline, all restarts fall to CrashLoopBackOff
					err := MakeRolloutErorr("init container %q is in %q", initStatus.Name, reason)
					return RolloutFatal(err)
			}
		}
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		// https://github.com/kubernetes/kubernetes/blob/4fda1207e347af92e649b59d60d48c7021ba0c54/pkg/kubelet/container/sync_result.go#L37
		// fail if the pod is containerruntimeerror (misconfigured env, missing image, etc)
		// fail if the pod is in crashloop backoff
		if containerStatus.State.Waiting != nil {
			reason := containerStatus.State.Waiting.Reason
			switch reason {
			case "ContainerCreating":
				fallthrough
			case "PodInitializing":
				err := MakeRolloutErorr("container %q is in %q", containerStatus.Name, reason)
				return RolloutErrorProgressing(err)

			case "CrashLoopBackOff":
				// TODO this should retry but have a deadline, all restarts fall to CrashLoopBackOff
				fallthrough
			case "RunContainerError":
				// TODO this should retry but have a deadline, all restarts fall to CrashLoopBackOff
				fallthrough

			case "ErrImagePull":
				fallthrough
			case "ImagePullBackOff":
				err := MakeRolloutErorr("container %q is in %q: %v", containerStatus.Name, reason, containerStatus.State.Waiting.Message)
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
