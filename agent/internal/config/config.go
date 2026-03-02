package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config contient toute la configuration de l'agent.
// Chargée depuis variables d'environnement (priorité) ou fichier config.yaml
type Config struct {
	// Identité
	ClusterID string
	TenantID  string

	// SaaS
	NATSUrl  string
	SaaSURL  string

	// Collecte
	CollectInterval  time.Duration
	PrometheusURL    string
	PrometheusToken  string

	// Sécurité mTLS
	TLSCertFile string
	TLSKeyFile  string
	TLSCAFile   string
	SigningKey   string

	// Comportement
	LogLevel        string
	IncludeSystemNS bool
	DryRun          bool
}

func Load() (*Config, error) {
	v := viper.New()

	// Valeurs par défaut
	v.SetDefault("collect_interval", "30s")
	v.SetDefault("log_level", "info")
	v.SetDefault("include_system_ns", false)
	v.SetDefault("dry_run", false)
	v.SetDefault("prometheus_url", "http://prometheus-operated.monitoring:9090")

	// Variables d'environnement
	v.SetEnvPrefix("GREEN")
	v.AutomaticEnv()

	// Mappings explicites env → config
	envMappings := map[string]string{
		"cluster_id":       "CLUSTER_ID",
		"tenant_id":        "TENANT_ID",
		"nats_url":         "NATS_URL",
		"saas_url":         "SAAS_URL",
		"collect_interval": "COLLECT_INTERVAL",
		"prometheus_url":   "PROMETHEUS_URL",
		"prometheus_token": "PROMETHEUS_TOKEN",
		"tls_cert_file":    "TLS_CERT_FILE",
		"tls_key_file":     "TLS_KEY_FILE",
		"tls_ca_file":      "TLS_CA_FILE",
		"signing_key":      "SIGNING_KEY",
		"log_level":        "LOG_LEVEL",
		"dry_run":          "DRY_RUN",
	}

	for key, env := range envMappings {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("bind env %s: %w", env, err)
		}
	}

	// Validation des champs obligatoires
	required := []string{"cluster_id", "tenant_id", "nats_url", "signing_key"}
	for _, field := range required {
		if v.GetString(field) == "" {
			return nil, fmt.Errorf("variable obligatoire manquante: %s (GREEN_%s)", field, field)
		}
	}

	interval, err := time.ParseDuration(v.GetString("collect_interval"))
	if err != nil {
		return nil, fmt.Errorf("collect_interval invalide: %w", err)
	}

	return &Config{
		ClusterID:        v.GetString("cluster_id"),
		TenantID:         v.GetString("tenant_id"),
		NATSUrl:          v.GetString("nats_url"),
		SaaSURL:          v.GetString("saas_url"),
		CollectInterval:  interval,
		PrometheusURL:    v.GetString("prometheus_url"),
		PrometheusToken:  v.GetString("prometheus_token"),
		TLSCertFile:      v.GetString("tls_cert_file"),
		TLSKeyFile:       v.GetString("tls_key_file"),
		TLSCAFile:        v.GetString("tls_ca_file"),
		SigningKey:        v.GetString("signing_key"),
		LogLevel:         v.GetString("log_level"),
		IncludeSystemNS:  v.GetBool("include_system_ns"),
		DryRun:           v.GetBool("dry_run"),
	}, nil
}
