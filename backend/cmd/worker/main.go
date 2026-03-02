package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/k8s-green/backend/internal/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Erreur chargement config")
	}

	level, _ := zerolog.ParseLevel(cfg.LogLevel)
	zerolog.SetGlobalLevel(level)

	db, err := repository.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Connexion PostgreSQL échouée")
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh
		log.Info().Msg("Signal reçu — arrêt gracieux worker...")
		cancel()
	}()

	consumer, err := worker.NewConsumer(cfg, db)
	if err != nil {
		log.Fatal().Err(err).Msg("Création consumer NATS échouée")
	}

	log.Info().Msg("Worker démarré")
	consumer.Start(ctx)
}
