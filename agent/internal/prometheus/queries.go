package prometheus
// CPUUsageP95 retourne la query PromQL pour le P95 CPU sur 24h (millicores)
func CPUUsageP95() string {
	return `quantile_over_time(0.95,
		rate(container_cpu_usage_seconds_total{container!="",container!="POD"}[5m])[24h:5m]
	) * 1000`
}
// MemUsageP95 retourne la query PromQL pour le P95 mémoire sur 24h (MiB)
func MemUsageP95() string {
	return `quantile_over_time(0.95,
		container_memory_working_set_bytes{container!="",container!="POD"}[24h:5m]
	) / 1048576`
}
// CPUThrottlingRate retourne le taux de throttling
func CPUThrottlingRate() string {
	return `sum by (pod,namespace,container)(rate(container_cpu_cfs_throttled_seconds_total[5m])) /
		sum by (pod,namespace,container)(rate(container_cpu_cfs_periods_total[5m]))`
}
