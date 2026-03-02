package repository

import (
	"context"
	"time"
)

// MetricsRepo persiste les métriques pod brutes dans TimescaleDB
type MetricsRepo struct{ db *Postgres }

type PodMetricRow struct {
	Time          time.Time
	ClusterID     string
	PodName       string
	ContainerName string
	Namespace     string
	NodeName      string
	CPURequestM   float64
	CPULimitM     float64
	MemRequestMi  float64
	MemLimitMi    float64
	CPUUsageP95M  float64
	MemUsageP95Mi float64
	CPUWasteM     float64
	MemWasteMi    float64
	CostWasteEur  float64
	HasLimits     bool
	RestartCount  int
}

func (r *Postgres) MetricsRepo() *MetricsRepo { return &MetricsRepo{r} }

func (m *MetricsRepo) BulkInsert(ctx context.Context, rows []PodMetricRow) error {
	if len(rows) == 0 {
		return nil
	}
	tx, err := m.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, r := range rows {
		_, err := tx.Exec(ctx,
			`INSERT INTO pod_metrics
			 (time, cluster_id, pod_name, container_name, namespace, node_name,
			  cpu_request_m, cpu_limit_m, mem_request_mi, mem_limit_mi,
			  cpu_usage_p95_m, mem_usage_p95_mi, cpu_waste_m, mem_waste_mi,
			  cost_waste_eur, has_limits, restart_count)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
			r.Time, r.ClusterID, r.PodName, r.ContainerName, r.Namespace, r.NodeName,
			r.CPURequestM, r.CPULimitM, r.MemRequestMi, r.MemLimitMi,
			r.CPUUsageP95M, r.MemUsageP95Mi, r.CPUWasteM, r.MemWasteMi,
			r.CostWasteEur, r.HasLimits, r.RestartCount,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (m *MetricsRepo) DeleteOlderThan(ctx context.Context, duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	_, err := m.db.Pool.Exec(ctx,
		`DELETE FROM pod_metrics WHERE time < $1`, cutoff,
	)
	return err
}
