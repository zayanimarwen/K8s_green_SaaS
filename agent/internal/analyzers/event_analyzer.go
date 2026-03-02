package analyzers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventAnalyzer struct{ c *Clients }
func NewEventAnalyzer(c *Clients) *EventAnalyzer { return &EventAnalyzer{c} }
func (a *EventAnalyzer) Name() string   { return "eventAnalyzer" }
func (a *EventAnalyzer) Optional() bool { return false }

// Events a ignorer
var ignoredReasons = map[string]bool{
	"Scheduled": true, "Pulled": true, "Started": true, "Created": true,
	"Killing": true, "Preempting": true, "ScalingReplicaSet": true,
}

func (a *EventAnalyzer) Analyze(ctx context.Context) ([]Issue, error) {
	var issues []Issue
	events, err := a.c.K8s.CoreV1().Events("").List(ctx, metav1.ListOptions{
		FieldSelector: "type=Warning",
	})
	if err != nil { return nil, err }

	seen := map[string]bool{}
	for _, e := range events.Items {
		if ignoredReasons[e.Reason] { continue }
		if time.Since(e.LastTimestamp.Time) > 2*time.Hour { continue }

		key := e.InvolvedObject.Namespace + "/" + e.InvolvedObject.Name + "/" + e.Reason
		if seen[key] { continue }
		seen[key] = true

		sev := Warning
		if e.Count > 10 { sev = Critical }

		issues = append(issues, issue(a.Name(),
			e.InvolvedObject.Kind,
			e.InvolvedObject.Namespace,
			e.InvolvedObject.Name,
			e.Reason, sev, e.Message))
	}
	return issues, nil
}
