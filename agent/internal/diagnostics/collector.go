package diagnostics

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// IssueType categorise les problèmes détectés
type IssueType string

const (
	IssueCrashLoop      IssueType = "CrashLoopBackOff"
	IssueOOMKilled      IssueType = "OOMKilled"
	IssuePending        IssueType = "Pending"
	IssueImagePullError IssueType = "ImagePullError"
	IssueEvicted        IssueType = "Evicted"
	IssueUnschedulable  IssueType = "Unschedulable"
	IssueHighRestart    IssueType = "HighRestartCount"
)

// K8sIssue représente un problème détecté dans le cluster
type K8sIssue struct {
	ID          string    `json:"id"`
	Type        IssueType `json:"type"`
	Severity    string    `json:"severity"` // critical, warning, info
	Namespace   string    `json:"namespace"`
	ResourceKind string   `json:"resource_kind"` // Pod, Node, PVC...
	ResourceName string   `json:"resource_name"`
	Message     string    `json:"message"`
	Details     string    `json:"details"`
	Events      []string  `json:"events,omitempty"`
	DetectedAt  time.Time `json:"detected_at"`
}

// DiagnosticsSnapshot contient tous les problèmes détectés
type DiagnosticsSnapshot struct {
	ClusterID   string     `json:"cluster_id"`
	TenantID    string     `json:"tenant_id"`
	CollectedAt time.Time  `json:"collected_at"`
	Issues      []K8sIssue `json:"issues"`
	HealthScore float64    `json:"health_score"` // 0-100
}

// Collector détecte les problèmes K8s
type Collector struct {
	client kubernetes.Interface
}

func NewCollector(client kubernetes.Interface) *Collector {
	return &Collector{client: client}
}

func (c *Collector) Collect(ctx context.Context, clusterID, tenantID string) (*DiagnosticsSnapshot, error) {
	snap := &DiagnosticsSnapshot{
		ClusterID:   clusterID,
		TenantID:    tenantID,
		CollectedAt: time.Now(),
	}

	// Analyser tous les pods
	pods, err := c.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	for _, pod := range pods.Items {
		issues := c.analyzePod(ctx, &pod)
		snap.Issues = append(snap.Issues, issues...)
	}

	// Analyser les events récents
	events, err := c.client.CoreV1().Events("").List(ctx, metav1.ListOptions{
		FieldSelector: "type=Warning",
	})
	if err == nil {
		issue := c.analyzeEvents(events)
		snap.Issues = append(snap.Issues, issue...)
	}

	// Analyser les nodes
	nodes, err := c.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err == nil {
		issues := c.analyzeNodes(nodes)
		snap.Issues = append(snap.Issues, issues...)
	}

	snap.HealthScore = c.computeHealthScore(snap.Issues)
	return snap, nil
}

func (c *Collector) analyzePod(ctx context.Context, pod *corev1.Pod) []K8sIssue {
	var issues []K8sIssue

	// Pod en phase Failed ou Evicted
	if pod.Status.Phase == corev1.PodFailed {
		reason := pod.Status.Reason
		if reason == "Evicted" {
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/evicted", pod.Namespace, pod.Name),
				Type:         IssueEvicted,
				Severity:     "warning",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Pod %s evicted: %s", pod.Name, pod.Status.Message),
				DetectedAt:   time.Now(),
			})
		}
	}

	// Pod en Pending depuis trop longtemps
	if pod.Status.Phase == corev1.PodPending {
		age := time.Since(pod.CreationTimestamp.Time)
		if age > 5*time.Minute {
			reason := "Unknown"
			if len(pod.Status.Conditions) > 0 {
				reason = string(pod.Status.Conditions[0].Reason)
			}
			issueType := IssuePending
			if strings.Contains(strings.ToLower(reason), "unschedulable") {
				issueType = IssueUnschedulable
			}
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/pending", pod.Namespace, pod.Name),
				Type:         issueType,
				Severity:     "critical",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Pod pending depuis %s: %s", age.Round(time.Second), reason),
				DetectedAt:   time.Now(),
			})
		}
	}

	// Analyser les containers
	for _, cs := range pod.Status.ContainerStatuses {
		// CrashLoopBackOff
		if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/%s/crashloop", pod.Namespace, pod.Name, cs.Name),
				Type:         IssueCrashLoop,
				Severity:     "critical",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Container %s en CrashLoopBackOff (restarts: %d)", cs.Name, cs.RestartCount),
				Details:      cs.State.Waiting.Message,
				DetectedAt:   time.Now(),
			})
		}

		// OOMKilled
		if cs.LastTerminationState.Terminated != nil &&
			cs.LastTerminationState.Terminated.Reason == "OOMKilled" {
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/%s/oom", pod.Namespace, pod.Name, cs.Name),
				Type:         IssueOOMKilled,
				Severity:     "critical",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Container %s tué par OOM — augmenter la memory limit", cs.Name),
				DetectedAt:   time.Now(),
			})
		}

		// ImagePullError
		if cs.State.Waiting != nil &&
			(cs.State.Waiting.Reason == "ImagePullBackOff" || cs.State.Waiting.Reason == "ErrImagePull") {
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/%s/imagepull", pod.Namespace, pod.Name, cs.Name),
				Type:         IssueImagePullError,
				Severity:     "critical",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Impossible de puller l'image pour %s: %s", cs.Name, cs.State.Waiting.Message),
				DetectedAt:   time.Now(),
			})
		}

		// High restart count
		if cs.RestartCount > 10 {
			issues = append(issues, K8sIssue{
				ID:           fmt.Sprintf("%s/%s/%s/restarts", pod.Namespace, pod.Name, cs.Name),
				Type:         IssueHighRestart,
				Severity:     "warning",
				Namespace:    pod.Namespace,
				ResourceKind: "Pod",
				ResourceName: pod.Name,
				Message:      fmt.Sprintf("Container %s a %d redémarrages", cs.Name, cs.RestartCount),
				DetectedAt:   time.Now(),
			})
		}
	}

	return issues
}

func (c *Collector) analyzeEvents(events *corev1.EventList) []K8sIssue {
	var issues []K8sIssue
	seen := map[string]bool{}

	for _, e := range events.Items {
		// Dédupliquer
		key := fmt.Sprintf("%s/%s/%s", e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason)
		if seen[key] { continue }
		seen[key] = true

		// Seulement les events récents (1h)
		if time.Since(e.LastTimestamp.Time) > time.Hour { continue }

		// Ignorer les events normaux
		if e.Reason == "Scheduled" || e.Reason == "Pulled" || e.Reason == "Started" { continue }

		issues = append(issues, K8sIssue{
			ID:           key,
			Type:         IssueType(e.Reason),
			Severity:     "warning",
			Namespace:    e.InvolvedObject.Namespace,
			ResourceKind: e.InvolvedObject.Kind,
			ResourceName: e.InvolvedObject.Name,
			Message:      e.Message,
			DetectedAt:   e.LastTimestamp.Time,
		})
	}
	return issues
}

func (c *Collector) analyzeNodes(nodes *corev1.NodeList) []K8sIssue {
	var issues []K8sIssue
	for _, node := range nodes.Items {
		for _, cond := range node.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status != corev1.ConditionTrue {
				issues = append(issues, K8sIssue{
					ID:           fmt.Sprintf("node/%s/notready", node.Name),
					Type:         IssuePending,
					Severity:     "critical",
					ResourceKind: "Node",
					ResourceName: node.Name,
					Message:      fmt.Sprintf("Node %s NotReady: %s", node.Name, cond.Message),
					DetectedAt:   time.Now(),
				})
			}
		}
	}
	return issues
}

func (c *Collector) computeHealthScore(issues []K8sIssue) float64 {
	if len(issues) == 0 { return 100.0 }
	penalty := 0.0
	for _, issue := range issues {
		switch issue.Severity {
		case "critical": penalty += 15
		case "warning":  penalty += 5
		case "info":     penalty += 1
		}
	}
	score := 100.0 - penalty
	if score < 0 { return 0 }
	return score
}
