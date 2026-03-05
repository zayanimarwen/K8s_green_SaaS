package handlers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
)

type AgentMetricsSnapshot struct {
	ClusterID   string        `json:"cluster_id"`
	TenantID    string        `json:"tenant_id"`
	CollectedAt time.Time     `json:"collected_at"`
	Pods        []AgentPod    `json:"pods"`
	Nodes       []AgentNode   `json:"nodes"`
	Deployments []AgentDeploy `json:"deployments"`
	Namespaces  []AgentNS     `json:"namespaces"`
}

type AgentPod struct {
	PodName       string            `json:"PodName"`
	ContainerName string            `json:"ContainerName"`
	Namespace     string            `json:"Namespace"`
	NodeName      string            `json:"NodeName"`
	Labels        map[string]string `json:"Labels"`
	CPURequestM   float64           `json:"CPURequestM"`
	CPULimitM     float64           `json:"CPULimitM"`
	MemRequestMi  float64           `json:"MemRequestMi"`
	MemLimitMi    float64           `json:"MemLimitMi"`
	HasLimits     bool              `json:"HasLimits"`
	RestartCount  int32             `json:"RestartCount"`
}

type AgentNode struct {
	Name             string  `json:"Name"`
	CPUAllocatableM  float64 `json:"CPUAllocatableM"`
	MemAllocatableMi float64 `json:"MemAllocatableMi"`
}

type AgentDeploy struct {
	Name      string `json:"Name"`
	Namespace string `json:"Namespace"`
	Replicas  *int32 `json:"Replicas"`
	Ready     int32  `json:"ReadyReplicas"`
	HasHPA    bool   `json:"HasHPA"`
}

type AgentNS struct {
	Name string `json:"Name"`
}

func (h *Handler) IngestMetrics(c *gin.Context) {
	tenantID := middleware.TenantID(c)
	_ = tenantID

	var snap AgentMetricsSnapshot
	if err := c.ShouldBindJSON(&snap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Stocker les métriques pods
	for _, pod := range snap.Pods {
		h.db.Pool.Exec(ctx,
			`INSERT INTO pod_metrics (
				time, cluster_id, pod_name, container_name, namespace, node_name,
				cpu_request_m, cpu_limit_m, mem_request_mi, mem_limit_mi,
				has_limits, restart_count
			) VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT DO NOTHING`,
			snap.ClusterID, pod.PodName, pod.ContainerName, pod.Namespace, pod.NodeName,
			pod.CPURequestM, pod.CPULimitM, pod.MemRequestMi, pod.MemLimitMi,
			pod.HasLimits, pod.RestartCount,
		)
	}

	// 2. Supprimer les anciennes issues non résolues de ce cluster
	h.db.Pool.Exec(ctx,
		`DELETE FROM cluster_issues WHERE cluster_id = $1 AND resolved = false`,
		snap.ClusterID,
	)

	// 3. Analyser et générer des issues
	issueCount := 0
	for _, pod := range snap.Pods {
		// Pod sans limits
		if !pod.HasLimits {
			h.db.Pool.Exec(ctx,
				`INSERT INTO cluster_issues (cluster_id, issue_type, severity, namespace, resource_kind, resource_name, message, details, detected_at)
				 VALUES ($1, 'MissingLimits', 'warning', $2, 'Pod', $3, $4, $5, NOW())`,
				snap.ClusterID, pod.Namespace, pod.PodName,
				fmt.Sprintf("Pod %s/%s n'a pas de CPU/Memory limits définis", pod.Namespace, pod.PodName),
				"Sans limits, ce pod peut consommer toutes les ressources du node et provoquer des OOM kills.",
			)
			issueCount++
		}

		// Pod en CrashLoop (restarts élevés)
		if pod.RestartCount >= 5 {
			severity := "warning"
			if pod.RestartCount >= 10 {
				severity = "critical"
			}
			h.db.Pool.Exec(ctx,
				`INSERT INTO cluster_issues (cluster_id, issue_type, severity, namespace, resource_kind, resource_name, message, details, detected_at)
				 VALUES ($1, 'HighRestarts', $2, $3, 'Pod', $4, $5, $6, NOW())`,
				snap.ClusterID, severity, pod.Namespace, pod.PodName,
				fmt.Sprintf("Pod %s/%s a redémarré %d fois", pod.Namespace, pod.PodName, pod.RestartCount),
				fmt.Sprintf("Restarts: %d. Vérifiez les logs avec: kubectl logs %s -n %s", pod.RestartCount, pod.PodName, pod.Namespace),
			)
			issueCount++
		}

		// Pod surdimensionné (request très élevée)
		if pod.CPURequestM > 1000 && pod.CPULimitM > 0 {
			ratio := pod.CPULimitM / pod.CPURequestM
			if ratio > 4 {
				h.db.Pool.Exec(ctx,
					`INSERT INTO cluster_issues (cluster_id, issue_type, severity, namespace, resource_kind, resource_name, message, details, detected_at)
					 VALUES ($1, 'Oversized', 'info', $2, 'Pod', $3, $4, $5, NOW())`,
					snap.ClusterID, pod.Namespace, pod.PodName,
					fmt.Sprintf("Pod %s/%s potentiellement surdimensionné (ratio limit/request: %.1fx)", pod.Namespace, pod.PodName, ratio),
					fmt.Sprintf("CPU Request: %.0fm, CPU Limit: %.0fm. Envisagez de réduire les requests.", pod.CPURequestM, pod.CPULimitM),
				)
				issueCount++
			}
		}
	}

	// 4. Mettre à jour last_seen_at
	h.db.Pool.Exec(ctx,
		`UPDATE clusters SET last_seen_at = NOW() WHERE id = $1`,
		snap.ClusterID,
	)

	c.JSON(http.StatusOK, gin.H{
		"cluster_id":   snap.ClusterID,
		"pods_ingested": len(snap.Pods),
		"nodes":        len(snap.Nodes),
		"issues_found": issueCount,
		"collected_at": snap.CollectedAt,
	})
}
