package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/k8s-green/backend/internal/analyzer"
	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/domain"
	"github.com/k8s-green/backend/internal/notifications"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

type Processor struct {
	cfg   *config.Config
	db    *repository.Postgres
	slack *notifications.SlackNotifier
}

func NewProcessor(cfg *config.Config, db *repository.Postgres) *Processor {
	var slack *notifications.SlackNotifier
	if cfg.SlackWebhook != "" {
		slack = notifications.NewSlack(cfg.SlackWebhook)
	}
	return &Processor{cfg: cfg, db: db, slack: slack}
}

// incomingSnapshot est le payload JSON reçu de l'agent
type incomingSnapshot struct {
	ClusterID   string    `json:"cluster_id"`
	TenantID    string    `json:"tenant_id"`
	CollectedAt time.Time `json:"collected_at"`
	Pods []struct {
		PodName, ContainerName, Namespace string
		CPURequestM, CPUUsageP95M         float64
		MemRequestMi, MemUsageP95Mi       float64
		HasLimits, IsHPAManaged           bool
		RestartCount                      int32
	} `json:"pods"`
}

func (p *Processor) Process(ctx context.Context, data []byte) error {
	var snap incomingSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}

	log.Info().
		Str("cluster", snap.ClusterID).
		Str("tenant", snap.TenantID).
		Int("pods", len(snap.Pods)).
		Msg("Traitement snapshot")

	// 1. Convertir en PodInput
	var pods []analyzer.PodInput
	for _, pod := range snap.Pods {
		pods = append(pods, analyzer.PodInput{
			PodName:       pod.PodName,
			ContainerName: pod.ContainerName,
			Namespace:     pod.Namespace,
			CPURequestM:   pod.CPURequestM,
			CPUUsageP95M:  pod.CPUUsageP95M,
			MemRequestMi:  pod.MemRequestMi,
			MemUsageP95Mi: pod.MemUsageP95Mi,
			HasLimits:     pod.HasLimits,
			IsHPAManaged:  pod.IsHPAManaged,
			RestartCount:  pod.RestartCount,
		})
	}

	// 2. Détection gaspillage
	detector := analyzer.NewWasteDetector(p.cfg.DefaultCostCPUH, p.cfg.DefaultCostMemH)
	var reports []analyzer.WasteReport
	for _, pod := range pods {
		reports = append(reports, detector.Analyze(pod))
	}

	// 3. Score global
	score := analyzer.ComputeScore(pods, 0, len(pods))

	// 4. Carbone
	cc := analyzer.NewCarbonCalculator("on-prem", "fr")
	carbon := cc.Compute(reports)

	// 5. Coût gaspillage total
	costReport := analyzer.ComputeSavings(reports)

	log.Info().
		Str("cluster", snap.ClusterID).
		Float64("score", score.Value).
		Str("grade", score.Grade).
		Float64("co2_kg", carbon.CO2KgAnnual).
		Float64("waste_eur", costReport.AnnualSavingsEur).
		Msg("Analyse terminée")

	// 6. Persister le score dans TimescaleDB
	domainScore := domain.GreenScore{
		ClusterID:      snap.ClusterID,
		Time:           time.Now(),
		Score:          score.Value,
		Grade:          score.Grade,
		Label:          score.Label,
		CPUEfficiency:  score.Breakdown.CPUEfficiency,
		MemEfficiency:  score.Breakdown.MemEfficiency,
		NodePacking:    score.Breakdown.NodePacking,
		HPACoverage:    score.Breakdown.HPACoverage,
		LimitCompliance: score.Breakdown.LimitCompliance,
		PodsAnalyzed:   len(pods),
		AnnualWasteEur: costReport.AnnualSavingsEur,
		CO2KgAnnual:    carbon.CO2KgAnnual,
	}
	if err := p.db.TimescaleRepo().InsertScore(ctx, snap.ClusterID, domainScore); err != nil {
		log.Error().Err(err).Str("cluster", snap.ClusterID).Msg("Échec persistance score")
		// Non-fatal : on continue quand même
	}

	// 7. Persister les métriques pod brutes
	var metricRows []repository.PodMetricRow
	for i, pod := range snap.Pods {
		if i >= len(reports) { break }
		r := reports[i]
		metricRows = append(metricRows, repository.PodMetricRow{
			Time:          snap.CollectedAt,
			ClusterID:     snap.ClusterID,
			PodName:       pod.PodName,
			ContainerName: pod.ContainerName,
			Namespace:     pod.Namespace,
			CPURequestM:   pod.CPURequestM,
			MemRequestMi:  pod.MemRequestMi,
			CPUUsageP95M:  pod.CPUUsageP95M,
			MemUsageP95Mi: pod.MemUsageP95Mi,
			CPUWasteM:     r.CPUWasteM,
			MemWasteMi:    r.MemWasteMi,
			CostWasteEur:  r.AnnualCostWasteEur,
			HasLimits:     pod.HasLimits,
			RestartCount:  int(pod.RestartCount),
		})
	}
	if len(metricRows) > 0 {
		if err := p.db.MetricsRepo().BulkInsert(ctx, metricRows); err != nil {
			log.Error().Err(err).Msg("Échec BulkInsert metrics")
		}
	}

	// 8. Alerte Slack si score sous le seuil
	if p.slack != nil && score.Value < p.cfg.ScoreAlertThreshold {
		p.slack.SendScanReport(snap.TenantID, snap.ClusterID, score.Value, score.Grade, costReport.AnnualSavingsEur)
	}

	// 9. Mettre à jour last_seen du cluster
	p.db.ClusterRepo().UpdateLastSeen(ctx, snap.ClusterID)

	return nil
}
