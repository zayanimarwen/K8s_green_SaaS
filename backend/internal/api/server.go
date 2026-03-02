package api

import (
	"context"
	"net/http"
	"time"

	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg *config.Config
	db  *repository.Postgres
	rdb *repository.Redis
}

func NewServer(cfg *config.Config, db *repository.Postgres, rdb *repository.Redis) *Server {
	return &Server{cfg: cfg, db: db, rdb: rdb}
}

func (s *Server) Run(ctx context.Context) error {
	router := NewRouter(s.cfg, s.db, s.rdb)

	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown sur ctx
	go func() {
		<-ctx.Done()
		log.Info().Msg("Arrêt serveur API...")
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutCtx)
	}()

	log.Info().Str("port", s.cfg.Port).Str("env", s.cfg.Env).Msg("API démarrée")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
