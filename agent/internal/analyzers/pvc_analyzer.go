package analyzers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PVCAnalyzer struct{ c *Clients }
func NewPVCAnalyzer(c *Clients) *PVCAnalyzer { return &PVCAnalyzer{c} }
func (a *PVCAnalyzer) Name() string   { return "pvcAnalyzer" }
func (a *PVCAnalyzer) Optional() bool { return false }

func (a *PVCAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	pvcs, err := a.c.K8s.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pvc := range pvcs.Items {
		ns, name := pvc.Namespace, pvc.Name

		// PVC Pending
		if pvc.Status.Phase == corev1.ClaimPending {
			age := time.Since(pvc.CreationTimestamp.Time)
			if age > 2*time.Minute {
				issues = append(issues, issue(a.Name(), "PersistentVolumeClaim", ns, name,
					"PVCPending", Critical,
					fmt.Sprintf("PVC en attente de binding depuis %s ??? verifier StorageClass", age.Round(time.Second))))
			}
		}

		// PVC Lost
		if pvc.Status.Phase == corev1.ClaimLost {
			issues = append(issues, issue(a.Name(), "PersistentVolumeClaim", ns, name,
				"PVCLost", Critical,
				"PVC perdu ??? le PersistentVolume a ete supprime ou recycle"))
		}

		// PVC sans StorageClass
		if pvc.Spec.StorageClassName == nil || *pvc.Spec.StorageClassName == "" {
			issues = append(issues, issue(a.Name(), "PersistentVolumeClaim", ns, name,
				"NoStorageClass", Warning,
				"PVC sans StorageClass explicite ??? utilisera la classe par defaut"))
		}
	}
	return issues, nil
}
