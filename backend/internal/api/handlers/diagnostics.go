package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/ai"
	"github.com/k8s-green/backend/internal/api/middleware"
)

// GetDiagnostics retourne le diagnostic IA du cluster
func (h *Handler) GetDiagnostics(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	ctx       := c.Request.Context()

	// Recuperer les issues stockees en DB depuis la derniere collecte agent
	rows, err := h.db.Pool.Query(ctx,
		`SELECT issue_type, severity, namespace, resource_kind, resource_name, message, details
		 FROM cluster_issues
		 WHERE cluster_id = $1 AND detected_at > NOW() - INTERVAL '1 hour'
		 ORDER BY severity DESC, detected_at DESC
		 LIMIT 20`,
		clusterID,
	)

	var issues []ai.IssueInput
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var issue ai.IssueInput
			rows.Scan(&issue.Type, &issue.Severity, &issue.Namespace,
				&issue.ResourceKind, &issue.ResourceName, &issue.Message, &issue.Details)
			issues = append(issues, issue)
		}
	}

	// Si pas de donnees DB, retourner un etat "pas encore scanne"
	if len(issues) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"cluster_id":   clusterID,
			"health_score": 100,
			"issues":       []interface{}{},
			"summary":      "Aucun probleme detecte. L agent n a pas encore envoye de donnees de diagnostic.",
			"analyzed_at":  time.Now(),
		})
		return
	}

	// Analyser avec le LLM si configure
	if h.aiSvc != nil {
		result, err := h.aiSvc.AnalyzeIssues(ctx, clusterID, tenantID, issues)
		if err == nil {
			c.JSON(http.StatusOK, result)
			return
		}
	}

	// Fallback sans LLM : retourner les issues brutes
	c.JSON(http.StatusOK, gin.H{
		"cluster_id":   clusterID,
		"health_score": 80,
		"issues":       issues,
		"summary":      "LLM non configure - analyse IA indisponible. Configurez AI_BACKEND dans les secrets.",
		"analyzed_at":  time.Now(),
	})
}

// AnalyzeSingleIssue analyse un probleme specifique en detail
func (h *Handler) AnalyzeSingleIssue(c *gin.Context) {
	var req ai.IssueInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.aiSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "LLM non configure. Ajoutez AI_BACKEND, AI_API_KEY dans les secrets Kubernetes.",
		})
		return
	}

	explanation, err := h.aiSvc.AnalyzeSingleIssue(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Analyse LLM echouee: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issue":       req,
		"explanation": explanation,
		"analyzed_at": time.Now(),
	})
}
