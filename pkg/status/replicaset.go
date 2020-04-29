package status

import (
	"dite.pro/rollout-status/pkg/client"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

func TestReplicaSetStatus(wrapper client.Kubernetes, replicaSet appsv1.ReplicaSet) RolloutStatus {
	for _, rsCondition := range replicaSet.Status.Conditions {
		if rsCondition.Type == appsv1.ReplicaSetReplicaFailure && rsCondition.Status == v1.ConditionTrue {
			err := MakeRolloutErorr("replicaset %q failed to create pods: %v", replicaSet.Name, rsCondition.Message)
			return RolloutFatal(err)
		}
	}

	podList, err := wrapper.ListV1Pods(&replicaSet)
	if err != nil {
		return RolloutFatal(err)
	}

	aggregatedStatus := RolloutStatus{
		Continue: true,
		Error:    nil,
	}
	for _, pod := range podList.Items {
		status := TestPodStatus(&pod)
		if status.Error != nil {
			aggregatedStatus.Error = status.Error
		}
		if !status.Continue {
			aggregatedStatus.Continue = false
			break
		}
	}

	return aggregatedStatus
}
