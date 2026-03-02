package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PDBAnalyzer struct{ c *Clients }
func NewPDBAnalyzer(c *Clients) *PDBAnalyzer { return &PDBAnalyzer{c} }
func (a *PDBAnalyzer) Name() string   { return "pdbAnalyzer" }
func (a *PDBAnalyzer) Optional() bool { return true }

func (a *PDBAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	pdbs, err := a.c.K8s.PolicyV1().PodDisruptionBudgets("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pdb := range pdbs.Items {
		ns, name := pdb.Namespace, pdb.Name

		// PDB qui bloque les disruptions
		if pdb.Status.DisruptionsAllowed == 0 && pdb.Status.CurrentHealthy < pdb.Status.DesiredHealthy {
			issues = append(issues, issue(a.Name(), "PodDisruptionBudget", ns, name,
				"PDBBlockingDisruptions", Critical,
				fmt.Sprintf("PDB bloque les disruptions: %d/%d pods healthy",
					pdb.Status.CurrentHealthy, pdb.Status.DesiredHealthy)))
		}

		// PDB trop restrictif (maxUnavailable=0)
		if pdb.Spec.MaxUnavailable != nil && pdb.Spec.MaxUnavailable.IntValue() == 0 {
			issues = append(issues, issue(a.Name(), "PodDisruptionBudget", ns, name,
				"PDBTooRestrictive", Warning,
				"PDB avec maxUnavailable=0 ??? maintenance du node impossible"))
		}
	}
	return issues, nil
}
