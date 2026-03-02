package analyzers

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Config active/desactive les analyseurs optionnels
type Config struct {
	EnableHPA            bool
	EnablePDB            bool
	EnableNetworkPolicy  bool
	EnableSecurity       bool
	EnableLog            bool
	EnableStorage        bool
	EnableGateway        bool
	EnableOLM            bool
}

// DefaultConfig active les analyseurs standards
func DefaultConfig() Config {
	return Config{
		EnableHPA:           true,
		EnablePDB:           true,
		EnableNetworkPolicy: false, // beaucoup de faux positifs sans CNI
		EnableSecurity:      true,
		EnableLog:           true,
		EnableStorage:       true,
		EnableGateway:       false, // Gateway API non installee par defaut
		EnableOLM:           false, // OLM non installe par defaut
	}
}

// Result resultat de tous les analyseurs
type Result struct {
	ClusterID   string    `json:"cluster_id"`
	TenantID    string    `json:"tenant_id"`
	CollectedAt time.Time `json:"collected_at"`
	Issues      []Issue   `json:"issues"`
	HealthScore float64   `json:"health_score"`
	Analyzers   []string  `json:"analyzers_run"`
	Duration    string    `json:"duration"`
}

// Registry lance tous les analyseurs en parallele
type Registry struct {
	analyzers []Analyzer
}

func NewRegistry(clients *Clients, cfg Config) *Registry {
	r := &Registry{}

	// Core analyzers (toujours actifs)
	r.analyzers = []Analyzer{
		NewPodAnalyzer(clients),
		NewPVCAnalyzer(clients),
		NewRSAnalyzer(clients),
		NewServiceAnalyzer(clients),
		NewEventAnalyzer(clients),
		NewIngressAnalyzer(clients),
		NewStatefulSetAnalyzer(clients),
		NewDeploymentAnalyzer(clients),
		NewJobAnalyzer(clients),
		NewCronJobAnalyzer(clients),
		NewNodeAnalyzer(clients),
		NewMutatingWebhookAnalyzer(clients),
		NewValidatingWebhookAnalyzer(clients),
		NewConfigMapAnalyzer(clients),
	}

	// Optional analyzers
	if cfg.EnableHPA           { r.analyzers = append(r.analyzers, NewHPAAnalyzer(clients)) }
	if cfg.EnablePDB           { r.analyzers = append(r.analyzers, NewPDBAnalyzer(clients)) }
	if cfg.EnableNetworkPolicy { r.analyzers = append(r.analyzers, NewNetworkPolicyAnalyzer(clients)) }
	if cfg.EnableSecurity      { r.analyzers = append(r.analyzers, NewSecurityAnalyzer(clients)) }
	if cfg.EnableLog           { r.analyzers = append(r.analyzers, NewLogAnalyzer(clients)) }
	if cfg.EnableStorage       { r.analyzers = append(r.analyzers, NewStorageAnalyzer(clients)) }
	if cfg.EnableGateway {
		r.analyzers = append(r.analyzers,
			NewGatewayClassAnalyzer(clients),
			NewGatewayAnalyzer(clients),
			NewHTTPRouteAnalyzer(clients),
		)
	}
	if cfg.EnableOLM { r.analyzers = append(r.analyzers, NewOLMAnalyzer(clients)) }

	return r
}

// Run lance tous les analyseurs en parallele et aggrege les resultats
func (r *Registry) Run(ctx context.Context, clusterID, tenantID string) *Result {
	start := time.Now()
	result := &Result{
		ClusterID:   clusterID,
		TenantID:    tenantID,
		CollectedAt: start,
	}

	type partial struct {
		name   string
		issues []Issue
		err    error
	}

	ch := make(chan partial, len(r.analyzers))
	var wg sync.WaitGroup

	for _, a := range r.analyzers {
		wg.Add(1)
		go func(a Analyzer) {
			defer wg.Done()
			issues, err := a.Analyze(ctx)
			if err != nil {
				log.Warn().Str("analyzer", a.Name()).Err(err).Msg("Analyzer error")
			}
			ch <- partial{name: a.Name(), issues: issues, err: err}
		}(a)
	}

	wg.Wait()
	close(ch)

	for p := range ch {
		result.Analyzers = append(result.Analyzers, p.name)
		result.Issues = append(result.Issues, p.issues...)
	}

	result.HealthScore = computeScore(result.Issues)
	result.Duration = time.Since(start).Round(time.Millisecond).String()

	log.Info().
		Str("cluster", clusterID).
		Int("issues", len(result.Issues)).
		Float64("health", result.HealthScore).
		Str("duration", result.Duration).
		Msg("Diagnostics complets")

	return result
}

func computeScore(issues []Issue) float64 {
	if len(issues) == 0 { return 100 }
	penalty := 0.0
	for _, i := range issues {
		switch i.Severity {
		case Critical: penalty += 15
		case Warning:  penalty += 5
		case Info:     penalty += 1
		}
	}
	if s := 100 - penalty; s >= 0 { return s }
	return 0
}
