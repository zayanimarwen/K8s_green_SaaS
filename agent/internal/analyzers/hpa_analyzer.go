package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HPAAnalyzer struct{ c *Clients }
func NewHPAAnalyzer(c *Clients) *HPAAnalyzer { return &HPAAnalyzer{c} }
func (a *HPAAnalyzer) Name() string   { return "hpaAnalyzer" }
func (a *HPAAnalyzer) Optional() bool { return true }

func (a *HPAAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	hpas, err := a.c.K8s.AutoscalingV2().HorizontalPodAutoscalers("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, hpa := range hpas.Items {
		ns, name := hpa.Namespace, hpa.Name

		// HPA au max depuis trop longtemps
		if hpa.Status.CurrentReplicas >= hpa.Spec.MaxReplicas {
			issues = append(issues, issue(a.Name(), "HorizontalPodAutoscaler", ns, name,
				"HPAAtMaxReplicas", Warning,
				fmt.Sprintf("HPA au maximum (%d/%d replicas) ??? envisager d augmenter maxReplicas",
					hpa.Status.CurrentReplicas, hpa.Spec.MaxReplicas)))
		}

		// HPA sans metriques
		scalingActive := false
		for _, cond := range hpa.Status.Conditions {
			if string(cond.Type) == "ScalingActive" && cond.Status == "True" {
				scalingActive = true
			}
		}
		if !scalingActive {
			issues = append(issues, issue(a.Name(), "HorizontalPodAutoscaler", ns, name,
				"HPAScalingInactive", Critical,
				"HPA inactif ??? metrics-server disponible ? Metriques accessibles ?"))
		}

		// Ecart trop grand min/max
		minReplicas := int32(1)
		if hpa.Spec.MinReplicas != nil { minReplicas = *hpa.Spec.MinReplicas }
		if hpa.Spec.MaxReplicas > minReplicas*10 {
			issues = append(issues, issue(a.Name(), "HorizontalPodAutoscaler", ns, name,
				"HPAWideRange", Info,
				fmt.Sprintf("HPA avec plage large: min=%d max=%d ??? scaling agressif possible", minReplicas, hpa.Spec.MaxReplicas)))
		}
	}
	return issues, nil
}
