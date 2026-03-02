package analyzers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceAnalyzer struct{ c *Clients }
func NewServiceAnalyzer(c *Clients) *ServiceAnalyzer { return &ServiceAnalyzer{c} }
func (a *ServiceAnalyzer) Name() string   { return "serviceAnalyzer" }
func (a *ServiceAnalyzer) Optional() bool { return false }

func (a *ServiceAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	services, err := a.c.K8s.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, svc := range services.Items {
		ns, name := svc.Namespace, svc.Name
		if svc.Spec.Type == corev1.ServiceTypeExternalName { continue }

		// Service sans endpoints
		ep, err := a.c.K8s.CoreV1().Endpoints(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil { continue }

		hasReady := false
		for _, subset := range ep.Subsets {
			if len(subset.Addresses) > 0 { hasReady = true; break }
		}
		if !hasReady && svc.Spec.Type != corev1.ServiceTypeExternalName {
			issues = append(issues, issue(a.Name(), "Service", ns, name,
				"NoEndpoints", Warning,
				"Service sans endpoints actifs ??? les selectors correspondent-ils aux pods ?"))
		}

		// LoadBalancer sans IP externe depuis trop longtemps
		if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
			if len(svc.Status.LoadBalancer.Ingress) == 0 {
				issues = append(issues, issue(a.Name(), "Service", ns, name,
					"LoadBalancerPending", Warning,
					fmt.Sprintf("Service LoadBalancer %s sans IP externe ??? cloud provider operationnel ?", name)))
			}
		}
	}
	return issues, nil
}
