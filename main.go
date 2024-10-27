package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"afonya/policies"
)

// SetupKubeClient initializes the Kubernetes clientset.
func SetupKubeClient() (*kubernetes.Clientset, error) {
	// Load kubeconfig from home directory
	var kubeconfig *string
	// Retrieve home directory
	if home := homedir.HomeDir(); home != "" {
		// If there's home, then the kubeconfig is usually at ~/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// If there's no home, then the path must be defined
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	
	// Build the config to connect to the K8s API using the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	
	// Create Clientset to access K8s resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func main() {
	clientset, err := SetupKubeClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	fmt.Println("Connected to Kubernetes cluster")

	pm := policies.NewPolicyManager(clientset)

	// Example: List network policies in "default" ns
	policies, err := pm.ListPolicies("default")
	if err != nil {
		log.Fatalf("Failed to list network policies: %v", err)
	}
	fmt.Println("Network Policies in 'default' namespace:")
	for _, policy := range policies {
		fmt.Printf(" - %s\n", policy.Name)
	}

	// Example: Create a new network policy
	newPolicy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "allow-http",
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					Ports: []v1.NetworkPolicyPort{
						{
							//Protocol: &v1.ProtocolTCP,
							Port:	  &intstr.IntOrString{IntVal: 80},
						},
					},
				},
			},
		},
	}

	if err := pm.CreatePolicy("default", newPolicy); err != nil {
		log.Fatalf("Failed  to create network policy: %v", err)
	}

	// Example: Delete a network policy
	//if err := pm.DeletePolicy("default", "allow-http"); err != nil {
	//	log.Fatalf("Failed to delete network policy: %v", err)
	//}
}
	
