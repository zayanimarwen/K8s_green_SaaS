package ai

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Backend string

const (
	BackendOpenAI  Backend = "openai"
	BackendClaude  Backend = "claude"
	BackendOllama  Backend = "ollama"
)

type LLMClient interface {
	Analyze(ctx context.Context, prompt string) (string, error)
	BackendName() Backend
}

type Config struct {
	Backend   Backend
	APIKey    string
	Model     string
	OllamaURL string
	MaxTokens int
}

// NewClient cree le client LLM ??? Ollama est le backend par defaut
func NewClient(cfg Config) (LLMClient, error) {
	// Defaut : Ollama si rien n est configure
	if cfg.Backend == "" {
		cfg.Backend = BackendOllama
		log.Info().Msg("AI backend non specifie ??? utilisation d Ollama (local)")
	}

	switch cfg.Backend {
	case BackendOllama:
		url := cfg.OllamaURL
		if url == "" { url = "http://localhost:11434" }
		model := cfg.Model
		if model == "" { model = "llama3.2" }
		client := NewOllamaClient(url, model, cfg.MaxTokens)

		// Verifier disponibilite et auto-pull si necessaire
		ctx := context.Background()
		ok := client.IsAvailable(ctx)
		if !ok {
			log.Warn().Str("url", url).Msg("Ollama non joignable - diagnostics IA desactives")
			return nil, fmt.Errorf("ollama non joignable sur %s", url)
		}
		if err := client.PullModel(ctx, model); err != nil {
			log.Warn().Err(err).Str("model", model).Msg("Pull modele echoue")
		}

		log.Info().Str("model", model).Str("url", url).Msg("Ollama pret")
		return client, nil

	case BackendOpenAI:
		if cfg.APIKey == "" { return nil, fmt.Errorf("AI_API_KEY (OpenAI) manquant") }
		model := cfg.Model
		if model == "" { model = "gpt-4o" }
		return NewOpenAIClient(cfg.APIKey, model, cfg.MaxTokens), nil

	case BackendClaude:
		if cfg.APIKey == "" { return nil, fmt.Errorf("AI_API_KEY (Anthropic) manquant") }
		model := cfg.Model
		if model == "" { model = "claude-3-5-sonnet-20241022" }
		return NewClaudeClient(cfg.APIKey, model, cfg.MaxTokens), nil

	default:
		return nil, fmt.Errorf("backend inconnu: %s (valeurs: ollama, openai, claude)", cfg.Backend)
	}
}
