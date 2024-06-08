package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		return
	}

	deployments, err := clientset.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}

	for _, deployment := range deployments.Items {
		if strings.Contains(deployment.Name, "database") {
			fmt.Printf("Restarting deployment: %s\n", deployment.Name)
			deployment.Spec.Template.Annotations = map[string]string{
				"kubectl.kubernetes.io/restartedAt": time.Now().Format(time.RFC3339),
			}
			if _, err := clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{}); err != nil {
				fmt.Printf("Error restarting deployment %s: %v\n", deployment.Name, err)
			}
		}
	}
}
