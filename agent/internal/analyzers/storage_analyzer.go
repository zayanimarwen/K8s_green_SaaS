package analyzers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StorageAnalyzer struct{ c *Clients }
func NewStorageAnalyzer(c *Clients) *StorageAnalyzer { return &StorageAnalyzer{c} }
func (a *StorageAnalyzer) Name() string   { return "storageAnalyzer" }
func (a *StorageAnalyzer) Optional() bool { return true }

func (a *StorageAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue

	// PersistentVolumes Failed
	pvs, err := a.c.K8s.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, pv := range pvs.Items {
		switch pv.Status.Phase {
		case "Failed":
			issues = append(issues, issue(a.Name(), "PersistentVolume", "", pv.Name,
				"PVFailed", Critical,
				fmt.Sprintf("PersistentVolume en echec: %s", pv.Status.Message)))
		case "Released":
			// PV released mais retain policy = gaspillage
			if pv.Spec.PersistentVolumeReclaimPolicy == "Retain" {
				issues = append(issues, issue(a.Name(), "PersistentVolume", "", pv.Name,
					"PVReleasedRetained", Warning,
					"PV libere avec politique Retain ??? supprimer manuellement ou recycle"))
			}
		}
	}

	// StorageClass sans provisioner valide
	scs, err := a.c.K8s.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, sc := range scs.Items {
		if sc.Provisioner == "no-provisioner" || sc.Provisioner == "kubernetes.io/no-provisioner" {
			issues = append(issues, issue(a.Name(), "StorageClass", "", sc.Name,
				"NoProvisioner", Info,
				fmt.Sprintf("StorageClass '%s' sans provisioner automatique ??? PVC a binder manuellement", sc.Name)))
		}
	}
	return issues, nil
}
