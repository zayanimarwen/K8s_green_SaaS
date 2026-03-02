package collector
import (
	"context"
	"github.com/k8s-green/agent/internal/publisher"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)
type NamespaceCollector struct { k8s kubernetes.Interface }
func NewNamespaceCollector(k8s kubernetes.Interface) *NamespaceCollector {
	return &NamespaceCollector{k8s: k8s}
}
func (nc *NamespaceCollector) Collect(ctx context.Context) ([]publisher.NamespaceMetrics, error) {
	nsList, err := nc.k8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }
	var result []publisher.NamespaceMetrics
	for _, ns := range nsList.Items {
		if systemNS[ns.Name] { continue }
		nm := publisher.NamespaceMetrics{Name: ns.Name, Labels: ns.Labels}
		q, _ := nc.k8s.CoreV1().ResourceQuotas(ns.Name).List(ctx, metav1.ListOptions{})
		nm.HasResourceQuota = len(q.Items) > 0
		l, _ := nc.k8s.CoreV1().LimitRanges(ns.Name).List(ctx, metav1.ListOptions{})
		nm.HasLimitRange = len(l.Items) > 0
		result = append(result, nm)
	}
	return result, nil
}
