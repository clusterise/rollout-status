package status

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

// https://github.com/kubernetes/kubernetes/blob/dde6e8e7465468c32642659cb708a5cc922add64/staging/src/k8s.io/kubectl/pkg/util/deployment/deployment.go#L36
const RevisionAnnotation = "deployment.kubernetes.io/revision"

func DeploymentStatus(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) RolloutStatus {
	log.Printf("checking status for deployment %v", deployment.Name)

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == v1.ConditionFalse {
			err := MakeRolloutErorr("deployment %v is not progressing: %v", deployment.Name, condition.Message)
			return RolloutFatal(err)
		}
	}

	replicasSetList, err := getReplicaSetsByDeployment(clientset, deployment)
	if err != nil {
		return RolloutFatal(err)
	}

	lastRevision, ok := deployment.Annotations[RevisionAnnotation]
	if !ok {
		return RolloutFatal(fmt.Errorf("missing annotation %v on deployment %v", RevisionAnnotation, deployment.Name))
	}
	log.Printf("  last revision is %v", lastRevision)

	groupErr := ErrorGroup{}
	aggregatedStatus := RolloutStatus{
		Continue: true,
	}
	for _, replicaSet := range replicasSetList.Items {
		rsRevision, ok := replicaSet.Annotations[RevisionAnnotation]
		if !ok {
			groupErr.Add(fmt.Errorf("missing annotation %v on replicaset %v", RevisionAnnotation, replicaSet.Name))
			aggregatedStatus.Continue = false
			break
		}

		if rsRevision != lastRevision {
			// RS that is live and we no longer want
			// TODO make sure there are no more pods running in this RS
			// or older RS or RS that never became live and was overwritten
			continue
		}

		status := TestReplicaSetStatus(clientset, replicaSet)
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

func getReplicaSetsByDeployment(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) (*appsv1.ReplicaSetList, error) {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return clientset.AppsV1().ReplicaSets(deployment.Namespace).List(options)
}
