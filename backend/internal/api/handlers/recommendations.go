package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/analyzer"
	"github.com/k8s-green/backend/internal/api/middleware"
)

func (h *Handler) GetRecommendations(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	cacheKey  := fmt.Sprintf("tenant:%s:reco:%s", tenantID, clusterID)
	ctx       := c.Request.Context()

	// Cache Redis TTL 5min
	if cached, err := h.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", cached)
		return
	}

	// Reconstruire les rapports depuis les métriques en DB
	// Pour l'instant on retourne des recommandations synthétiques depuis le dernier score
	score, err := h.db.TimescaleRepo().GetLatestScore(ctx, clusterID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"recommendations": []interface{}{}})
		return
	}

	// Générer des recommandations synthétiques basées sur le score breakdown
	var recs []analyzer.Recommendation

	if score.CPUEfficiency < 50 {
		recs = append(recs, analyzer.Recommendation{
			Priority:         "HIGH",
			Type:             "rightsizing",
			Title:            "CPU très sous-utilisé — rightsizing critique",
			Description:      fmt.Sprintf("Efficacité CPU: %.0f%% — réduire les requests CPU", score.CPUEfficiency),
			SavingsEurAnnual: score.AnnualWasteEur * 0.6,
			Confidence:       0.90,
		})
	}
	if score.LimitCompliance < 80 {
		recs = append(recs, analyzer.Recommendation{
			Priority:    "HIGH",
			Type:        "missing_limits",
			Title:       "Pods sans resource limits",
			Description: fmt.Sprintf("%.0f%% des pods ont des limits — risque OOM", score.LimitCompliance),
			Confidence:  1.0,
		})
	}
	if score.HPACoverage < 50 {
		recs = append(recs, analyzer.Recommendation{
			Priority:         "MEDIUM",
			Type:             "add_hpa",
			Title:            "Ajouter HPA sur les deployments stateless",
			Description:      fmt.Sprintf("Couverture HPA: %.0f%% — autoscaling manquant", score.HPACoverage),
			SavingsEurAnnual: score.AnnualWasteEur * 0.15,
			Confidence:       0.80,
		})
	}

	resp := gin.H{
		"cluster_id":      clusterID,
		"count":           len(recs),
		"generated_at":    time.Now(),
		"recommendations": recs,
	}

	if b, err := json.Marshal(resp); err == nil {
		h.rdb.Client().Set(ctx, cacheKey, b, 5*time.Minute)
		c.Header("X-Cache", "MISS")
		c.Data(http.StatusOK, "application/json", b)
		return
	}
	c.JSON(http.StatusOK, resp)
}
