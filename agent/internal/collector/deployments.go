package collector
import (
	"context"
	"github.com/k8s-green/agent/internal/config"
	"github.com/k8s-green/agent/internal/publisher"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)
type DeploymentCollector struct { k8s kubernetes.Interface; cfg *config.Config }
func NewDeploymentCollector(k8s kubernetes.Interface, cfg *config.Config) *DeploymentCollector {
	return &DeploymentCollector{k8s: k8s, cfg: cfg}
}
func (dc *DeploymentCollector) Collect(ctx context.Context) ([]publisher.DeploymentMetrics, error) {
	deps, _ := dc.k8s.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	hpas, _ := dc.k8s.AutoscalingV2().HorizontalPodAutoscalers("").List(ctx, metav1.ListOptions{})
	hpaSet := map[string]bool{}
	for _, h := range hpas.Items { hpaSet[h.Namespace+"/"+h.Spec.ScaleTargetRef.Name] = true }
	var result []publisher.DeploymentMetrics
	for _, d := range deps.Items {
		if !dc.cfg.IncludeSystemNS && systemNS[d.Namespace] { continue }
		result = append(result, publisher.DeploymentMetrics{
			Name: d.Name, Namespace: d.Namespace, Labels: d.Labels,
			Replicas: d.Spec.Replicas, ReadyReplicas: d.Status.ReadyReplicas,
			HasHPA: hpaSet[d.Namespace+"/"+d.Name],
		})
	}
	return result, nil
}
