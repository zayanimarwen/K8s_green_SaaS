package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatefulSetAnalyzer struct{ c *Clients }
func NewStatefulSetAnalyzer(c *Clients) *StatefulSetAnalyzer { return &StatefulSetAnalyzer{c} }
func (a *StatefulSetAnalyzer) Name() string   { return "statefulSetAnalyzer" }
func (a *StatefulSetAnalyzer) Optional() bool { return false }

func (a *StatefulSetAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	sets, err := a.c.K8s.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, sts := range sets.Items {
		ns, name := sts.Namespace, sts.Name
		desired := int32(1)
		if sts.Spec.Replicas != nil { desired = *sts.Spec.Replicas }

		// Replicas insuffisantes
		if sts.Status.ReadyReplicas < desired {
			issues = append(issues, issue(a.Name(), "StatefulSet", ns, name,
				"InsufficientReplicas", Critical,
				fmt.Sprintf("StatefulSet %d/%d replicas ready", sts.Status.ReadyReplicas, desired)))
		}

		// Mise a jour bloquee
		if sts.Status.UpdatedReplicas < desired && sts.Status.CurrentReplicas == desired {
			issues = append(issues, issue(a.Name(), "StatefulSet", ns, name,
				"UpdateStuck", Warning,
				"Mise a jour du StatefulSet bloquee ??? verifier les rolling update strategy"))
		}

		// Service headless manquant
		if sts.Spec.ServiceName != "" {
			_, err := a.c.K8s.CoreV1().Services(ns).Get(ctx, sts.Spec.ServiceName, metav1.GetOptions{})
			if err != nil {
				issues = append(issues, issue(a.Name(), "StatefulSet", ns, name,
					"HeadlessServiceMissing", Critical,
					fmt.Sprintf("Service headless '%s' manquant pour StatefulSet", sts.Spec.ServiceName)))
			}
		}
	}
	return issues, nil
}
