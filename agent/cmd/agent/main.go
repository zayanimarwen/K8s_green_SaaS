package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k8s-green/agent/internal/collector"
	"github.com/k8s-green/agent/internal/config"
	"github.com/k8s-green/agent/internal/publisher"
	"github.com/k8s-green/agent/internal/security"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var version = "dev"

func main() {
	// ── Logging structuré JSON ──────────────────────────────────────
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	log.Info().
		Str("version", version).
		Str("component", "agent").
		Msg("K8s Green Agent démarrage")

	// ── Configuration ───────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Échec chargement configuration")
	}

	log.Info().
		Str("cluster_id", cfg.ClusterID).
		Str("tenant_id", cfg.TenantID).
		Str("saas_url", cfg.SaaSURL).
		Str("schedule", cfg.CollectInterval.String()).
		Msg("Configuration chargée")

	// ── Sécurité mTLS + Signer ──────────────────────────────────────
	tlsConfig, err := security.LoadMTLS(cfg.TLSCertFile, cfg.TLSKeyFile, cfg.TLSCAFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Échec chargement certificats mTLS")
	}

	signer, err := security.NewSigner(cfg.SigningKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Échec init signer HMAC")
	}

	// ── Publisher NATS ──────────────────────────────────────────────
	pub, err := publisher.NewNATSPublisher(cfg.NATSUrl, cfg.TenantID, cfg.ClusterID, tlsConfig, signer)
	if err != nil {
		log.Fatal().Err(err).Msg("Échec connexion NATS")
	}
	defer pub.Close()

	// ── Collector ───────────────────────────────────────────────────
	col, err := collector.New(cfg, pub)
	if err != nil {
		log.Fatal().Err(err).Msg("Échec init collector")
	}

	// ── Graceful shutdown ───────────────────────────────────────────
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("Signal reçu — arrêt gracieux")
		cancel()
	}()

	// ── Démarrage de la boucle de collecte ──────────────────────────
	log.Info().Msg("Agent démarré — début collecte")
	if err := col.Run(ctx); err != nil && err != context.Canceled {
		log.Error().Err(err).Msg("Erreur collector")
		os.Exit(1)
	}

	log.Info().Msg("Agent arrêté proprement")
}
