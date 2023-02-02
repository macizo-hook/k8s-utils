/*
This application is a simple Go program that performs a smoke test against Kubernetes resources.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Create the client config
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println(err)
	}

	// Create a client object to interface with the cluster
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Check the total number of pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Check the readiness of deployments
	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// Oops, no deployments
	if deployments.Items == nil {
		log.Println("No deployments found")
		return
	}
	for _, deployment := range deployments.Items {
		// Check the deployments for readiness
		ready := deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
		fmt.Printf("Deployment %s is %s\n", deployment.Name, map[bool]string{true: "ready", false: "not ready"}[ready])
	}

}
