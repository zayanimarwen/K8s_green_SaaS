package domain

import "time"

type PodMetric struct {
	Time          time.Time `json:"time"`
	ClusterID     string    `json:"cluster_id"`
	PodName       string    `json:"pod_name"`
	ContainerName string    `json:"container_name"`
	Namespace     string    `json:"namespace"`
	CPURequestM   float64   `json:"cpu_request_m"`
	MemRequestMi  float64   `json:"mem_request_mi"`
	CPUUsageP95M  float64   `json:"cpu_usage_p95_m"`
	MemUsageP95Mi float64   `json:"mem_usage_p95_mi"`
	HasLimits     bool      `json:"has_limits"`
	RestartCount  int       `json:"restart_count"`
}
