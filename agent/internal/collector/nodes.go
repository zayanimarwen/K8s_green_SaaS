package collector
import (
	"context"
	"github.com/k8s-green/agent/internal/publisher"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)
type NodeCollector struct { k8s kubernetes.Interface }
func NewNodeCollector(k8s kubernetes.Interface) *NodeCollector { return &NodeCollector{k8s: k8s} }
func (nc *NodeCollector) Collect(ctx context.Context) ([]publisher.NodeMetrics, error) {
	nodes, err := nc.k8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }
	var result []publisher.NodeMetrics
	for _, n := range nodes.Items {
		alloc := n.Status.Allocatable
		result = append(result, publisher.NodeMetrics{
			Name: n.Name, Labels: n.Labels,
			CPUAllocatableM: float64(alloc.Cpu().MilliValue()),
			MemAllocatableMi: float64(alloc.Memory().Value()) / 1048576,
			CPUCapacityM: float64(n.Status.Capacity.Cpu().MilliValue()),
			MemCapacityMi: float64(n.Status.Capacity.Memory().Value()) / 1048576,
			InstanceType: n.Labels["node.kubernetes.io/instance-type"],
			Region: n.Labels["topology.kubernetes.io/region"],
			Zone: n.Labels["topology.kubernetes.io/zone"],
		})
	}
	return result, nil
}
