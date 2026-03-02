package publisher
import "time"
type MetricsSnapshot struct {
	ClusterID string `json:"cluster_id"`; TenantID string `json:"tenant_id"`
	CollectedAt time.Time `json:"collected_at"`
	Pods []PodMetrics `json:"pods"`; Nodes []NodeMetrics `json:"nodes"`
	Deployments []DeploymentMetrics `json:"deployments"`; Namespaces []NamespaceMetrics `json:"namespaces"`
}
type PodMetrics struct {
	PodName, ContainerName, Namespace, NodeName string
	Labels map[string]string
	CPURequestM, CPULimitM, MemRequestMi, MemLimitMi float64
	CPUUsageP95M, MemUsageP95Mi float64
	HasLimits, IsHPAManaged bool; RestartCount int32
}
type NodeMetrics struct {
	Name string; Labels map[string]string
	CPUAllocatableM, MemAllocatableMi float64
	CPUCapacityM, MemCapacityMi float64
	InstanceType, Region, Zone string
}
type DeploymentMetrics struct {
	Name, Namespace string; Labels map[string]string
	Replicas *int32; ReadyReplicas int32; HasHPA bool
}
type NamespaceMetrics struct {
	Name string; Labels map[string]string
	HasResourceQuota, HasLimitRange bool
}
