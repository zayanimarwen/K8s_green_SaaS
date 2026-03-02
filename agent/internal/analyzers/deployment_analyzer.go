package analyzers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentAnalyzer struct{ c *Clients }
func NewDeploymentAnalyzer(c *Clients) *DeploymentAnalyzer { return &DeploymentAnalyzer{c} }
func (a *DeploymentAnalyzer) Name() string   { return "deploymentAnalyzer" }
func (a *DeploymentAnalyzer) Optional() bool { return false }

func (a *DeploymentAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	deployments, err := a.c.K8s.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, d := range deployments.Items {
		ns, name := d.Namespace, d.Name
		desired := int32(1)
		if d.Spec.Replicas != nil { desired = *d.Spec.Replicas }

		// Aucune replica disponible
		if desired > 0 && d.Status.AvailableReplicas == 0 {
			issues = append(issues, issue(a.Name(), "Deployment", ns, name,
				"NoAvailableReplicas", Critical,
				fmt.Sprintf("Deployment '%s' sans replicas disponibles (%d desires)", name, desired)))
		}

		// Replicas partielles
		if d.Status.AvailableReplicas > 0 && d.Status.AvailableReplicas < desired {
			issues = append(issues, issue(a.Name(), "Deployment", ns, name,
				"PartialReplicas", Warning,
				fmt.Sprintf("Deployment %d/%d replicas disponibles", d.Status.AvailableReplicas, desired)))
		}

		// Rollout bloque
		for _, cond := range d.Status.Conditions {
			if cond.Type == "Progressing" && cond.Status == "False" {
				age := time.Since(cond.LastUpdateTime.Time)
				if age > 5*time.Minute {
					issues = append(issues, issue(a.Name(), "Deployment", ns, name,
						"RolloutStuck", Critical,
						fmt.Sprintf("Rollout bloque depuis %s: %s", age.Round(time.Second), cond.Message)))
				}
			}
		}
	}
	return issues, nil
}
