package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecurityAnalyzer struct{ c *Clients }
func NewSecurityAnalyzer(c *Clients) *SecurityAnalyzer { return &SecurityAnalyzer{c} }
func (a *SecurityAnalyzer) Name() string   { return "securityAnalyzer" }
func (a *SecurityAnalyzer) Optional() bool { return true }

func (a *SecurityAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	pods, err := a.c.K8s.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pod := range pods.Items {
		ns, name := pod.Namespace, pod.Name

		for _, c := range pod.Spec.Containers {
			if c.SecurityContext == nil {
				issues = append(issues, issue(a.Name(), "Pod", ns, name,
					"NoSecurityContext", Warning,
					fmt.Sprintf("Container '%s' sans SecurityContext ??? best practice non respectee", c.Name)))
				continue
			}

			// Root container
			if c.SecurityContext.RunAsNonRoot == nil || !*c.SecurityContext.RunAsNonRoot {
				if c.SecurityContext.RunAsUser == nil || *c.SecurityContext.RunAsUser == 0 {
					issues = append(issues, issue(a.Name(), "Pod", ns, name,
						"ContainerRunningAsRoot", Critical,
						fmt.Sprintf("Container '%s' tourne potentiellement en root ??? risque securite", c.Name)))
				}
			}

			// Privileged container
			if c.SecurityContext.Privileged != nil && *c.SecurityContext.Privileged {
				issues = append(issues, issue(a.Name(), "Pod", ns, name,
					"PrivilegedContainer", Critical,
					fmt.Sprintf("Container '%s' en mode privileged ??? acces total au node", c.Name)))
			}

			// AllowPrivilegeEscalation
			if c.SecurityContext.AllowPrivilegeEscalation == nil ||
				*c.SecurityContext.AllowPrivilegeEscalation {
				issues = append(issues, issue(a.Name(), "Pod", ns, name,
					"PrivilegeEscalationAllowed", Warning,
					fmt.Sprintf("Container '%s': AllowPrivilegeEscalation non desactive", c.Name)))
			}
		}

		// ServiceAccount par defaut
		if pod.Spec.ServiceAccountName == "" || pod.Spec.ServiceAccountName == "default" {
			issues = append(issues, issue(a.Name(), "Pod", ns, name,
				"DefaultServiceAccount", Info,
				"Pod utilise le ServiceAccount 'default' ??? creer un SA dedie avec RBAC minimal"))
		}
	}
	return issues, nil
}
