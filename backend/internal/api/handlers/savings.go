package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
)

func (h *Handler) GetSavings(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	cacheKey  := fmt.Sprintf("tenant:%s:savings:%s", tenantID, clusterID)
	ctx       := c.Request.Context()

	if cached, err := h.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", cached)
		return
	}

	score, err := h.db.TimescaleRepo().GetLatestScore(ctx, clusterID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"cluster_id": clusterID, "annual_savings_eur": 0})
		return
	}

	// Ventilation des économies potentielles
	annual := score.AnnualWasteEur
	resp := gin.H{
		"cluster_id":          clusterID,
		"annual_savings_eur":  annual,
		"monthly_savings_eur": annual / 12,
		"breakdown": gin.H{
			"rightsizing_eur":        annual * 0.60, // 60% du gaspillage = surdimensionnement
			"node_consolidation_eur": annual * 0.25, // 25% = consolidation nœuds
			"hpa_automation_eur":     annual * 0.15, // 15% = autoscaling
		},
		"computed_at": time.Now(),
	}

	if b, err := json.Marshal(resp); err == nil {
		h.rdb.Client().Set(ctx, cacheKey, b, 5*time.Minute)
		c.Header("X-Cache", "MISS")
		c.Data(http.StatusOK, "application/json", b)
		return
	}
	c.JSON(http.StatusOK, resp)
}
