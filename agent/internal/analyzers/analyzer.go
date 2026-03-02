package analyzers

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/dynamic"
)

// Severity niveaux
const (
	Critical = "critical"
	Warning  = "warning"
	Info     = "info"
)

// Issue represente un probleme detecte
type Issue struct {
	ID           string    `json:"id"`
	Analyzer     string    `json:"analyzer"`
	Type         string    `json:"type"`
	Severity     string    `json:"severity"`
	Namespace    string    `json:"namespace"`
	ResourceKind string    `json:"resource_kind"`
	ResourceName string    `json:"resource_name"`
	Message      string    `json:"message"`
	Details      string    `json:"details,omitempty"`
	DetectedAt   time.Time `json:"detected_at"`
}

// Analyzer interface que tous les analyseurs implementent
type Analyzer interface {
	Name() string
	Analyze(ctx context.Context) ([]Issue, error)
	Optional() bool
}

// Clients K8s partages entre analyseurs
type Clients struct {
	K8s     kubernetes.Interface
	Dynamic dynamic.Interface
}

// issue helper
func issue(analyzer, kind, ns, name, issueType, severity, msg string) Issue {
	return Issue{
		ID:           ns + "/" + name + "/" + issueType,
		Analyzer:     analyzer,
		Type:         issueType,
		Severity:     severity,
		Namespace:    ns,
		ResourceKind: kind,
		ResourceName: name,
		Message:      msg,
		DetectedAt:   time.Now(),
	}
}
