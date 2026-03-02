package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	RedisURL    string
	NATSUrl     string
	CORSOrigins string
	SigningKey   string
	LogLevel     string

	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
	JWTIssuer            string

	SlackWebhook string
	TeamsWebhook string
	SMTPHost     string
	SMTPPort     int
	SMTPFrom     string
	SMTPPassword string

	DefaultCostCPUH     float64
	DefaultCostMemH     float64
	ScoreAlertThreshold float64

	// AI / LLM — Ollama par defaut
	AIBackend  string // ollama | openai | claude  (defaut: ollama)
	AIAPIKey   string // requis seulement pour openai et claude
	AIModel    string // defaut: llama3.2 pour ollama, gpt-4o pour openai
	OllamaURL  string // defaut: http://localhost:11434
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetDefault("PORT",                  "9000")
	v.SetDefault("ENV",                   "development")
	v.SetDefault("LOG_LEVEL",             "info")
	v.SetDefault("CORS_ORIGINS",          "*")
	v.SetDefault("DEFAULT_COST_CPU_H",    0.05)
	v.SetDefault("DEFAULT_COST_MEM_H",    0.006)
	v.SetDefault("SCORE_ALERT_THRESHOLD", 50.0)
	v.SetDefault("SMTP_PORT",             587)
	// Ollama comme backend IA par defaut
	v.SetDefault("AI_BACKEND",  "ollama")
	v.SetDefault("AI_MODEL",    "llama3.2")
	v.SetDefault("OLLAMA_URL",  "http://host.docker.internal:11434")
	v.AutomaticEnv()

	for _, key := range []string{"DATABASE_URL", "REDIS_URL", "NATS_URL"} {
		if v.GetString(key) == "" {
			return nil, fmt.Errorf("variable obligatoire manquante: %s", key)
		}
	}

	return &Config{
		Env:                  v.GetString("ENV"),
		Port:                 v.GetString("PORT"),
		LogLevel:             v.GetString("LOG_LEVEL"),
		DatabaseURL:          v.GetString("DATABASE_URL"),
		RedisURL:             v.GetString("REDIS_URL"),
		NATSUrl:              v.GetString("NATS_URL"),
		CORSOrigins:          v.GetString("CORS_ORIGINS"),
		SigningKey:           v.GetString("SIGNING_KEY"),
		KeycloakURL:          v.GetString("KEYCLOAK_URL"),
		KeycloakRealm:        v.GetString("KEYCLOAK_REALM"),
		KeycloakClientID:     v.GetString("KEYCLOAK_CLIENT_ID"),
		KeycloakClientSecret: v.GetString("KEYCLOAK_CLIENT_SECRET"),
		JWTIssuer:            v.GetString("JWT_ISSUER"),
		SlackWebhook:         v.GetString("SLACK_WEBHOOK"),
		TeamsWebhook:         v.GetString("TEAMS_WEBHOOK"),
		SMTPHost:             v.GetString("SMTP_HOST"),
		SMTPPort:             v.GetInt("SMTP_PORT"),
		SMTPFrom:             v.GetString("SMTP_FROM"),
		SMTPPassword:         v.GetString("SMTP_PASSWORD"),
		DefaultCostCPUH:      v.GetFloat64("DEFAULT_COST_CPU_H"),
		DefaultCostMemH:      v.GetFloat64("DEFAULT_COST_MEM_H"),
		ScoreAlertThreshold:  v.GetFloat64("SCORE_ALERT_THRESHOLD"),
		AIBackend:            v.GetString("AI_BACKEND"),
		AIAPIKey:             v.GetString("AI_API_KEY"),
		AIModel:              v.GetString("AI_MODEL"),
		OllamaURL:            v.GetString("OLLAMA_URL"),
	}, nil
}
