package status

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func TestReplicaSetStatus(clientset *kubernetes.Clientset, replicaSet appsv1.ReplicaSet) RolloutStatus {
	log.Printf("  checking status for replicaset %v", replicaSet.Name)

	for _, rsCondition := range replicaSet.Status.Conditions {
		if rsCondition.Type == appsv1.ReplicaSetReplicaFailure && rsCondition.Status == v1.ConditionTrue {
			err := MakeRolloutErorr("replicaset %v failed to create pods: %v", replicaSet.Name, rsCondition.Message)
			return RolloutFatal(err)
		}
	}

	podList, err := getPodsByReplicaSet(clientset, &replicaSet)
	if err != nil {
		return RolloutFatal(err)
	}

	groupErr := ErrorGroup{}
	aggregatedStatus := RolloutStatus{
		Continue: true,
	}
	for _, pod := range podList.Items {
		status := TestPodStatus(&pod)
		if status.Error != nil {
			groupErr.Add(err)
		}
		if !status.Continue {
			aggregatedStatus.Continue = false
			break
		}
	}
	if !groupErr.Empty() {
		aggregatedStatus.Error = groupErr
	}

	return aggregatedStatus
}

func getPodsByReplicaSet(clientset *kubernetes.Clientset, replicasSet *appsv1.ReplicaSet) (*v1.PodList, error) {
	selector, err := metav1.LabelSelectorAsSelector(replicasSet.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return clientset.CoreV1().Pods(replicasSet.Namespace).List(options)
}
