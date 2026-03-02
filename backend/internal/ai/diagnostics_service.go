package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

type DiagnosticsResult struct {
	ID          string  `json:"id"`
	ClusterID   string  `json:"cluster_id"`
	TenantID    string  `json:"tenant_id"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
	Issues      []Issue `json:"issues"`
	HealthScore float64 `json:"health_score"`
	Summary     string  `json:"summary"`
	Backend     string  `json:"llm_backend"`
	Model       string  `json:"llm_model"`
}

type Issue struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Severity      string    `json:"severity"`
	Namespace     string    `json:"namespace"`
	ResourceKind  string    `json:"resource_kind"`
	ResourceName  string    `json:"resource_name"`
	Message       string    `json:"message"`
	AIExplanation string    `json:"ai_explanation,omitempty"`
	DetectedAt    time.Time `json:"detected_at"`
}

type DiagnosticsService struct {
	mu     sync.RWMutex
	llm    LLMClient
	ollama *OllamaClient // reference directe pour gestion modeles
	rdb    *repository.Redis
}

func NewDiagnosticsService(llm LLMClient, rdb *repository.Redis) *DiagnosticsService {
	svc := &DiagnosticsService{llm: llm, rdb: rdb}
	// Conserver reference Ollama si applicable
	if o, ok := llm.(*OllamaClient); ok {
		svc.ollama = o
	}
	return svc
}

// OllamaClient retourne le client Ollama si actif
func (s *DiagnosticsService) OllamaClient() *OllamaClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ollama
}

// SetModel change le modele Ollama actif a chaud (sans redemarrage)
func (s *DiagnosticsService) SetModel(modelName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ollama != nil {
		s.ollama.model = modelName
		log.Info().Str("model", modelName).Msg("Modele Ollama change a chaud")
	}
}

func (s *DiagnosticsService) AnalyzeIssues(ctx context.Context, clusterID, tenantID string, rawIssues []IssueInput) (*DiagnosticsResult, error) {
	s.mu.RLock()
	llm := s.llm
	s.mu.RUnlock()

	cacheKey := fmt.Sprintf("tenant:%s:diagnostics:%s", tenantID, clusterID)

	if cached, err := s.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		var result DiagnosticsResult
		if json.Unmarshal(cached, &result) == nil {
			log.Debug().Str("cluster", clusterID).Msg("Diagnostics depuis cache")
			return &result, nil
		}
	}

	model := ""
	if o, ok := llm.(*OllamaClient); ok {
		model = o.ModelName()
	}

	result := &DiagnosticsResult{
		ID:          fmt.Sprintf("diag-%s-%d", clusterID, time.Now().Unix()),
		ClusterID:   clusterID,
		TenantID:    tenantID,
		AnalyzedAt:  time.Now(),
		HealthScore: s.computeHealth(rawIssues),
		Backend:     string(llm.BackendName()),
		Model:       model,
	}

	for _, ri := range rawIssues {
		result.Issues = append(result.Issues, Issue{
			ID:           fmt.Sprintf("%s/%s/%s", ri.Namespace, ri.ResourceName, ri.Type),
			Type:         ri.Type,
			Severity:     ri.Severity,
			Namespace:    ri.Namespace,
			ResourceKind: ri.ResourceKind,
			ResourceName: ri.ResourceName,
			Message:      ri.Message,
			DetectedAt:   time.Now(),
		})
	}

	if len(rawIssues) > 0 {
		prompt := BuildDiagnosticPrompt(rawIssues)
		explanation, err := llm.Analyze(ctx, prompt)
		if err != nil {
			log.Warn().Err(err).Msg("LLM analyze failed")
			result.Summary = "Analyse IA indisponible ??? " + err.Error()
		} else {
			result.Summary = explanation
		}
	} else {
		result.Summary = "Aucun probleme detecte. Cluster en bonne sante."
	}

	if b, err := json.Marshal(result); err == nil {
		s.rdb.Client().Set(ctx, cacheKey, b, 2*time.Minute)
	}
	return result, nil
}

func (s *DiagnosticsService) AnalyzeSingleIssue(ctx context.Context, issue IssueInput) (string, error) {
	s.mu.RLock()
	llm := s.llm
	s.mu.RUnlock()
	prompt := BuildSingleIssuePrompt(issue)
	return llm.Analyze(ctx, prompt)
}

func (s *DiagnosticsService) computeHealth(issues []IssueInput) float64 {
	if len(issues) == 0 { return 100 }
	penalty := 0.0
	for _, i := range issues {
		switch i.Severity {
		case "critical": penalty += 20
		case "warning":  penalty += 8
		default:         penalty += 2
		}
	}
	if s := 100 - penalty; s >= 0 { return s }
	return 0
}
