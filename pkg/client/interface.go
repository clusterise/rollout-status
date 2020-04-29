package client

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type Kubernetes interface {
	ListAppsV1Deployments(namespace, selector string) (*appsv1.DeploymentList, error)
	ListAppsV1ReplicaSets(deployment *appsv1.Deployment) (*appsv1.ReplicaSetList, error)
	ListV1Pods(replicasSet *appsv1.ReplicaSet) (*v1.PodList, error)
}
