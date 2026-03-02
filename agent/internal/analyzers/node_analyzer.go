package analyzers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeAnalyzer struct{ c *Clients }
func NewNodeAnalyzer(c *Clients) *NodeAnalyzer { return &NodeAnalyzer{c} }
func (a *NodeAnalyzer) Name() string   { return "nodeAnalyzer" }
func (a *NodeAnalyzer) Optional() bool { return false }

func (a *NodeAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	nodes, err := a.c.K8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }

	for _, node := range nodes.Items {
		name := node.Name

		// Node conditions
		for _, cond := range node.Status.Conditions {
			switch cond.Type {
			case corev1.NodeReady:
				if cond.Status != corev1.ConditionTrue {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"NodeNotReady", Critical,
						fmt.Sprintf("Node NotReady: %s", cond.Message)))
				}
			case corev1.NodeMemoryPressure:
				if cond.Status == corev1.ConditionTrue {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"MemoryPressure", Critical,
						"Node sous pression memoire ??? eviction de pods possible"))
				}
			case corev1.NodeDiskPressure:
				if cond.Status == corev1.ConditionTrue {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"DiskPressure", Critical,
						"Node sous pression disque ??? nettoyer les images et logs"))
				}
			case corev1.NodePIDPressure:
				if cond.Status == corev1.ConditionTrue {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"PIDPressure", Warning,
						"Node avec trop de processus ??? risque de saturation"))
				}
			}
		}

		// Node cordonne (pod scheduling desactive)
		if node.Spec.Unschedulable {
			issues = append(issues, issue(a.Name(), "Node", "", name,
				"NodeCordoned", Warning,
				"Node marque Unschedulable (cordon) ??? nouveau scheduling impossible"))
		}

		// Utilisation CPU/Memoire elevee
		allocatable := node.Status.Allocatable
		requests := node.Status.Capacity

		if cpu, ok := allocatable[corev1.ResourceCPU]; ok {
			if req, ok2 := requests[corev1.ResourceCPU]; ok2 {
				pct := float64(req.MilliValue()) / float64(cpu.MilliValue()) * 100
				if pct > 90 {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"HighCPUUsage", Warning,
						fmt.Sprintf("Node a %.0f%% CPU utilise", pct)))
				}
			}
		}

		if mem, ok := allocatable[corev1.ResourceMemory]; ok {
			if req, ok2 := requests[corev1.ResourceMemory]; ok2 {
				_ = resource.MustParse("1Mi") // import force
				pct := float64(req.Value()) / float64(mem.Value()) * 100
				if pct > 90 {
					issues = append(issues, issue(a.Name(), "Node", "", name,
						"HighMemoryUsage", Warning,
						fmt.Sprintf("Node a %.0f%% memoire utilisee", pct)))
				}
			}
		}
	}
	return issues, nil
}
