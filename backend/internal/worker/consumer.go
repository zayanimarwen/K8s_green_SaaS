package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// Consumer consomme les snapshots d'agents via NATS JetStream
type Consumer struct {
	cfg       *config.Config
	db        *repository.Postgres
	nc        *nats.Conn
	js        nats.JetStreamContext
	processor *Processor
}

func NewConsumer(cfg *config.Config, db *repository.Postgres) (*Consumer, error) {
	nc, err := nats.Connect(cfg.NATSUrl,
		nats.Name("green-backend-worker"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Warn().Err(err).Msg("NATS déconnecté")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info().Str("url", nc.ConnectedUrl()).Msg("NATS reconnecté")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("jetstream init: %w", err)
	}

	// Créer le stream si absent (idempotent)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:       "METRICS",
		Subjects:   []string{"metrics.>"},
		MaxAge:     7 * 24 * time.Hour, // Rétention 7 jours
		Storage:    nats.FileStorage,
		Replicas:   1,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Stream METRICS déjà existant ou erreur")
	}

	return &Consumer{
		cfg:       cfg,
		db:        db,
		nc:        nc,
		js:        js,
		processor: NewProcessor(cfg, db),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	log.Info().Msg("Consumer démarré — écoute metrics.>")

	sub, err := c.js.QueueSubscribe(
		"metrics.>",
		"backend-workers",
		c.handle,
		nats.Durable("backend-worker"),
		nats.AckExplicit(),
		nats.MaxDeliver(3),
		nats.AckWait(60*time.Second),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Subscribe failed")
	}
	defer sub.Unsubscribe()

	// Lancer le scheduler en parallèle
	sched := NewScheduler(c.db)
	go sched.Start(ctx)

	<-ctx.Done()
	log.Info().Msg("Consumer arrêté")
	c.nc.Drain()
}

func (c *Consumer) handle(msg *nats.Msg) {
	// Vérifier signature HMAC
	sig := msg.Header.Get("X-Signature")
	if sig == "" {
		log.Warn().Str("subject", msg.Subject).Msg("Message sans signature — ignoré")
		msg.Nak()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.processor.Process(ctx, msg.Data); err != nil {
		log.Error().Err(err).Str("subject", msg.Subject).Msg("Process failed")
		msg.Nak()
		return
	}

	msg.Ack()
}
