package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkPolicyAnalyzer struct{ c *Clients }
func NewNetworkPolicyAnalyzer(c *Clients) *NetworkPolicyAnalyzer { return &NetworkPolicyAnalyzer{c} }
func (a *NetworkPolicyAnalyzer) Name() string   { return "networkPolicyAnalyzer" }
func (a *NetworkPolicyAnalyzer) Optional() bool { return true }

func (a *NetworkPolicyAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	policies, err := a.c.K8s.NetworkingV1().NetworkPolicies("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	// Detecter les namespaces sans NetworkPolicy (risque securite)
	namespacesWithPolicy := map[string]bool{}
	for _, np := range policies.Items {
		namespacesWithPolicy[np.Namespace] = true
	}

	namespaces, err := a.c.K8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	ignored := map[string]bool{"kube-system": true, "kube-public": true, "kube-node-lease": true}
	for _, ns := range namespaces.Items {
		if ignored[ns.Name] { continue }
		if !namespacesWithPolicy[ns.Name] {
			issues = append(issues, issue(a.Name(), "Namespace", "", ns.Name,
				"NoNetworkPolicy", Info,
				fmt.Sprintf("Namespace '%s' sans NetworkPolicy ??? tout le trafic est autorise", ns.Name)))
		}
	}
	return issues, nil
}
