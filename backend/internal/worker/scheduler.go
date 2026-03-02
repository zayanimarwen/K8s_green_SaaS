package worker

import (
	"context"
	"time"

	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	db *repository.Postgres
}

func NewScheduler(db *repository.Postgres) *Scheduler {
	return &Scheduler{db: db}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Info().Msg("Scheduler démarré")

	// Cleanup métriques > 90 jours toutes les 24h
	cleanupTicker := time.NewTicker(24 * time.Hour)
	// Check rapports hebdo chaque heure
	weeklyTicker := time.NewTicker(time.Hour)

	defer cleanupTicker.Stop()
	defer weeklyTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Scheduler arrêté")
			return

		case <-cleanupTicker.C:
			log.Info().Msg("Scheduler: cleanup métriques > 90 jours")
			if err := s.db.MetricsRepo().DeleteOlderThan(ctx, 90*24*time.Hour); err != nil {
				log.Error().Err(err).Msg("Cleanup failed")
			}

		case t := <-weeklyTicker.C:
			// Générer les rapports hebdo chaque lundi à 8h UTC
			if t.Weekday() == time.Monday && t.Hour() == 8 {
				log.Info().Msg("Scheduler: génération rapports hebdomadaires")
				s.generateWeeklyReports(ctx)
			}
		}
	}
}

func (s *Scheduler) generateWeeklyReports(ctx context.Context) {
	// Récupérer tous les tenants actifs
	tenants, err := s.db.TenantRepo().List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Impossible de lister les tenants")
		return
	}
	for _, tenant := range tenants {
		log.Info().Str("tenant", tenant.ID).Msg("Génération rapport hebdo")
		// Récupérer les clusters du tenant et envoyer le rapport
		// Les emails sont envoyés via notifications.EmailNotifier
	}
}
