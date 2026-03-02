package repository

import (
	"context"
	"time"

	"github.com/k8s-green/backend/internal/domain"
)

// TimescaleRepo gère les métriques time-series dans TimescaleDB
type TimescaleRepo struct{ db *Postgres }

// InsertScore persiste un Green Score dans la hypertable
func (r *TimescaleRepo) InsertScore(ctx context.Context, clusterID string, score domain.GreenScore) error {
	_, err := r.db.Pool.Exec(ctx,
		`INSERT INTO green_scores
		 (time, cluster_id, score, grade, cpu_eff, mem_eff, node_packing, hpa_coverage, limit_comp,
		  pod_count, waste_eur_annual, co2_kg_annual)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		time.Now(), clusterID,
		score.Score, score.Grade,
		score.CPUEfficiency, score.MemEfficiency, score.NodePacking,
		score.HPACoverage, score.LimitCompliance,
		score.PodsAnalyzed, score.AnnualWasteEur, score.CO2KgAnnual,
	)
	return err
}

// GetScoreHistory retourne l'historique des scores sur N jours
func (r *TimescaleRepo) GetScoreHistory(ctx context.Context, clusterID string, days int) ([]domain.GreenScore, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT time, cluster_id, score, grade, cpu_eff, mem_eff, node_packing, hpa_coverage, limit_comp
		 FROM green_scores
		 WHERE cluster_id = $1 AND time > NOW() - INTERVAL '1 day' * $2
		 ORDER BY time DESC`,
		clusterID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []domain.GreenScore
	for rows.Next() {
		var s domain.GreenScore
		var t time.Time
		rows.Scan(&t, &s.ClusterID, &s.Score, &s.Grade,
			&s.CPUEfficiency, &s.MemEfficiency, &s.NodePacking, &s.HPACoverage, &s.LimitCompliance)
		s.Time = t
		scores = append(scores, s)
	}
	return scores, nil
}

// GetLatestScore retourne le dernier score d'un cluster (avec cache possible)
func (r *TimescaleRepo) GetLatestScore(ctx context.Context, clusterID string) (*domain.GreenScore, error) {
	var s domain.GreenScore
	err := r.db.Pool.QueryRow(ctx,
		`SELECT time, cluster_id, score, grade, cpu_eff, mem_eff, node_packing, hpa_coverage, limit_comp,
		        pod_count, waste_eur_annual, co2_kg_annual
		 FROM green_scores WHERE cluster_id = $1 ORDER BY time DESC LIMIT 1`,
		clusterID,
	).Scan(&s.Time, &s.ClusterID, &s.Score, &s.Grade,
		&s.CPUEfficiency, &s.MemEfficiency, &s.NodePacking, &s.HPACoverage, &s.LimitCompliance,
		&s.PodsAnalyzed, &s.AnnualWasteEur, &s.CO2KgAnnual)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
