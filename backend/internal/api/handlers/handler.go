package handlers

import (
	"github.com/k8s-green/backend/internal/ai"
	"github.com/k8s-green/backend/internal/config"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	db    *repository.Postgres
	rdb   *repository.Redis
	cfg   *config.Config
	aiSvc *ai.DiagnosticsService
}

func NewHandler(db *repository.Postgres, rdb *repository.Redis, cfg *config.Config) *Handler {
	h := &Handler{db: db, rdb: rdb, cfg: cfg}

	// Initialiser le service LLM si configure
	if cfg.AIBackend != "" {
		aiConfig := ai.Config{
			Backend:   ai.Backend(cfg.AIBackend),
			APIKey:    cfg.AIAPIKey,
			Model:     cfg.AIModel,
			OllamaURL: cfg.OllamaURL,
			MaxTokens: 1500,
		}
		client, err := ai.NewClient(aiConfig)
		if err != nil {
			log.Warn().Err(err).Msg("LLM non initialise — diagnostics IA desactives")
		} else {
			h.aiSvc = ai.NewDiagnosticsService(client, rdb)
			log.Info().Str("backend", cfg.AIBackend).Str("model", cfg.AIModel).Msg("LLM initialise")
		}
	}
	return h
}
