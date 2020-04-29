package status

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestRollout(clientset *kubernetes.Clientset, namespace, selector string) RolloutStatus {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return RolloutFatal(err)
	}

	if len(deployments.Items) == 0 {
		err = MakeRolloutErorr("Selector %q did not match any deployments in namespace %q", selector, namespace)
		return RolloutFatal(err)
	}

	groupErr := ErrorGroup{}
	aggregatedStatus := RolloutStatus{
		Continue: true,
	}

	//https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/cmd/rollout/rollout_status.go
	//https://github.com/kubernetes/kubernetes/blob/47daccb272c1a98c7b005dc1c19a88dbb643a3ee/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L59
	for _, deployment := range deployments.Items {
		// todo run all of those in goroutines concurrently
		status := DeploymentStatus(clientset, &deployment)
		if status.Error != nil {
			groupErr.Add(status.Error)
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
