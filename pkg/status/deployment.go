package status

import (
	"dite.pro/rollout-status/pkg/client"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"log"
)

// https://github.com/kubernetes/kubernetes/blob/dde6e8e7465468c32642659cb708a5cc922add64/staging/src/k8s.io/kubectl/pkg/util/deployment/deployment.go#L36
const RevisionAnnotation = "deployment.kubernetes.io/revision"

func DeploymentStatus(wrapper client.Kubernetes, deployment *appsv1.Deployment) RolloutStatus {
	log.Printf("checking status for deployment %v", deployment.Name)

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == v1.ConditionFalse {
			err := MakeRolloutErorr("deployment %q is not progressing: %v", deployment.Name, condition.Message)
			return RolloutFatal(err)
		}
	}

	replicasSetList, err := wrapper.ListAppsV1ReplicaSets(deployment)
	if err != nil {
		return RolloutFatal(err)
	}

	lastRevision, ok := deployment.Annotations[RevisionAnnotation]
	if !ok {
		return RolloutFatal(fmt.Errorf("missing annotation %q on deployment %q", RevisionAnnotation, deployment.Name))
	}
	log.Printf("  last revision is %v", lastRevision)

	aggregatedStatus := RolloutStatus{
		Continue: true,
		Error: nil,
	}
	for _, replicaSet := range replicasSetList.Items {
		rsRevision, ok := replicaSet.Annotations[RevisionAnnotation]
		if !ok {
			aggregatedStatus.Error = fmt.Errorf("missing annotation %q on replicaset %q", RevisionAnnotation, replicaSet.Name)
			aggregatedStatus.Continue = false
			break
		}

		if rsRevision != lastRevision {
			// RS that is live and we no longer want
			// TODO make sure there are no more pods running in this RS
			// or older RS or RS that never became live and was overwritten
			continue
		}

		status := TestReplicaSetStatus(wrapper, replicaSet)
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
