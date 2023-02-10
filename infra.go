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

	// Check the readiness of services
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if services.Items == nil {
		log.Println("No services found")
		return
	}
	for _, service := range services.Items {

		// Check if the service is ready
		ready := service.Spec.ClusterIP != ""
		fmt.Printf("Service %s is %s\n", service.Name, map[bool]string{true: "ready", false: "not ready"}[ready])

		// Check the readiness of endpoints
		endpoints, err := clientset.CoreV1().Endpoints(service.Namespace).Get(context.TODO(), service.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get endpoints for service %s: %v", service.Name, err)
			continue
		}
		ready = len(endpoints.Subsets) > 0
		fmt.Printf("Endpoints of service %s are %s\n", service.Name, map[bool]string{true: "ready", false: "not ready"}[ready])
	}
}
