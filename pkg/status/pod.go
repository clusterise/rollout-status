package status

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func TestPodStatus(pod *v1.Pod) RolloutStatus {
	aggregatedStatus := RolloutStatus{
		Continue: true,
		Error:    nil,
	}
	for _, initStatus := range pod.Status.InitContainerStatuses {
		status := TestContainerStatus(&initStatus)
		if status.Error != nil {
			if !status.Continue {
				return status
			}
			aggregatedStatus.Error = status.Error
		}
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		status := TestContainerStatus(&containerStatus)
		if status.Error != nil {
			if !status.Continue {
				return status
			}
			aggregatedStatus.Error = status.Error
		}
	}
	if aggregatedStatus.Error != nil {
		return aggregatedStatus
	}

	if pod.Status.Phase == v1.PodPending {
		for _, condition := range pod.Status.Conditions {
			// fail if the pod is pending for X time
			if condition.Type == v1.PodScheduled {
				deadline := metav1.NewTime(time.Now().Add(time.Minute * -3)) // TODO configure
				err := MakeRolloutErorr(FailureScheduling, "Failed to schedule pod: %v", condition.Message)

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
