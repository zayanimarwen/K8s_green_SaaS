package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RSAnalyzer struct{ c *Clients }
func NewRSAnalyzer(c *Clients) *RSAnalyzer { return &RSAnalyzer{c} }
func (a *RSAnalyzer) Name() string   { return "rsAnalyzer" }
func (a *RSAnalyzer) Optional() bool { return false }

func (a *RSAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	rsList, err := a.c.K8s.AppsV1().ReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, rs := range rsList.Items {
		ns, name := rs.Namespace, rs.Name

		// RS orphelin (pas de owner Deployment)
		if len(rs.OwnerReferences) == 0 && *rs.Spec.Replicas > 0 {
			issues = append(issues, issue(a.Name(), "ReplicaSet", ns, name,
				"OrphanReplicaSet", Warning,
				"ReplicaSet sans Deployment proprietaire"))
		}

		// Replicas insuffisantes
		desired := int32(0)
		if rs.Spec.Replicas != nil { desired = *rs.Spec.Replicas }
		if desired > 0 && rs.Status.ReadyReplicas < desired {
			issues = append(issues, issue(a.Name(), "ReplicaSet", ns, name,
				"InsufficientReplicas", Warning,
				fmt.Sprintf("ReplicaSet: %d/%d replicas ready", rs.Status.ReadyReplicas, desired)))
		}
	}
	return issues, nil
}
