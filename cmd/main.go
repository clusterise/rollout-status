package main

import (
	"dite.pro/rollout-status/pkg/client"
	"dite.pro/rollout-status/pkg/status"
	"encoding/json"
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

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
	wrapper := client.FromClientset(clientset)

	rollout := status.TestRollout(wrapper, *namespace, *selector)
	out, err := json.MarshalIndent(rollout, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(out))
	//log.Println("Rollout successfully completed")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
