package policies

import (
	"context"
	"fmt"
	"k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PolicyManager manages network policies within the cluster.
type PolicyManager struct {
	clientset *kubernetes.Clientset
}

// NewPolicyManager creates a new PolicyManager.
func NewPolicyManager(clientset *kubernetes.Clientset) *PolicyManager {
	return &PolicyManager{clientset: clientset}
}

// CreatePolicy creates a new network policy in the specified namespace.
func (pm *PolicyManager) CreatePolicy(namespace string, policy *v1.NetworkPolicy) error {
	_, err := pm.clientset.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create network policy: %v", err)
	}
	fmt.Println("Network policy created successfully")
	return nil
}

// DeletePolicy deletes a network policy by name in a given namespace.
func (pm *PolicyManager) DeletePolicy(namespace, policyName string) error {
	err := pm.clientset.NetworkingV1().NetworkPolicies(namespace).Delete(context.TODO(), policyName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete network policy: %v", err)
	}
	fmt.Println("Network policy deleted successfully")
	return nil
}

// ListPolicies lists all network policies in the specified namespace.
func (pm *PolicyManager) ListPolicies(namespace string) ([]v1.NetworkPolicy, error) {
	policies, err := pm.clientset.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list network policies: %v", err)
	}
	return policies.Items, nil
}
