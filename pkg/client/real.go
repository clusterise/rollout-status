package client

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesImpl struct {
	clientset *kubernetes.Clientset
}

func (impl KubernetesImpl) ListAppsV1Deployments(namespace, selector string) (*appsv1.DeploymentList, error) {
	return impl.clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{LabelSelector: selector})
}

func (impl KubernetesImpl) ListAppsV1ReplicaSets(deployment *appsv1.Deployment) (*appsv1.ReplicaSetList, error) {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return impl.clientset.AppsV1().ReplicaSets(deployment.Namespace).List(options)
}

func (impl KubernetesImpl) ListV1Pods(replicasSet *appsv1.ReplicaSet) (*v1.PodList, error) {
	selector, err := metav1.LabelSelectorAsSelector(replicasSet.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return impl.clientset.CoreV1().Pods(replicasSet.Namespace).List(options)
}

func FromClientset(clientset *kubernetes.Clientset) *KubernetesImpl {
	return &KubernetesImpl{
		clientset: clientset,
	}
}
