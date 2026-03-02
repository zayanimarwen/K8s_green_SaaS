package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
)

func (h *Handler) GetCarbon(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	cacheKey  := fmt.Sprintf("tenant:%s:carbon:%s", tenantID, clusterID)
	ctx       := c.Request.Context()

	if cached, err := h.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", cached)
		return
	}

	// Récupérer le dernier score qui contient co2_kg_annual
	score, err := h.db.TimescaleRepo().GetLatestScore(ctx, clusterID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"cluster_id": clusterID,
			"co2_kg_annual": 0,
			"message": "En attente du premier scan",
		})
		return
	}

	// Equivalences carbone
	resp := gin.H{
		"cluster_id":       clusterID,
		"co2_kg_annual":    score.CO2KgAnnual,
		"kwh_annual":       score.CO2KgAnnual / 0.065,           // 65gCO2/kWh France
		"equivalent_km_car": score.CO2KgAnnual * 6.5,            // 154gCO2/km voiture essence
		"equivalent_trees": score.CO2KgAnnual / 21.77,           // ~21.77kgCO2/arbre/an
		"provider":         "on-prem",
		"region":           "fr",
		"carbon_intensity": 65,                                  // gCO2/kWh France RTE 2023
		"computed_at":      time.Now(),
	}

	if b, err := json.Marshal(resp); err == nil {
		h.rdb.Client().Set(ctx, cacheKey, b, 5*time.Minute)
		c.Header("X-Cache", "MISS")
		c.Data(http.StatusOK, "application/json", b)
		return
	}
	c.JSON(http.StatusOK, resp)
}
