package collector
import (
	"context"
	"github.com/k8s-green/agent/internal/config"
	"github.com/k8s-green/agent/internal/publisher"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)
var systemNS = map[string]bool{"kube-system": true, "kube-public": true, "kube-node-lease": true}
type PodCollector struct { k8s kubernetes.Interface; cfg *config.Config }
func NewPodCollector(k8s kubernetes.Interface, cfg *config.Config) *PodCollector {
	return &PodCollector{k8s: k8s, cfg: cfg}
}
func (pc *PodCollector) Collect(ctx context.Context) ([]publisher.PodMetrics, error) {
	pods, err := pc.k8s.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil { return nil, err }
	var result []publisher.PodMetrics
	for _, pod := range pods.Items {
		if !pc.cfg.IncludeSystemNS && systemNS[pod.Namespace] { continue }
		if pod.Status.Phase != corev1.PodRunning { continue }
		for _, c := range pod.Spec.Containers {
			req, lim := c.Resources.Requests, c.Resources.Limits
			result = append(result, publisher.PodMetrics{
				PodName: pod.Name, ContainerName: c.Name, Namespace: pod.Namespace,
				NodeName: pod.Spec.NodeName, Labels: pod.Labels,
				CPURequestM: float64(req.Cpu().MilliValue()),
				MemRequestMi: float64(req.Memory().Value()) / 1048576,
				CPULimitM: float64(lim.Cpu().MilliValue()),
				MemLimitMi: float64(lim.Memory().Value()) / 1048576,
				HasLimits: !lim.Cpu().IsZero() && !lim.Memory().IsZero(),
				RestartCount: getRestartCount(pod.Status.ContainerStatuses, c.Name),
			})
		}
	}
	return result, nil
}
func getRestartCount(statuses []corev1.ContainerStatus, name string) int32 {
	for _, s := range statuses { if s.Name == name { return s.RestartCount } }
	return 0
}
