package main

import (
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)

// https://github.com/kubernetes/kubernetes/blob/dde6e8e7465468c32642659cb708a5cc922add64/staging/src/k8s.io/kubectl/pkg/util/deployment/deployment.go#L36
const RevisionAnnotation = "deployment.kubernetes.io/revision"

func main() {
	namespace := flag.String("namespace", "", "Namespace to watch")
	selector := flag.String("selector", "", "Label selector to watch, kubectl format such as release=foo,component=frontend")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deployments, err := clientset.AppsV1().Deployments(*namespace).List(metav1.ListOptions{LabelSelector: *selector})
	if err != nil {
		panic(err.Error())
	}

	//https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/cmd/rollout/rollout_status.go
	//https://github.com/kubernetes/kubernetes/blob/47daccb272c1a98c7b005dc1c19a88dbb643a3ee/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L59
	for _, deployment := range deployments.Items {
		// todo run all of those in goroutines concurrently
		err := deploymentStatus(clientset, &deployment)
		// todo check what kind of error is this
		// if deployment error, print different message
		if err != nil {
			panic(err.Error())
		}
	}

	log.Println("completed")
}

func getReplicaSetsByDeployment(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) (*appsv1.ReplicaSetList, error) {
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return clientset.AppsV1().ReplicaSets(deployment.Namespace).List(options)
}

func getPodsByReplicaSet(clientset *kubernetes.Clientset, replicasSet *appsv1.ReplicaSet) (*v1.PodList, error) {
	selector, err := metav1.LabelSelectorAsSelector(replicasSet.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	return clientset.CoreV1().Pods(replicasSet.Namespace).List(options)
}

func deploymentStatus(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) error {
	log.Printf("checking status for deployment %v", deployment.Name)

	replicasSetList, err := getReplicaSetsByDeployment(clientset, deployment)
	if err != nil {
		return err
	}

	lastRevision, ok := deployment.Annotations[RevisionAnnotation]
	if !ok {
		return fmt.Errorf("missing annotation %v on deployment %v", RevisionAnnotation, deployment.Name)
	}
	log.Printf("  last revision is %v", lastRevision)

	for _, replicasSet := range replicasSetList.Items {
		rsRevision, ok := replicasSet.Annotations[RevisionAnnotation]
		if !ok {
			return fmt.Errorf("missing annotation %v on replicaset %v", RevisionAnnotation, replicasSet.Name)
		}

		if rsRevision != lastRevision {
			// RS that is live and we no longer want
			// TODO make sure there are no more pods running in this RS
			// or older RS or RS that never became live and was overwritten
			continue
		}

		log.Printf("  checking status for replicaset %v", replicasSet.Name)

		podList, err := getPodsByReplicaSet(clientset, &replicasSet)
		if err != nil {
			return err
		}
		for _, pod := range podList.Items {
			err = testPodStatus(&pod)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testPodStatus(pod *v1.Pod) error {
	log.Printf("    checking status for pod %v", pod.Name)

	if pod.Status.Phase == v1.PodPending {
		for _, condition := range pod.Status.Conditions {
			// fail if the pod is pending for X time
			if condition.Type == v1.PodScheduled {
				deadline := metav1.NewTime(time.Now().Add(time.Minute * -3)) // TODO configure
				if condition.LastTransitionTime.Before(&deadline) {
					return fmt.Errorf("failed to scheduled pod: %v", condition.Message)
				}
			}
		}
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		// https://github.com/kubernetes/kubernetes/blob/4fda1207e347af92e649b59d60d48c7021ba0c54/pkg/kubelet/container/sync_result.go#L37
		// fail if the pod is containerruntimeerror (misconfigured env, missing image, etc)
		// fail if the pod is in crashloop backoff
		if containerStatus.State.Waiting != nil {
			switch containerStatus.State.Waiting.Reason {
			case "CrashLoopBackOff":
				fallthrough
			case "ImagePullBackOff":
				return fmt.Errorf("container %v is in %v", containerStatus.Name, containerStatus.State.Waiting.Reason)
			}
		}
	}

	// TODO fail if the pod is not running and state did not change for X time
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
